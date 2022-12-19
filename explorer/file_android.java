package org.gioui.x.explorer;

import java.io.InputStream;
import java.io.OutputStream;
import java.io.Closeable;
import java.io.Flushable;

public class file_android {
    public String err;
    public Object handler;

    public void setHandle(Object f) {
        this.handler = f;
    }

    public static int fileRead(file_android fa, byte[] b) {
        try {
            return ((InputStream) fa.handler).read(b, 0, b.length);
        } catch (Exception e) {
            fa.err = e.toString();
            return 0;
        }
    }

    public static boolean fileWrite(file_android fa, byte[] b) {
        try {
            ((OutputStream) fa.handler).write(b);
            return true;
        } catch (Exception e) {
            fa.err = e.toString();
            return false;
        }
    }

    public static boolean fileClose(file_android fa) {
        try {
            if (fa.handler instanceof Flushable) {
                ((Flushable) fa.handler).flush();
            }
            if (fa.handler instanceof Closeable) {
                ((Closeable) fa.handler).close();
            }
            return true;
        } catch (Exception e) {
            fa.err = e.toString();
            return false;
        }
    }

    public static String getError(file_android fa) {
        return fa.err;
    }
}