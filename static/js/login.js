document.querySelector("form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const username = document.querySelector('input[name="username"]').value.trim();
    const password = document.querySelector('input[name="password"]').value;

    sessionStorage.removeItem(sessionKeyFor(username));

    const stored = loadWrappedKeypair(username);
    if (stored) {
        try {
            const { wrappingKey } = await deriveWrappingKey(username, password, stored.salt);
            const rawPriv = await unwrapPrivateKeyRaw(stored.wrappedPrivate, stored.iv, wrappingKey);
            sessionStorage.setItem(sessionKeyFor(username), b64encode(rawPriv));
        } catch (err) {
            console.warn("Could not unlock local encryption key on login:", err);
            sessionStorage.removeItem(sessionKeyFor(username));
        }
    }

    e.target.submit();
});