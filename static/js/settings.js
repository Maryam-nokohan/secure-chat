function esc(str) {
  if (!str) return "";
  return str.replace(/[&<>'"]/g, (t) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", "'": "&#39;", '"': "&quot;" }[t] || t));
}

async function loadProfile() {
  try {
    const res = await fetch("/profile");
    if (!res.ok) return;
    const p = await res.json();
    document.getElementById("username-input").value = p.username || "";
    document.getElementById("bio-input").value = p.bio || "";
    document.getElementById("public-id").value = p.public_id || "";
    document.getElementById("avatar-initial").textContent = (p.username || "?")[0].toUpperCase();
    if (p.avatar_url) {
      const img = document.getElementById("avatar-preview");
      img.src = p.avatar_url + "?t=" + Date.now();
      img.classList.remove("d-none");
      document.getElementById("avatar-fallback").classList.add("d-none");
    }
  } catch (e) {
    console.error(e);
  }
}

function copyPublicId() {
  const val = document.getElementById("public-id").value;
  navigator.clipboard.writeText(val).then(() => {
    const el = document.getElementById("copy-ok");
    el.classList.remove("d-none");
    setTimeout(() => el.classList.add("d-none"), 1500);
  });
}

async function saveProfile() {
  const errEl = document.getElementById("profile-error");
  const okEl = document.getElementById("profile-ok");
  errEl.classList.add("d-none");
  okEl.classList.add("d-none");
  const username = document.getElementById("username-input").value.trim();
  const bio = document.getElementById("bio-input").value;
  try {
    const res = await fetch("/settings/profile", {
      method: "PUT",
      headers: { "Content-Type": "application/json", "X-CSRF-Token": CSRF_TOKEN },
      body: JSON.stringify({ username, bio }),
    });
    const data = await res.json();
    if (!res.ok) {
      errEl.textContent = data.error || "Failed to save.";
      errEl.classList.remove("d-none");
      return;
    }
    okEl.classList.remove("d-none");
  } catch {
    errEl.textContent = "Network error.";
    errEl.classList.remove("d-none");
  }
}

document.getElementById("avatar-input").addEventListener("change", async (e) => {
  const file = e.target.files[0];
  if (!file) return;
  const errEl = document.getElementById("avatar-error");
  errEl.classList.add("d-none");

  const form = new FormData();
  form.append("avatar", file);
  try {
    const res = await fetch("/settings/avatar", {
      method: "POST",
      headers: { "X-CSRF-Token": CSRF_TOKEN },
      body: form,
    });
    const data = await res.json();
    if (!res.ok) {
      errEl.textContent = data.error || "Upload failed.";
      errEl.classList.remove("d-none");
      return;
    }
    const img = document.getElementById("avatar-preview");
    img.src = data.avatar_url + "?t=" + Date.now();
    img.classList.remove("d-none");
    document.getElementById("avatar-fallback").classList.add("d-none");
  } catch {
    errEl.textContent = "Network error.";
    errEl.classList.remove("d-none");
  }
});

window.onload = loadProfile;