package com.inkeliz.hyperlink_android;

import android.app.Activity;
import android.view.View;
import android.content.Context;
import android.content.Intent;
import android.net.Uri;

public class hyperlink_android {

    public static void open(View view, String url) {
        Intent intent = new Intent(Intent.ACTION_VIEW, Uri.parse(url));
        Activity activity = (Activity)view.getContext();

        activity.runOnUiThread(new Runnable() {
            public void run() {
                activity.startActivity(intent);
            }
        });

    }

}