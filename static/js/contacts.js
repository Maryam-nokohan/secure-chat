async function loadContacts() {
  try {
    const res = await fetch("/contacts");
    if (!res.ok) return;
    renderContacts(await res.json());
  } catch (e) {
    console.error(e);
  }
}

function renderContacts(contacts) {
  const list = document.getElementById("contacts-list");
  const pending = document.getElementById("pending-requests");
  if (!list || !pending) return;
  list.innerHTML = "";
  pending.innerHTML = "";

  contacts.forEach((ct) => {
    if (ct.status === "accepted" && ct.room_id) {
      const btn = document.createElement("button");
      btn.className = "list-group-item list-group-item-action room-item py-1 px-2 text-truncate";
      btn.innerHTML = `<i class="bi bi-person-circle me-1"></i>${esc(ct.username)}`;
      btn.onclick = () => switchRoom(ct.room_id, ct.username);
      list.appendChild(btn);
    } else if (ct.status === "pending" && ct.incoming) {
      const row = document.createElement("div");
      row.className = "d-flex align-items-center justify-content-between small border-bottom py-1 px-1";
      row.innerHTML = `<span>${esc(ct.username)}</span>
        <span>
          <button class="btn btn-sm btn-outline-success py-0 px-1" onclick="respondContact('${ct.id}', true)"><i class="bi bi-check"></i></button>
          <button class="btn btn-sm btn-outline-danger py-0 px-1" onclick="respondContact('${ct.id}', false)"><i class="bi bi-x"></i></button>
        </span>`;
      pending.appendChild(row);
    }
  });
}

async function sendContactRequest() {
  const input = document.getElementById("add-contact-id");
  const errEl = document.getElementById("add-contact-error");
  const publicId = input.value.trim();
  errEl.classList.add("d-none");
  if (!publicId) return;
  try {
    const res = await fetch("/contacts", {
      method: "POST",
      headers: { "Content-Type": "application/json", "X-CSRF-Token": CSRF_TOKEN },
      body: JSON.stringify({ public_id: publicId }),
    });
    const data = await res.json();
    if (!res.ok) {
      errEl.textContent = data.error || "Failed to send request.";
      errEl.classList.remove("d-none");
      return;
    }
    input.value = "";
    await loadContacts();
  } catch {
    errEl.textContent = "Network error.";
    errEl.classList.remove("d-none");
  }
}

async function respondContact(id, accept) {
  try {
    await fetch(`/contacts/${id}/${accept ? "accept" : "decline"}`, {
      method: "POST",
      headers: { "X-CSRF-Token": CSRF_TOKEN },
    });
    await Promise.all([loadContacts(), loadRooms()]);
  } catch (e) {
    console.error(e);
  }
}