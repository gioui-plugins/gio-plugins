package com.inkeliz.pushnotification_android;

import android.Manifest;
import android.app.Activity;
import android.content.pm.PackageManager;
import android.os.Build;
import android.util.Log;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;
import android.view.View;

import com.google.android.gms.tasks.OnCompleteListener;
import com.google.android.gms.tasks.Task;
import com.google.firebase.FirebaseApp;
import com.google.firebase.FirebaseOptions;
import com.google.firebase.messaging.FirebaseMessaging;

public class pushnotification_android {
    private static final String TAG = "GioPushNotification";

    public pushnotification_android() {
    }

    public void initialize(View view, String appId, String projectId, String apiKey, String senderId) {
        Activity activity = (Activity) view.getContext();

        if (FirebaseApp.getApps(activity).isEmpty()) {
            FirebaseOptions options = new FirebaseOptions.Builder()
                    .setApplicationId(appId)
                    .setProjectId(projectId)
                    .setApiKey(apiKey)
                    .setGcmSenderId(senderId)
                    .build();
            FirebaseApp.initializeApp(activity, options);
        }
    }

    public void requestPermission(Activity activity) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            if (ContextCompat.checkSelfPermission(activity,
                    Manifest.permission.POST_NOTIFICATIONS) != PackageManager.PERMISSION_GRANTED) {
                ActivityCompat.requestPermissions(activity, new String[] { Manifest.permission.POST_NOTIFICATIONS },
                        1001);
            }
        }
    }

    public void getToken(View view, long handler) {
        Activity activity = (Activity) view.getContext();

        requestPermission(activity); // Ensure permission is requested

        FirebaseMessaging.getInstance().getToken()
                .addOnCompleteListener(new OnCompleteListener<String>() {
                    @Override
                    public void onComplete(Task<String> task) {
                        if (!task.isSuccessful()) {
                            Log.w(TAG, "Fetching FCM registration token failed", task.getException());
                            onError(handler,
                                    task.getException() != null ? task.getException().getMessage() : "Unknown error");
                            return;
                        }

                        // Get new FCM registration token
                        String token = task.getResult();
                        onTokenReceived(handler, token);
                    }
                });
    }

    private static native void onTokenReceived(long handler, String token);

    private static native void onError(long handler, String message);
}
