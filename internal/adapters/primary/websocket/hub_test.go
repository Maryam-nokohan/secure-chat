package websocket

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/maryam-nokohan/secure-chat/pkg"
)

func newTestHubClient(id, username string, bufferSize int) *Client {
	return &Client{
		ID:       id,
		Username: username,
		Send:     make(chan []byte, bufferSize),
	}
}

func startTestHub(t *testing.T) *Hub {
	t.Helper()
	pkg.Init()

	hub := NewHub()

	go hub.Run()

	return hub
}

func waitForMessage(
	t *testing.T,
	ch <-chan []byte,
	timeout time.Duration,
) []byte {
	t.Helper()

	select {
	case msg := <-ch:
		return msg
	case <-time.After(timeout):
		t.Fatal("timed out waiting for message")
		return nil
	}
}

func assertNoMessage(
	t *testing.T,
	ch <-chan []byte,
	timeout time.Duration,
) {
	t.Helper()

	select {
	case msg := <-ch:
		t.Fatalf("expected no message, but received %q", string(msg))
	case <-time.After(timeout):
	}
}

func waitForClientRegistration(
	t *testing.T,
	hub *Hub,
	client *Client,
) {
	t.Helper()

	hub.Register <- client

	deadline := time.After(time.Second)

	for {
		select {
		case <-deadline:
			t.Fatalf(
				"client %s was not registered within timeout",
				client.ID,
			)

		default:
			hub.mu.Lock()
			_, exists := hub.Clients[client.ID]
			hub.mu.Unlock()

			if exists {
				return
			}

			time.Sleep(time.Millisecond)
		}
	}
}

func TestNewHub_InitializesAllCollectionsAndChannels(t *testing.T) {
	hub := NewHub()

	if hub == nil {
		t.Fatal("expected non-nil hub")
	}

	if hub.Clients == nil {
		t.Fatal("Clients map must be initialized")
	}

	if hub.Rooms == nil {
		t.Fatal("Rooms map must be initialized")
	}

	if hub.Register == nil {
		t.Fatal("Register channel must be initialized")
	}

	if hub.Unregister == nil {
		t.Fatal("Unregister channel must be initialized")
	}

	if hub.Broadcast == nil {
		t.Fatal("Broadcast channel must be initialized")
	}
}

func TestHub_RegisterClient(t *testing.T) {
	hub := startTestHub(t)

	client := newTestHubClient("user-1", "alice", 10)

	waitForClientRegistration(t, hub, client)

	hub.mu.Lock()
	registeredClient, exists := hub.Clients["user-1"]
	hub.mu.Unlock()

	if !exists {
		t.Fatal("client was not registered")
	}

	if registeredClient != client {
		t.Fatal("registered client pointer does not match")
	}
}

func TestHub_RegisterMultipleClients(t *testing.T) {
	hub := startTestHub(t)

	client1 := newTestHubClient("user-1", "alice", 10)
	client2 := newTestHubClient("user-2", "bob", 10)
	client3 := newTestHubClient("user-3", "charlie", 10)

	waitForClientRegistration(t, hub, client1)
	waitForClientRegistration(t, hub, client2)
	waitForClientRegistration(t, hub, client3)

	hub.mu.Lock()
	defer hub.mu.Unlock()

	if len(hub.Clients) != 3 {
		t.Fatalf(
			"expected 3 registered clients, got %d",
			len(hub.Clients),
		)
	}
}

func TestHub_JoinRoom(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "room-1")

	hub.mu.Lock()
	defer hub.mu.Unlock()

	if client.Room != "room-1" {
		t.Fatalf(
			"expected client room room-1, got %q",
			client.Room,
		)
	}

	room, exists := hub.Rooms["room-1"]

	if !exists {
		t.Fatal("room was not created")
	}

	if room["user-1"] != client {
		t.Fatal("client was not added to room")
	}
}

func TestHub_JoinRoomCreatesRoomAutomatically(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "new-room")

	hub.mu.Lock()
	_, exists := hub.Rooms["new-room"]
	hub.mu.Unlock()

	if !exists {
		t.Fatal("JoinRoom should create a missing room")
	}
}

