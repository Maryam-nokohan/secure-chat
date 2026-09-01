document.querySelector("form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const username = document.querySelector('input[name="username"]').value.trim();
    const password = document.querySelector('input[name="password"]').value;

    try {
        const stored = loadWrappedKeypair(username);
        if (stored) {

            const { wrappingKey } = await deriveWrappingKey(username, password, stored.salt);
            const rawPriv = await unwrapPrivateKeyRaw(stored.wrappedPrivate, stored.iv, wrappingKey);
            saveRawPrivateKey(username, b64encode(rawPriv));
            clearPendingPassword(username);
        } else {
            savePendingPassword(username, password);
        }
    } catch (err) {
        console.warn("Could not unlock local encryption key on login:", err);
        savePendingPassword(username, password);
    }

    e.target.submit();
});