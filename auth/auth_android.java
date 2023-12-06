package com.inkeliz.auth_android;

import android.util.Log;
import android.app.Activity;
import android.view.View;
import android.content.Context;
import android.content.Intent;
import android.app.Fragment;
import android.app.FragmentManager;
import android.app.FragmentTransaction;
import android.app.PendingIntent;
import java.util.ArrayList;
import java.util.List;
import com.google.android.gms.auth.api.identity.*;
import com.google.android.gms.auth.api.signin.*;
import com.google.android.gms.common.api.ApiException;
import com.google.android.gms.common.api.Scope;
import com.google.android.gms.tasks.Task;
import com.google.android.gms.tasks.OnCompleteListener;
import android.content.pm.PackageManager;
import android.content.pm.ResolveInfo;
import android.net.Uri;
import android.os.Bundle;
import androidx.browser.customtabs.CustomTabsIntent;

public class auth_android {
    final googleauth_android_fragment frag = new googleauth_android_fragment();

    // Functions defined on Golang.
    static public native void NativeAuthCallback(long handler, String idToken);

    // Request code for Google Sign In.
    static private final int REQUEST_CODE = 1_820_989;
    static private long LastHandler;

    private SignInClient oneTapClient;

    public static class googleauth_android_fragment extends Fragment {
        Context context;
        SignInClient oneTapClient;

        @Override public void onAttach(Context ctx) {
            this.context = ctx;
            super.onAttach(ctx);
        }

        @Override public void onActivityResult(int requestCode, int resultCode, Intent data) {
            super.onActivityResult(requestCode, resultCode, data);

            Activity activity = this.getActivity();
            activity.runOnUiThread(new Runnable() {
                public void run() {
                    if(requestCode != REQUEST_CODE){
                        return;
                    }

                    try {
                        SignInCredential acc = oneTapClient.getSignInCredentialFromIntent(data);
                        String idToken = acc.getGoogleIdToken();

                        NativeAuthCallback(LastHandler, idToken);
                    } catch (Exception e) {
                        Log.wtf("auth_android", e);
                    }
                }
            });

        }
    }

    public void openGeneral(View view, String url, long handler) {
        CustomTabsIntent.Builder builder = new CustomTabsIntent.Builder();
        CustomTabsIntent customTabsIntent = builder.build();

        PackageManager packageManager = view.getContext().getPackageManager();
        List<ResolveInfo> resolvedActivityList = packageManager.queryIntentActivities(customTabsIntent.intent, PackageManager.MATCH_DEFAULT_ONLY);

        for (ResolveInfo info : resolvedActivityList) {
            if (info.activityInfo.packageName.toLowerCase().contains("com.android.chrome")) {
                customTabsIntent.intent.setPackage("com.android.chrome");
                break;
            }
        }

        ((Activity) view.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                customTabsIntent.launchUrl((Activity) view.getContext(), Uri.parse(url));
            }
        });
    }

    public void openNative(View view, String clientID, String nonce, long handler) {
        LastHandler = handler;

        if (view == null) {
            return;
        }

        ((Activity) view.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                Context context = view.getContext();
                registerFrag(context, view);

                oneTapClient.signOut().addOnCompleteListener((Activity) view.getContext(), new OnCompleteListener<Void>() {
                    @Override
                    public void onComplete(Task<Void> task) {
                        GetSignInIntentRequest request = GetSignInIntentRequest.builder().setServerClientId(clientID).setNonce(nonce).build();

                        oneTapClient.getSignInIntent(request).addOnCompleteListener((Activity) view.getContext(), new OnCompleteListener<PendingIntent>() {
                            @Override
                            public void onComplete(Task<PendingIntent> task) {
                                try {
                                    PendingIntent res = task.getResult(ApiException.class);
                                    try {
                                        frag.startIntentSenderForResult(res.getIntentSender(), REQUEST_CODE, null, 0, 0, 0, null);
                                    } catch (Exception e) {
                                        Log.wtf("auth_android", e);
                                    }
                                } catch (ApiException e) {
                                    Log.wtf("auth_android", e);
                                }
                            }
                        });
                    }
                });
            }
        });
    }

    private void registerFrag(Context context, View view) {
        final Context ctx = view.getContext();
        final FragmentManager fm;

        try {
            fm = (FragmentManager) ctx.getClass().getMethod("getFragmentManager").invoke(ctx);
        } catch (Exception e) {
            e.printStackTrace();
            return;
        }

        if (fm.findFragmentByTag("googleauth_android_fragment") != null) {
            return; // Already exists;
        }

        oneTapClient = Identity.getSignInClient(context);
        frag.oneTapClient = oneTapClient;

        FragmentTransaction ft = fm.beginTransaction();
        ft.add(frag, "googleauth_android_fragment");
        ft.commitNow();
    }
}