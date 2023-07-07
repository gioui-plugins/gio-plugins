// SPDX-License-Identifier: Unlicense OR MIT

package com.inkeliz.deeplink;

import android.app.Activity;
import android.os.Bundle;
import android.content.Intent;
import android.content.pm.PackageManager;

public class deeplink_android extends Activity {

    // This function is defined in Golang.
    static public native void ReceiveScheme(String scheme);

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        this.onNewIntent(this.getIntent());
    }

    @Override
    protected void onNewIntent(Intent intent) {
        super.onNewIntent(intent);

        // That must come first to initialize the Golang runtime, otherwise
        // ReceiveScheme will not be called.
        try {
            Intent intent = new Intent(this, Class.forName(this.getMainActivityName()));
            startActivity(intent);
        } catch (Exception e) {
            e.printStackTrace();
        }

        if(intent.getData() != null) {
            ReceiveScheme(intent.getData().toString());
        }
    }

    // That will be org.gioui.GioActivity in Gio, but we can't use it here
    // otherwise will not work on non-gio apps.
    private String getMainActivityName() {
        String mainActivityName = "";

        try {
            PackageManager packageManager = getPackageManager();
            Intent intent = new Intent(Intent.ACTION_MAIN);
            intent.addCategory(Intent.CATEGORY_LAUNCHER);
            intent.setPackage(getPackageName());

            mainActivityName = packageManager.resolveActivity(intent, 0).activityInfo.name;
        } catch (Exception e) {
            e.printStackTrace();
        }

        return mainActivityName;
    }
}


