let ws;
let currentRoom = null;
let currentRoomName = "";
let roomJoinReady = false;
let onlineUsers = {};

const chatBox = document.getElementById("chat-box");
const messageInput = document.getElementById("message-input");

function fmtTime(iso) {
  if (!iso) return "";
  const d = new Date(iso);
  if (isNaN(d.getTime())) return iso;
  return d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

// WebSocket
function connectWebSocket() {
  const proto = location.protocol === "https:" ? "wss://" : "ws://";
  ws = new WebSocket(`${proto}${location.host}/ws`);

  ws.onopen = () => {
    reconnectAttempts = 0;
    setConnStatus(true);
    if (currentRoom) sendJoin(currentRoom);
  };

  ws.onclose = () => {
    setConnStatus(false);
    roomJoinReady = false;
    messageInput.disabled = true;
    scheduleReconnect();
  };
  ws.onmessage = async ({ data }) => {
    let msg;
    try {
      msg = JSON.parse(data);
    } catch {
      return;
    }

    switch (msg.type) {
      case "joined":
        if (msg.room_id === currentRoom) {
          roomJoinReady = true;
          messageInput.disabled = false;
          messageInput.placeholder = "Type a message…";
        }
        break;
      case "message":
        if (msg.room_id === currentRoom) await appendLiveMessage(msg);
        break;
      case "member_joined":
        if (msg.room_id === currentRoom) {

          await loadRoomProfile(msg.room_id);
          if (msg.user_id !== CURRENT_UID) {
            appendSystemMessage(`${msg.username} joined the room`);
          }
        }
        break;

      case "presence":
        if (msg.online) {
          onlineUsers[msg.user_id] = msg.username;
        } else {
          delete onlineUsers[msg.user_id];
        }
        refreshMembersPanel();

        break;
    }
  };
}

function sendJoin(roomID) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    roomJoinReady = false;
    messageInput.disabled = true;
    messageInput.placeholder = "Joining room…";
    ws.send(JSON.stringify({ type: "join", room_id: roomID }));
  }
}

async function sendMessage(e) {
  e.preventDefault();
  const text = messageInput.value.trim();
  if (!text) return;
  if (
    !ws ||
    ws.readyState !== WebSocket.OPEN ||
    !roomJoinReady ||
    !currentRoom
  ) {
    appendSystemMessage("You're offline — message not sent. Reconnecting…");
    return;
  }
  if (!myPrivateKey) {
    appendSystemMessage(
      "Encryption is locked. Please unlock to send messages.",
    );
    showUnlockPrompt();
    return;
  }

  if (!currentRoomProfile || currentRoomProfile.id !== currentRoom) {
    appendSystemMessage("Still loading room info — try again in a moment.");
    return;
  }

  const { ciphertext, nonce, rawKey } = await encryptMessage(text);
  const keys = {};
  for (const m of currentRoomProfile.members) {
    const pubKey = await getMemberPublicKey(m.id, m.public_key);
    if (!pubKey) continue;
    keys[m.id] = await encryptKeyForRecipient(rawKey, pubKey);
  }

  ws.send(
    JSON.stringify({
      type: "message",
      room_id: currentRoom,
      ciphertext,
      nonce,
      keys,
    }),
  );
  messageInput.value = "";
}

function setConnStatus(online) {
  document.getElementById("conn-status").innerHTML = online
    ? '<i class="bi bi-circle-fill text-success"></i> Online'
    : '<i class="bi bi-circle-fill text-danger"></i> Offline';
}

// Room switching
async function switchRoom(roomID, roomName) {
  if (currentRoom === roomID) return;

  currentRoom = roomID;
  currentRoomName = roomName;
  roomJoinReady = false;
  messageInput.disabled = true;

  document.getElementById("active-room-title").textContent = `# ${roomName}`;
  document.getElementById("invite-btn").classList.remove("d-none");
  document
    .querySelectorAll(".room-item")
    .forEach((el) => el.classList.remove("active"));
  const btn = document.querySelector(`.room-item[data-id="${roomID}"]`);
  if (btn) btn.classList.add("active");

  chatBox.innerHTML = `<div class="message-system">Joining #${esc(roomName)}…</div>`;

  sendJoin(roomID);

  await Promise.all([loadHistory(roomID), loadRoomProfile(roomID)]);
}

