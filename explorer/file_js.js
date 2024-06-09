(() => {
    if (go.importObject.gojs == null) {
        go.importObject.gojs = {};
    }
})();
(() => {
    Object.assign(go.importObject.gojs, {
        // fileSlice(uint32, uint32, js.Value, js.Func, js.Func) uint32
        "github.com/gioui-plugins/gio-plugins/explorer.fileSlice": (sp) => {
            sp = (sp >>> 0);
            // uint32:
            let _start = go.mem.getUint32(sp + 8, true);
            // uint32:
            let _end = go.mem.getUint32(sp + 8 + 4, true);
            // js.Value:
            let _refBuffer = go.mem.getUint32(sp + 8 + 8, true);
            // js.Func:
            let _refSuccess = go.mem.getUint32(sp + 8 + 8 + 8 + 8, true);
            // js.func:
            let _refFailure = go.mem.getUint32(sp + 8 + 8 + 8 + 8 + 8 + 8 + 8, true);

            go._values[_refBuffer].slice(_start, _end).arrayBuffer().then(go._values[_refSuccess], go._values[_refFailure])
        },
        // fileRead(js.Value, []byte) uint32
        "github.com/gioui-plugins/gio-plugins/explorer.fileRead": (sp) => {
            sp = (sp >>> 0);
            // js.Value:
            let _ref = go.mem.getUint32(sp + 8, true);
            // []byte:
            let _slicePointer = go.mem.getUint32(sp + 8 + 8 + 8, true) + go.mem.getInt32(sp + 8 + 8 + 8 + 4, true) * 4294967296;
            //let _sliceLength = go.mem.getUint32(sp + 8 + 8 + 8 + 8, true) + go.mem.getInt32(sp + 8 + 8 + 8 + 8 + 4, true) * 4294967296;

            let subArray = new Uint8Array(go._values[_ref]);
            for (let i = 0; i < subArray.length; i++) {
                go.mem.setUint8(_slicePointer + i, subArray[i]);
            }

            // output:
            go.mem.setUint32(sp + 8 + 8 + 8 + 8 + 8 + 8, subArray.length, true)
        },
        // fileWrite(js.Value, []byte)
        "github.com/gioui-plugins/gio-plugins/explorer.fileWrite": (sp) => {
            sp = (sp >>> 0);
            // js.Value:
            let _ref = go.mem.getUint32(sp + 8, true);
            // []byte:
            let _slicePointer = go.mem.getUint32(sp + 8 + 8 + 8, true) + go.mem.getInt32(sp + 8 + 8 + 8 + 4, true) * 4294967296;
            let _sliceLength = go.mem.getUint32(sp + 8 + 8 + 8 + 8, true) + go.mem.getInt32(sp + 8 + 8 + 8 + 8 + 4, true) * 4294967296;

            let jsArray = go._values[_ref];
            let goSlice = new Uint8Array(go._inst.exports.mem.buffer, _slicePointer, _sliceLength);

            let newArray = new Uint8Array(jsArray.length + _sliceLength);
            newArray.set(jsArray);
            newArray.set(goSlice, jsArray.length);
            go._values[_ref] = newArray;
        },
        // writableWrite(js.Value, js.Value, []byte)
        "github.com/gioui-plugins/gio-plugins/explorer.writableWrite": (sp) => {
            sp = (sp >>> 0);
            // js.Value:
            sp += 8
            let _refWritable = go.mem.getUint32(sp, true);
            sp += 8 + 8
            let _refSuccess = go.mem.getUint32(sp, true);
            sp += 8 + 8
            let _refFailure = go.mem.getUint32(sp, true);
            sp += 8 + 8
            // []byte:
            let _slicePointer = go.mem.getUint32(sp, true) + go.mem.getInt32(sp + 4, true) * 4294967296;
            sp += 8
            let _sliceLength = go.mem.getUint32(sp, true) + go.mem.getInt32(sp + 4, true) * 4294967296;

            go._values[_refWritable].write({
                type: "write",
                data: new Uint8Array(go._inst.exports.mem.buffer, _slicePointer, _sliceLength),
                size: _sliceLength,
            }).then(go._values[_refSuccess]).catch(go._values[_refFailure]);
        }
    });
})();