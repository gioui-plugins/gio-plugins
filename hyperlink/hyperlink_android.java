package com.inkeliz.hyperlink_android;

import android.app.Activity;
import android.view.View;
import android.content.Context;
import android.content.Intent;
import android.net.Uri;
import android.content.pm.PackageManager;
import android.content.pm.PackageInfo;

public class hyperlink_android {

    public static void open(View view, String url, String packageName) {
        Intent intent = new Intent(Intent.ACTION_VIEW, Uri.parse(url));
        Activity activity = (Activity)view.getContext();

        if (packageName.length() > 0 && isInstalled(activity, packageName)) {
            intent.setPackage(packageName);
        }

        activity.runOnUiThread(new Runnable() {
            public void run() {
                activity.startActivity(intent);
            }
        });

    }

    private static boolean isInstalled(Activity activity, String packageName) {
        PackageManager packageManager = activity.getApplicationContext().getPackageManager();
        Intent intentForCheck = new Intent(Intent.ACTION_VIEW);
        if (intentForCheck.resolveActivity(packageManager) != null) {
            try {
                packageManager.getPackageInfo(packageName, PackageManager.GET_ACTIVITIES);
                return true;
            } catch (Exception e) {}
        }
        return false;
    }

}