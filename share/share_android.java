package com.inkeliz.share_android;

import android.app.Activity;
import android.view.View;
import android.content.Intent;

public class share_android {
    static public void openShare(View view, Intent i) {
        ((Activity) view.getContext()).startActivity(Intent.createChooser(i, null));
    }
    static public void shareText(View view, String title, String text) {
        Intent i = new Intent(Intent.ACTION_SEND);
        i.setType("text/plain");
        i.putExtra(Intent.EXTRA_TEXT, text);

        openShare(view, i);
    }
    static public void shareWebsite(View view, String title, String text, String link) {
        Intent i = new Intent(Intent.ACTION_SEND);
        i.setType("text/plain");
        i.putExtra(Intent.EXTRA_TITLE, text);
        i.putExtra(Intent.EXTRA_TEXT, link);

        openShare(view, i);
    }
}