func TestHub_JoinRoomSwitchesClientFromPreviousRoom(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "room-1")
	hub.JoinRoom(client, "room-2")

	hub.mu.Lock()
	defer hub.mu.Unlock()

	if client.Room != "room-2" {
		t.Fatalf(
			"expected current room room-2, got %q",
			client.Room,
		)
	}

	if _, exists := hub.Rooms["room-1"]["user-1"]; exists {
		t.Fatal("client must be removed from previous room")
	}

	if hub.Rooms["room-2"]["user-1"] != client {
		t.Fatal("client was not added to new room")
	}
}

func TestHub_JoinRoomMovesClientWithoutDuplicatingMembership(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "room-1")
	hub.JoinRoom(client, "room-1")

	hub.mu.Lock()
	defer hub.mu.Unlock()

	if len(hub.Rooms["room-1"]) != 1 {
		t.Fatalf(
			"expected exactly one room membership, got %d",
			len(hub.Rooms["room-1"]),
		)
	}
}

func TestHub_BroadcastToRoomOnlyReachesRoomMembers(t *testing.T) {
	hub := NewHub()

	alice := newTestHubClient("alice", "alice", 10)
	bob := newTestHubClient("bob", "bob", 10)
	charlie := newTestHubClient("charlie", "charlie", 10)

	hub.JoinRoom(alice, "room-a")
	hub.JoinRoom(bob, "room-a")
	hub.JoinRoom(charlie, "room-b")

	payload := []byte(`{"type":"message","text":"hello"}`)

	hub.BroadcastToRoom("room-a", payload)

	gotAlice := waitForMessage(t, alice.Send, time.Second)
	gotBob := waitForMessage(t, bob.Send, time.Second)

	if string(gotAlice) != string(payload) {
		t.Fatalf(
			"alice received wrong payload: %s",
			gotAlice,
		)
	}

	if string(gotBob) != string(payload) {
		t.Fatalf(
			"bob received wrong payload: %s",
			gotBob,
		)
	}

	assertNoMessage(t, charlie.Send, 50*time.Millisecond)
}

func TestHub_BroadcastToRoomDoesNotLeakMessagesBetweenRooms(t *testing.T) {
	hub := NewHub()

	roomAUser := newTestHubClient("user-a", "alice", 10)
	roomBUser := newTestHubClient("user-b", "bob", 10)

	hub.JoinRoom(roomAUser, "room-a")
	hub.JoinRoom(roomBUser, "room-b")

	payload := []byte(`{"type":"message","room_id":"room-a"}`)

	hub.BroadcastToRoom("room-a", payload)

	got := waitForMessage(t, roomAUser.Send, time.Second)

	if string(got) != string(payload) {
		t.Fatalf("room A received wrong payload")
	}

	assertNoMessage(t, roomBUser.Send, 50*time.Millisecond)
}

func TestHub_BroadcastToRoomWithUnknownRoomDoesNothing(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "room-a")

	hub.BroadcastToRoom(
		"does-not-exist",
		[]byte(`{"type":"message"}`),
	)

	assertNoMessage(t, client.Send, 50*time.Millisecond)
}

func TestHub_BroadcastToRoomWithFullClientBufferDoesNotBlock(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 1)

	hub.JoinRoom(client, "room-a")

	client.Send <- []byte("already-full")

	done := make(chan struct{})

	go func() {
		hub.BroadcastToRoom(
			"room-a",
			[]byte(`{"type":"message","text":"new"}`),
		)

		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("BroadcastToRoom blocked on a full client buffer")
	}

	got := waitForMessage(t, client.Send, time.Second)

	if string(got) != "already-full" {
		t.Fatalf(
			"existing buffered message was unexpectedly changed: %q",
			string(got),
		)
	}
}

