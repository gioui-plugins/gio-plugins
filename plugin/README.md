Plugin
-------

This package is required to interact with any other package in this repository.

You MUST use `gioplugins.Hijack()` in your event-loop:

```diff
for { // Gio main event loop
+    evt := gioplugins.Hijack(window)

    switch evt := evt.(type) {
        // ...
    }
}
```