// API helpers
async function loadHistory(roomID) {
  try {
    const res = await fetch(`/rooms/${roomID}/messages`);
    if (!res.ok) return;
    const msgs = await res.json();
    if (roomID !== currentRoom) return;
    chatBox.innerHTML = "";
    if (!msgs.length) {
      chatBox.innerHTML = `<div class="message-system">No messages yet. Say hello!</div>`;
      return;
    }
    for (const m of [...msgs].reverse()) {
      let content = "[sent before you joined this room]";
      if (m.encrypted_key && myPrivateKey) {
        try {
          content = await decryptMessage(
            m.ciphertext,
            m.nonce,
            m.encrypted_key,
            myPrivateKey,
          );
        } catch (e) {
          console.error("decrypt failed", e);
          content = "[unable to decrypt]";
        }
      }
      appendHistoryMessage({ ...m, content });
    }
    chatBox.scrollTop = chatBox.scrollHeight;
  } catch (e) {
    console.error(e);
  }
}

let currentRoomProfile = null;

async function loadRoomProfile(roomID) {
  try {
    const res = await fetch(`/rooms/${roomID}/profile`);
    if (!res.ok) return;
    const profile = await res.json();

    if (roomID !== currentRoom) return;
    currentRoomProfile = profile;
    currentRoomProfile.members.forEach((m) => {
      if (m.online) onlineUsers[m.id] = m.username;
    });
    refreshMembersPanel();
  } catch (e) {
    console.error(e);
  }
}

function refreshMembersPanel() {
  const list = document.getElementById("members-list");
  if (!currentRoomProfile) {
    list.innerHTML = "";
    return;
  }

  const members = currentRoomProfile.members || [];
  document.getElementById("room-member-count").textContent =
    `${members.length} member${members.length !== 1 ? "s" : ""}`;

  list.innerHTML = members
    .map((m) => {
      const isOnline = !!onlineUsers[m.id];
      return `<div class="member-item">
      <span class="online-dot ${isOnline ? "online" : "offline"}"></span>
      <span>${esc(m.username)}</span>
      ${m.id === CURRENT_UID ? '<span class="badge bg-secondary ms-auto" style="font-size:.65rem">you</span>' : ""}
  </div>`;
    })
    .join("");
}

async function createRoom() {
  const name = document.getElementById("new-room-name").value.trim();
  const errEl = document.getElementById("create-room-error");
  errEl.classList.add("d-none");
  if (!name) {
    errEl.textContent = "Name required.";
    errEl.classList.remove("d-none");
    return;
  }

  try {
    const res = await fetch("/rooms", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": CSRF_TOKEN,
      },
      body: JSON.stringify({ name }),
    });
    const data = await res.json();
    if (!res.ok) {
      errEl.textContent = data.error;
      errEl.classList.remove("d-none");
      return;
    }

    bootstrap.Modal.getInstance(
      document.getElementById("createRoomModal"),
    ).hide();
    document.getElementById("new-room-name").value = "";
    addRoomToSidebar(data.id, data.name, false);
    await switchRoom(data.id, data.name);
  } catch (e) {
    errEl.textContent = "Network error.";
    errEl.classList.remove("d-none");
  }
}

async function joinByCode() {
  const raw = document.getElementById("invite-input").value.trim();
  if (!raw) return;
  const match = raw.match(/\/join\/([a-f0-9]+)/);
  const code = match ? match[1] : raw;
  try {
    const res = await fetch(`/join/${code}`, {
      headers: { "X-CSRF-Token": CSRF_TOKEN },
    });
    const data = await res.json();
    if (!res.ok) {
      alert(data.error || "Invalid code.");
      return;
    }
    document.getElementById("invite-input").value = "";
    addRoomToSidebar(data.id, data.name, false);
    await switchRoom(data.id, data.name);
  } catch {
    alert("Network error.");
  }
}
let reconnectAttempts = 0;