func TestHub_GetOnlineUserIDs(t *testing.T) {
	hub := NewHub()

	alice := newTestHubClient("alice", "alice", 10)
	bob := newTestHubClient("bob", "bob", 10)

	hub.mu.Lock()
	hub.Clients[alice.ID] = alice
	hub.Clients[bob.ID] = bob
	hub.mu.Unlock()

	online := hub.GetOnlineUserIDs()

	if len(online) != 2 {
		t.Fatalf(
			"expected 2 online users, got %d",
			len(online),
		)
	}

	if !online["alice"] {
		t.Fatal("alice should be online")
	}

	if !online["bob"] {
		t.Fatal("bob should be online")
	}

	if online["charlie"] {
		t.Fatal("charlie should not be online")
	}
}

func TestHub_GetOnlineUserIDsReturnsCopy(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("alice", "alice", 10)

	hub.mu.Lock()
	hub.Clients[client.ID] = client
	hub.mu.Unlock()

	online := hub.GetOnlineUserIDs()

	delete(online, "alice")

	hub.mu.Lock()
	_, stillOnline := hub.Clients["alice"]
	hub.mu.Unlock()

	if !stillOnline {
		t.Fatal(
			"modifying GetOnlineUserIDs result must not modify hub state",
		)
	}
}

func TestHub_UnregisterRemovesClient(t *testing.T) {
	hub := startTestHub(t)

	client := newTestHubClient("user-1", "alice", 10)

	waitForClientRegistration(t, hub, client)

	hub.Unregister <- client

	select {
	case _, ok := <-client.Send:
		if ok {
			t.Fatal("expected Send channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("client Send channel was not closed")
	}

	hub.mu.Lock()
	_, exists := hub.Clients["user-1"]
	hub.mu.Unlock()

	if exists {
		t.Fatal("client still exists after unregister")
	}
}

func TestHub_UnregisterRemovesClientFromAllRooms(t *testing.T) {
	hub := startTestHub(t)

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "room-a")
	hub.JoinRoom(client, "room-b")

	waitForClientRegistration(t, hub, client)

	hub.Unregister <- client

	time.Sleep(50 * time.Millisecond)

	hub.mu.Lock()
	defer hub.mu.Unlock()

	for roomID, room := range hub.Rooms {
		if _, exists := room[client.ID]; exists {
			t.Fatalf(
				"client still exists in room %s after unregister",
				roomID,
			)
		}
	}
}

func TestHub_UnregisterUnknownClientDoesNotCloseSendChannel(t *testing.T) {
	hub := startTestHub(t)

	client := newTestHubClient("not-registered", "alice", 10)

	hub.Unregister <- client

	time.Sleep(50 * time.Millisecond)

	select {
	case client.Send <- []byte("still-open"):
	default:
		t.Fatal("Send channel should not have been closed")
	}
}

func TestHub_GlobalBroadcastReachesEveryClient(t *testing.T) {
	hub := startTestHub(t)

	alice := newTestHubClient("alice", "alice", 10)
	bob := newTestHubClient("bob", "bob", 10)
	charlie := newTestHubClient("charlie", "charlie", 10)

	waitForClientRegistration(t, hub, alice)
	waitForClientRegistration(t, hub, bob)
	waitForClientRegistration(t, hub, charlie)

	payload := []byte(`{"type":"global","text":"hello everyone"}`)

	hub.Broadcast <- payload

	for _, client := range []*Client{
		alice,
		bob,
		charlie,
	} {
		got := waitForMessage(t, client.Send, time.Second)

		if string(got) != string(payload) {
			t.Fatalf(
				"client %s received wrong payload: %s",
				client.ID,
				got,
			)
		}
	}
}

func TestHub_GlobalBroadcastReachesClientsRegardlessOfRoom(t *testing.T) {
	hub := startTestHub(t)

	alice := newTestHubClient("alice", "alice", 10)
	bob := newTestHubClient("bob", "bob", 10)

	hub.JoinRoom(alice, "room-a")
	hub.JoinRoom(bob, "room-b")

	waitForClientRegistration(t, hub, alice)
	waitForClientRegistration(t, hub, bob)

	payload := []byte(`{"type":"global"}`)

	hub.Broadcast <- payload

	gotAlice := waitForMessage(t, alice.Send, time.Second)
	gotBob := waitForMessage(t, bob.Send, time.Second)

	if string(gotAlice) != string(payload) {
		t.Fatal("alice received wrong broadcast")
	}

	if string(gotBob) != string(payload) {
		t.Fatal("bob received wrong broadcast")
	}
}

func TestHub_GlobalBroadcastWithFullClientBufferDoesNotBlockHub(t *testing.T) {
	hub := startTestHub(t)

	full := newTestHubClient("full", "full", 1)
	normal := newTestHubClient("normal", "normal", 10)

	waitForClientRegistration(t, hub, full)
	waitForClientRegistration(t, hub, normal)

	full.Send <- []byte("already-full")

	payload := []byte(`{"type":"broadcast"}`)

	hub.Broadcast <- payload

	got := waitForMessage(t, normal.Send, time.Second)

	if string(got) != string(payload) {
		t.Fatalf("normal client did not receive broadcast")
	}

	hub.mu.Lock()
	_, fullStillRegistered := hub.Clients["full"]
	hub.mu.Unlock()

	if fullStillRegistered {
		t.Fatal(
			"client with full Send buffer should be removed by global broadcast",
		)
	}
}

func TestHub_PresenceEventFormat(t *testing.T) {
	hub := NewHub()

	existing := newTestHubClient("existing", "existing-user", 10)
	newClient := newTestHubClient("new", "new-user", 10)

	hub.JoinRoom(existing, "room-a")
	hub.JoinRoom(newClient, "room-a")

	hub.mu.Lock()
	hub.Clients[existing.ID] = existing
	hub.mu.Unlock()

	hub.broadcastPresence(newClient, true)

	raw := waitForMessage(t, existing.Send, time.Second)

	var event PresenceEvent

	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf(
			"presence event is not valid JSON: %v",
			err,
		)
	}

	if event.Type != "presence" {
		t.Fatalf(
			"expected event type presence, got %q",
			event.Type,
		)
	}

	if event.UserID != "new" {
		t.Fatalf(
			"expected user ID new, got %q",
			event.UserID,
		)
	}

	if event.Username != "new-user" {
		t.Fatalf(
			"expected username new-user, got %q",
			event.Username,
		)
	}

	if !event.Online {
		t.Fatal("expected online=true")
	}
}

