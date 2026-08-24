const CSRF = window.CSRF_TOKEN;

async function api(path, opts = {}) {
  opts.headers = Object.assign({ "Content-Type": "application/json" }, opts.headers || {});
  if (opts.method && opts.method !== "GET") opts.headers["X-CSRF-Token"] = CSRF;
  const res = await fetch("/admin/api" + path, opts);
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || "request failed");
  return data;
}

function esc(str) {
  if (!str) return "";
  return String(str).replace(/[&<>'"]/g, (t) => ({ "&":"&amp;","<":"&lt;",">":"&gt;","'":"&#39;",'"':"&quot;" }[t] || t));
}

function formatUptime(sec) {
  const h = Math.floor(sec / 3600), m = Math.floor((sec % 3600) / 60), s = sec % 60;
  return `${h}h ${m}m ${s}s`;
}

async function loadStats() {
  try {
    const s = await api("/stats");
    document.getElementById("stat-users").textContent = s.total_users;
    document.getElementById("stat-online").textContent = s.online_users;
    document.getElementById("stat-rooms").textContent = s.total_rooms;
    document.getElementById("stat-messages").textContent = s.total_messages;
    document.getElementById("stat-uptime").textContent = formatUptime(s.uptime_seconds);
  } catch (e) { console.error(e); }
}

async function loadUsers() {
  try {
    const users = await api("/users?limit=100");
    document.getElementById("users-body").innerHTML = users.map(u => `
      <tr>
        <td><span class="online-dot ${u.online ? "online" : "offline"}"></span></td>
        <td>${esc(u.username)}</td>
        <td><span class="badge ${u.role === "admin" ? "bg-danger" : "bg-secondary"}">${esc(u.role)}</span></td>
        <td class="small text-muted">${new Date(u.created_at).toLocaleString()}</td>
        <td class="text-end">
          <button class="btn btn-sm btn-outline-secondary" onclick="editUser('${u.id}','${esc(u.username)}')"><i class="bi bi-pencil"></i></button>
          ${u.role === "admin"
            ? `<button class="btn btn-sm btn-outline-warning" onclick="setRole('${u.id}','user')">Demote</button>`
            : `<button class="btn btn-sm btn-outline-primary" onclick="setRole('${u.id}','admin')">Promote</button>`}
          <button class="btn btn-sm btn-outline-danger" onclick="deleteUser('${u.id}','${esc(u.username)}')"><i class="bi bi-trash"></i></button>
        </td>
      </tr>`).join("");
  } catch (e) { console.error(e); }
}

async function deleteUser(id, name) {
  if (!confirm(`Delete user "${name}"? This can be undone from the Actions panel.`)) return;
  try { await api(`/users/${id}`, { method: "DELETE" }); await Promise.all([loadUsers(), loadHistory()]); }
  catch (e) { alert(e.message); }
}

async function setRole(id, role) {
  try { await api(`/users/${id}/role`, { method: "POST", body: JSON.stringify({ role }) }); await Promise.all([loadUsers(), loadHistory()]); }
  catch (e) { alert(e.message); }
}

function editUser(id, currentUsername) {
  document.getElementById("edit-user-id").value = id;
  document.getElementById("edit-username").value = currentUsername;
  document.getElementById("edit-bio").value = "";
  new bootstrap.Modal(document.getElementById("editUserModal")).show();
}

async function submitEditUser() {
  const id = document.getElementById("edit-user-id").value;
  const username = document.getElementById("edit-username").value.trim();
  const bio = document.getElementById("edit-bio").value;
  try {
    await api(`/users/${id}`, { method: "PUT", body: JSON.stringify({ username, bio }) });
    bootstrap.Modal.getInstance(document.getElementById("editUserModal")).hide();
    await Promise.all([loadUsers(), loadHistory()]);
  } catch (e) { alert(e.message); }
}

async function loadHistory() {
  try {
    const history = await api("/history");
    document.getElementById("history-list").innerHTML = history.length
      ? history.slice().reverse().map(h => `<div class="small border-bottom py-1">${esc(h.description)}</div>`).join("")
      : `<div class="text-muted small">No actions yet.</div>`;
  } catch (e) { console.error(e); }
}

async function undoLast() {
  try {
    const res = await api("/undo", { method: "POST" });
    alert("Undone: " + res.undone);
    await Promise.all([loadUsers(), loadHistory()]);
  } catch (e) { alert(e.message); }
}

async function loadLogs() {
  try {
    const lines = await api("/logs?limit=200");
    const el = document.getElementById("logs-box");
    const atBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 10;
    el.textContent = lines.join("");
    if (atBottom) el.scrollTop = el.scrollHeight;
  } catch (e) { console.error(e); }
}

let newAdminPublicPEM = null;

async function generateAdminKeyPair() {
  const keyPair = await crypto.subtle.generateKey(
    { name: "RSA-OAEP", modulusLength: 2048, publicExponent: new Uint8Array([1,0,1]), hash: "SHA-256" },
    true, ["encrypt", "decrypt"]
  );
  const pubRaw = await crypto.subtle.exportKey("spki", keyPair.publicKey);
  const pubB64 = btoa(String.fromCharCode(...new Uint8Array(pubRaw)));
  newAdminPublicPEM = `-----BEGIN PUBLIC KEY-----\n${pubB64.match(/.{1,64}/g).join('\n')}\n-----END PUBLIC KEY-----`;
  document.getElementById("new-admin-key-status").textContent = "Encryption keys ready.";
}

function openCreateAdminModal() {
  document.getElementById("new-admin-username").value = "";
  document.getElementById("new-admin-password").value = "";
  document.getElementById("new-admin-key-status").textContent = "Generating encryption keys…";
  generateAdminKeyPair();
  new bootstrap.Modal(document.getElementById("createAdminModal")).show();
}

async function submitCreateAdmin() {
  const username = document.getElementById("new-admin-username").value.trim();
  const password = document.getElementById("new-admin-password").value;
  if (!newAdminPublicPEM) { alert("Encryption keys still generating, wait a moment."); return; }
  try {
    await api("/users", { method: "POST", body: JSON.stringify({ username, password, public_key: newAdminPublicPEM }) });
    bootstrap.Modal.getInstance(document.getElementById("createAdminModal")).hide();
    await loadUsers();
  } catch (e) { alert(e.message); }
}

window.onload = () => {
  loadStats(); loadUsers(); loadHistory(); loadLogs();
  setInterval(loadStats, 5000);
  setInterval(loadLogs, 4000);
  setInterval(loadUsers, 15000);
};