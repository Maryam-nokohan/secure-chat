let generatedPublicPEM = null;
let generatedPrivateRaw = null;

async function generateKeyPair() {
    try {
        const keyPair = await window.crypto.subtle.generateKey(
            { name: "RSA-OAEP", modulusLength: 2048,
              publicExponent: new Uint8Array([1, 0, 1]),
              hash: "SHA-256" },
            true, ["encrypt", "decrypt"]
        );

        const pubRaw = await window.crypto.subtle.exportKey("spki", keyPair.publicKey);
        const pubB64 = b64encode(pubRaw);
        generatedPublicPEM = `-----BEGIN PUBLIC KEY-----\n${pubB64.match(/.{1,64}/g).join('\n')}\n-----END PUBLIC KEY-----`;
        generatedPrivateRaw = await window.crypto.subtle.exportKey("pkcs8", keyPair.privateKey);

        document.getElementById("public_key").value = generatedPublicPEM;
        document.getElementById("register-btn").disabled = false;
        document.getElementById("key-status").textContent = "Encryption keys ready.";
    } catch (e) {
        document.getElementById("key-status").textContent = "⚠ Key generation failed: " + e.message;
        document.getElementById("key-status").style.color = "red";
    }
}
generateKeyPair();

document.querySelector("form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const username = document.querySelector('input[name="username"]').value.trim();
    const password = document.querySelector('input[name="password"]').value;

    if (!generatedPrivateRaw) {
        alert("Encryption keys are still generating, please wait a moment.");
        return;
    }

    try {
        const { wrappingKey, saltB64 } = await deriveWrappingKey(username, password);
        const { wrapped, iv } = await wrapPrivateKey(generatedPrivateRaw, wrappingKey);
        saveWrappedKeypair(username, generatedPublicPEM, wrapped, iv, saltB64);
    } catch (err) {
        console.warn("Failed to save local encryption key (will still register):", err);
    }

    e.target.submit();
});