let generatedPublicPEM = null;
let generatedPrivateRaw = null;

async function generateKeyPair() {
  try {
    const keyPair = await window.crypto.subtle.generateKey(
      { name: "RSA-OAEP", modulusLength: 2048, publicExponent: new Uint8Array([1, 0, 1]), hash: "SHA-256" },
      true,
      ["encrypt", "decrypt"],
    );
    const pubRaw = await window.crypto.subtle.exportKey("spki", keyPair.publicKey);
    const pubB64 = b64encode(pubRaw);
    generatedPublicPEM = `-----BEGIN PUBLIC KEY-----\n${pubB64.match(/.{1,64}/g).join("\n")}\n-----END PUBLIC KEY-----`;
    generatedPrivateRaw = await window.crypto.subtle.exportKey("pkcs8", keyPair.privateKey);

    document.getElementById("public_key").value = generatedPublicPEM;
    document.getElementById("submit-btn").disabled = false;
    document.getElementById("key-status").textContent = "Ready.";
  } catch (e) {
    document.getElementById("key-status").textContent = "Key generation failed: " + e.message;
  }
}
generateKeyPair();

document.querySelector("form").addEventListener("submit", async function (e) {
  e.preventDefault();

  const username = window.CURRENT_USER;
  const password = document.getElementById("password").value;
  if (!generatedPrivateRaw) {
    alert("Encryption keys are still generating, please wait a moment.");
    return;
  }
  if (!password) {
    alert("Please enter a password.");
    return;
  }

  try {
    const { wrappingKey, saltB64 } = await deriveWrappingKey(username, password);
    const { wrapped, iv } = await wrapPrivateKey(generatedPrivateRaw, wrappingKey);

    document.getElementById("wrapped_private_key").value = wrapped;
    document.getElementById("private_key_iv").value = iv;
    document.getElementById("private_key_salt").value = saltB64;

    saveWrappedKeypair(username, generatedPublicPEM, wrapped, iv, saltB64);
    saveRawPrivateKey(username, b64encode(generatedPrivateRaw));

    if (!loadWrappedKeypair(username)) {
      throw new Error("key not persisted locally");
    }
  } catch (err) {
    console.error("Failed to prepare encryption key:", err);
    alert("Your encryption key could not be created on this device. Make sure this browser allows local storage and try again.");
    return;
  }

  e.target.submit();
});