function scheduleReconnect() {
  reconnectAttempts++;
  const delay = Math.min(1000 * 2 ** reconnectAttempts, 15000);
  document.getElementById("conn-status").innerHTML =
    '<i class="bi bi-arrow-repeat text-warning"></i> Reconnecting…';

  setTimeout(async () => {
    try {
      const res = await fetch("/profile", { cache: "no-store" });
      if (res.status === 401) {
        window.location.href = "/login";
        return;
      }
    } catch {}
    connectWebSocket();
  }, delay);
}

function addRoomToSidebar(id, name, unread) {
  const list = document.getElementById("rooms-list");
  const existing = document.querySelector(`.room-item[data-id="${id}"]`);
  if (existing) {
    existing.querySelector(".unread-dot")?.remove();
    if (unread)
      existing.insertAdjacentHTML(
        "beforeend",
        '<span class="unread-dot ms-auto"></span>',
      );
    return;
  }
  const btn = document.createElement("button");
  btn.className =
    "list-group-item list-group-item-action room-item py-1 px-2 text-truncate d-flex align-items-center";
  btn.dataset.id = id;
  btn.innerHTML =
    `<i class="bi bi-hash me-1"></i>${esc(name)}` +
    (unread ? '<span class="unread-dot ms-auto"></span>' : "");
  btn.onclick = () => switchRoom(id, name);
  list.appendChild(btn);
}

async function loadRooms() {
  try {
    const res = await fetch("/rooms");
    if (!res.ok) return;
    const rooms = await res.json();
    rooms.forEach((r) => addRoomToSidebar(r.id, r.name, r.unread));
    if (rooms.length) await switchRoom(rooms[0].id, rooms[0].name);
  } catch (e) {
    console.error(e);
  }
}

// Profile modals
async function showProfile() {
  try {
    const res = await fetch("/profile");
    const data = await res.json();
    document.getElementById("up-initial").textContent =
      data.username[0].toUpperCase();
    document.getElementById("up-username").textContent = data.username;
    new bootstrap.Modal(document.getElementById("userProfileModal")).show();
  } catch {
    alert("Failed to load profile.");
  }
}

function showRoomProfile() {
  if (!currentRoomProfile) return;
  document.getElementById("rp-room-name").textContent = currentRoomProfile.name;
  document.getElementById("rp-invite-url").value =
    currentRoomProfile.invite_url;
  document.getElementById("rp-copy-ok").classList.add("d-none");

  const el = document.getElementById("rp-members-list");
  el.innerHTML = (currentRoomProfile.members || [])
    .map((m) => {
      const online = !!onlineUsers[m.id];
      return `<div class="member-item border-bottom">
      <span class="online-dot ${online ? "online" : "offline"}"></span>
      <span>${esc(m.username)}</span>
      ${m.id === CURRENT_UID ? '<span class="badge bg-secondary ms-auto" style="font-size:.65rem">you</span>' : ""}
  </div>`;
    })
    .join("");

  new bootstrap.Modal(document.getElementById("roomProfileModal")).show();
}

function copyRoomInvite() {
  const url = document.getElementById("rp-invite-url").value;
  navigator.clipboard.writeText(url).then(() => {
    const el = document.getElementById("rp-copy-ok");
    el.classList.remove("d-none");
    setTimeout(() => el.classList.add("d-none"), 2000);
  });
}

async function appendLiveMessage(msg) {
  const myKey = msg.keys && msg.keys[CURRENT_UID];
  let content = "[unable to decrypt]";
  if (myKey && myPrivateKey) {
    try {
      content = await decryptMessage(
        msg.ciphertext,
        msg.nonce,
        myKey,
        myPrivateKey,
      );
    } catch (e) {
      console.error("decrypt failed", e);
    }
  }
  const isSelf = msg.username === CURRENT_USER;
  chatBox.appendChild(
    bubble(isSelf ? "You" : msg.username, content, isSelf, fmtTime(msg.time)),
  );
  chatBox.scrollTop = chatBox.scrollHeight;
}
function appendSystemMessage(text) {
  const div = document.createElement("div");
  div.className = "message-system";
  div.textContent = text;
  chatBox.appendChild(div);
  chatBox.scrollTop = chatBox.scrollHeight;
}

