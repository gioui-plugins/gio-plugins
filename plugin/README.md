Plugin
-------

This package is required to interact with any other package in this repository.

You MUST use `plugin.Install()` in your event-loop:

```diff
for evt := range w.Events() { // Gio main event loop
+    plugin.Install(w, evt)

    switch evt := evt.(type) {
        // ...
    }
}
```