self.addEventListener('push', (event) => {
    if (!event.data) return;

    try {
        const payload = event.data.json();

        if (payload.notification) {
            const data = payload.notification;

            const options = {
                body: data.body || "",
                lang: data.lang || "en",
                dir: data.dir || "ltr",
                silent: data.silent || false,
                data: {
                    url: data.navigate
                },
                icon: data.icon || '',
                badge: data.badge || ''
            };

            event.waitUntil(
                self.registration.showNotification(data.title, options)
            );
        }
    } catch (e) {
        console.error("Push payload was not JSON or failed to parse:", e);
    }
});

self.addEventListener('notificationclick', (event) => {
    const clickedNotification = event.notification;
    clickedNotification.close();

    const urlToOpen = clickedNotification.data?.url;

    if (urlToOpen) {
        event.waitUntil(
            clients.matchAll({ type: 'window' }).then((windowClients) => {
                for (let client of windowClients) {
                    if (client.url === urlToOpen && 'focus' in client) {
                        return client.focus();
                    }
                }

                if (clients.openWindow) {
                    return clients.openWindow(urlToOpen);
                }
            })
        );
    }
});