function appendHistoryMessage(msg) {
  const isSelf = msg.sender_id === CURRENT_UID;
  chatBox.appendChild(
    bubble(
      isSelf ? "You" : msg.username,
      msg.content,
      isSelf,
      fmtTime(msg.time),
    ),
  );
}
function bubble(label, content, isSelf, time) {
  const row = document.createElement("div");
  row.className = `d-flex w-100 ${isSelf ? "justify-content-end" : "justify-content-start"}`;
  const timeHtml = time
    ? `<div class="small opacity-50 mt-1">${esc(time)}</div>`
    : "";
  if (isSelf) {
    row.innerHTML = `<div class="message-bubble bg-primary text-white shadow-sm">
      <div class="small fw-light opacity-75">You</div>
      <div>${esc(content)}</div>${timeHtml}</div>`;
  } else {
    row.innerHTML = `<div class="message-bubble bg-light border text-dark shadow-sm">
      <div class="small fw-bold text-primary">${esc(label)}</div>
      <div>${esc(content)}</div>${timeHtml}</div>`;
  }
  return row;
}

function esc(str) {
  if (!str) return "";
  return str.replace(
    /[&<>'"]/g,
    (t) =>
      ({
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        "'": "&#39;",
        '"': "&quot;",
      })[t] || t,
  );
}

window.onload = async () => {
  await loadMyPrivateKey();
  connectWebSocket();
  loadRooms();
};
let myPrivateKey = null;
const publicKeyCache = {};

async function loadMyPrivateKey() {
  const rawB64 = sessionStorage.getItem(sessionKeyFor(CURRENT_USER));
  if (!rawB64) {
    showUnlockPrompt();
    return;
  }
  try {
    myPrivateKey = await crypto.subtle.importKey(
      "pkcs8",
      b64decode(rawB64),
      { name: "RSA-OAEP", hash: "SHA-256" },
      false,
      ["decrypt"],
    );
  } catch {
    showUnlockPrompt();
  }
}

function showUnlockPrompt() {
  new bootstrap.Modal(document.getElementById("unlockModal")).show();
}

async function handleUnlock() {
  const password = document.getElementById("unlock-password").value;
  const errEl = document.getElementById("unlock-error");
  errEl.classList.add("d-none");

  const stored = loadWrappedKeypair(CURRENT_USER);
  if (!stored) {
    errEl.textContent =
      "No encryption key found for this account on this device.";
    errEl.classList.remove("d-none");
    return;
  }
  try {
    const { wrappingKey } = await deriveWrappingKey(
      CURRENT_USER,
      password,
      stored.salt,
    );
    const rawPriv = await unwrapPrivateKeyRaw(
      stored.wrappedPrivate,
      stored.iv,
      wrappingKey,
    );
    myPrivateKey = await crypto.subtle.importKey(
      "pkcs8",
      rawPriv,
      { name: "RSA-OAEP", hash: "SHA-256" },
      false,
      ["decrypt"],
    );
    sessionStorage.setItem(sessionKeyFor(CURRENT_USER), b64encode(rawPriv));
    bootstrap.Modal.getInstance(document.getElementById("unlockModal")).hide();
    if (currentRoom) await loadHistory(currentRoom);
  } catch {
    errEl.textContent =
      "Wrong password, or no key for this account on this device.";
    errEl.classList.remove("d-none");
  }
}

async function getMemberPublicKey(userID, pemString) {
  if (publicKeyCache[userID]) return publicKeyCache[userID];
  if (!pemString) return null;
  const key = await importPublicKeyFromPEM(pemString);
  publicKeyCache[userID] = key;
  return key;
}