func TestHub_PresenceOnlineOnlyReachesUsersInSameRoom(t *testing.T) {
	hub := NewHub()

	existingA := newTestHubClient("existing-a", "alice", 10)
	existingB := newTestHubClient("existing-b", "bob", 10)
	newClient := newTestHubClient("new", "new-user", 10)

	hub.JoinRoom(existingA, "room-a")
	hub.JoinRoom(existingB, "room-b")
	hub.JoinRoom(newClient, "room-a")

	hub.mu.Lock()
	hub.Clients[existingA.ID] = existingA
	hub.Clients[existingB.ID] = existingB
	hub.mu.Unlock()

	hub.broadcastPresence(newClient, true)

	_ = waitForMessage(t, existingA.Send, time.Second)

	assertNoMessage(t, existingB.Send, 50*time.Millisecond)
}

func TestHub_PresenceOnlineDoesNotSendEventToTheJoiningClient(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "room-a")

	hub.mu.Lock()
	hub.Clients[client.ID] = client
	hub.mu.Unlock()

	hub.broadcastPresence(client, true)

	assertNoMessage(t, client.Send, 50*time.Millisecond)
}

func TestHub_PresenceOfflineReachesRemainingRoomMembers(t *testing.T) {
	hub := NewHub()

	alice := newTestHubClient("alice", "alice", 10)
	bob := newTestHubClient("bob", "bob", 10)

	hub.JoinRoom(alice, "room-a")
	hub.JoinRoom(bob, "room-a")

	hub.mu.Lock()
	hub.Clients[alice.ID] = alice
	hub.Clients[bob.ID] = bob
	hub.mu.Unlock()

	hub.broadcastPresence(alice, false)

	raw := waitForMessage(t, bob.Send, time.Second)

	var event PresenceEvent

	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf(
			"invalid presence JSON: %v",
			err,
		)
	}

	if event.Type != "presence" {
		t.Fatalf("expected presence event")
	}

	if event.UserID != "alice" {
		t.Fatalf(
			"expected alice, got %q",
			event.UserID,
		)
	}

	if event.Online {
		t.Fatal("expected offline event")
	}

	assertNoMessage(t, alice.Send, 50*time.Millisecond)
}

