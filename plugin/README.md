Plugin
-------

This package is required to interact with any other package in this repository.

You MUST use `gioplugins.Event()` in your event-loop:

```diff
for { // Gio main event loop
+    evt := gioplugins.Event(window)

    switch evt := evt.(type) {
        // ...
    }
}
```