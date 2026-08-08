function b64encode(buf) {
  return btoa(String.fromCharCode(...new Uint8Array(buf)));
}
function b64decode(b64) {
  const bin = atob(b64);
  const bytes = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);
  return bytes.buffer;
}
function pemToArrayBuffer(pem) {
  const b64 = pem.replace(/-----BEGIN [^-]+-----/, "")
                 .replace(/-----END [^-]+-----/, "")
                 .replace(/\s+/g, "");
  return b64decode(b64);
}

async function deriveWrappingKey(username, password, saltB64) {
  const enc = new TextEncoder();
  const salt = saltB64 ? new Uint8Array(b64decode(saltB64)) : crypto.getRandomValues(new Uint8Array(16));
  const baseKey = await crypto.subtle.importKey(
    "raw", enc.encode(username + ":" + password), "PBKDF2", false, ["deriveKey"]
  );
  const wrappingKey = await crypto.subtle.deriveKey(
    { name: "PBKDF2", salt, iterations: 210000, hash: "SHA-256" },
    baseKey, { name: "AES-GCM", length: 256 }, false, ["encrypt", "decrypt"]
  );
  return { wrappingKey, saltB64: b64encode(salt) };
}

async function wrapPrivateKey(privateKeyPKCS8Buf, wrappingKey) {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const wrapped = await crypto.subtle.encrypt({ name: "AES-GCM", iv }, wrappingKey, privateKeyPKCS8Buf);
  return { wrapped: b64encode(wrapped), iv: b64encode(iv) };
}

async function unwrapPrivateKeyRaw(wrappedB64, ivB64, wrappingKey) {
  const iv = new Uint8Array(b64decode(ivB64));
  return crypto.subtle.decrypt({ name: "AES-GCM", iv }, wrappingKey, b64decode(wrappedB64));
}

async function importPublicKeyFromPEM(pem) {
  return crypto.subtle.importKey(
    "spki", pemToArrayBuffer(pem), { name: "RSA-OAEP", hash: "SHA-256" }, true, ["encrypt"]
  );
}

async function encryptMessage(plaintext) {
  const aesKey = await crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, true, ["encrypt", "decrypt"]);
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const ciphertext = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv }, aesKey, new TextEncoder().encode(plaintext)
  );
  const rawKey = await crypto.subtle.exportKey("raw", aesKey);
  return { ciphertext: b64encode(ciphertext), nonce: b64encode(iv), rawKey };
}

async function encryptKeyForRecipient(rawAesKey, recipientPublicKey) {
  const encrypted = await crypto.subtle.encrypt({ name: "RSA-OAEP" }, recipientPublicKey, rawAesKey);
  return b64encode(encrypted);
}

async function decryptMessage(ciphertextB64, nonceB64, encryptedKeyB64, myPrivateKey) {
  const rawAesKey = await crypto.subtle.decrypt({ name: "RSA-OAEP" }, myPrivateKey, b64decode(encryptedKeyB64));
  const aesKey = await crypto.subtle.importKey("raw", rawAesKey, "AES-GCM", false, ["decrypt"]);
  const iv = new Uint8Array(b64decode(nonceB64));
  const plainBuf = await crypto.subtle.decrypt({ name: "AES-GCM", iv }, aesKey, b64decode(ciphertextB64));
  return new TextDecoder().decode(plainBuf);
}

function storageKeyFor(username) { return "e2ee_keypair_" + username; }

function sessionKeyFor(username) { return "e2ee_privkey_pkcs8:" + username; }

function saveWrappedKeypair(username, publicKeyPEM, wrappedPrivate, iv, saltB64) {
  localStorage.setItem(storageKeyFor(username), JSON.stringify({
    publicKeyPEM, wrappedPrivate, iv, salt: saltB64,
  }));
}

function loadWrappedKeypair(username) {
  const raw = localStorage.getItem(storageKeyFor(username));
  return raw ? JSON.parse(raw) : null;
}