func TestHub_PresenceOfflineReachesUsersInOtherRooms(t *testing.T) {
	hub := NewHub()

	alice := newTestHubClient("alice", "alice", 10)
	bob := newTestHubClient("bob", "bob", 10)

	hub.JoinRoom(alice, "room-a")
	hub.JoinRoom(bob, "room-b")

	hub.mu.Lock()
	hub.Clients[alice.ID] = alice
	hub.Clients[bob.ID] = bob
	hub.mu.Unlock()

	hub.broadcastPresence(alice, false)

	raw := waitForMessage(t, bob.Send, time.Second)

	var event PresenceEvent

	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf(
			"invalid presence JSON: %v",
			err,
		)
	}

	if event.UserID != "alice" {
		t.Fatalf("expected alice presence event")
	}

	if event.Online {
		t.Fatal("expected offline event")
	}
}

func TestHub_BroadcastToRoomDoesNotSendToUnregisteredClient(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "room-a")

	hub.BroadcastToRoom(
		"room-a",
		[]byte(`{"type":"message"}`),
	)

	got := waitForMessage(t, client.Send, time.Second)

	if string(got) != `{"type":"message"}` {
		t.Fatalf("unexpected message: %s", got)
	}
}

func TestHub_LeavingRoomByJoiningAnotherRoomStopsOldRoomMessages(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)
	other := newTestHubClient("user-2", "bob", 10)

	hub.JoinRoom(client, "room-a")
	hub.JoinRoom(other, "room-b")

	hub.JoinRoom(client, "room-b")

	hub.BroadcastToRoom(
		"room-a",
		[]byte(`{"room":"room-a"}`),
	)

	assertNoMessage(t, client.Send, 50*time.Millisecond)

	hub.BroadcastToRoom(
		"room-b",
		[]byte(`{"room":"room-b"}`),
	)

	gotClient := waitForMessage(t, client.Send, time.Second)
	gotOther := waitForMessage(t, other.Send, time.Second)

	if string(gotClient) != `{"room":"room-b"}` {
		t.Fatalf(
			"client received wrong room message: %s",
			gotClient,
		)
	}

	if string(gotOther) != `{"room":"room-b"}` {
		t.Fatalf(
			"other client received wrong room message: %s",
			gotOther,
		)
	}
}

func TestHub_UnregisterClosesClientSendChannelOnlyOnce(t *testing.T) {
	hub := startTestHub(t)

	client := newTestHubClient("user-1", "alice", 10)

	waitForClientRegistration(t, hub, client)

	hub.Unregister <- client

	select {
	case _, ok := <-client.Send:
		if ok {
			t.Fatal("expected closed Send channel")
		}
	case <-time.After(time.Second):
		t.Fatal("Send channel was not closed")
	}
	hub.Unregister <- client

	time.Sleep(50 * time.Millisecond)
}

func TestHub_RoomBroadcastPreservesPayloadExactly(t *testing.T) {
	hub := NewHub()

	client := newTestHubClient("user-1", "alice", 10)

	hub.JoinRoom(client, "room-a")

	payload := []byte{
		0,
		1,
		2,
		3,
		255,
		'{',
		'}',
	}

	hub.BroadcastToRoom("room-a", payload)

	got := waitForMessage(t, client.Send, time.Second)

	if len(got) != len(payload) {
		t.Fatalf(
			"payload length changed: expected %d, got %d",
			len(payload),
			len(got),
		)
	}

	for i := range payload {
		if got[i] != payload[i] {
			t.Fatalf(
				"payload changed at byte %d: expected %d, got %d",
				i,
				payload[i],
				got[i],
			)
		}
	}
}
