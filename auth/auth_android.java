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

import android.os.CancellationSignal;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Executors;

import androidx.credentials.Credential;
import androidx.credentials.CredentialManager;
import androidx.credentials.CredentialManagerCallback;
import androidx.credentials.CustomCredential;
import androidx.credentials.GetCredentialRequest;
import androidx.credentials.GetCredentialResponse;
import androidx.credentials.PrepareGetCredentialResponse;
import androidx.credentials.exceptions.GetCredentialException;

import com.google.android.libraries.identity.googleid.GetGoogleIdOption;
import com.google.android.libraries.identity.googleid.GoogleIdTokenCredential;
import com.google.android.libraries.identity.googleid.GetSignInWithGoogleOption;

import android.content.pm.PackageManager;
import android.content.pm.ResolveInfo;
import android.net.Uri;
import android.os.Bundle;
import androidx.browser.customtabs.CustomTabsIntent;
import androidx.browser.customtabs.CustomTabsServiceConnection;
import androidx.browser.customtabs.CustomTabsClient;
import android.content.ComponentName;

public class auth_android {
    static private long LastHandler = 0;

    static private CustomTabsServiceConnection connection = null;
    static private CustomTabsClient client = null;

    // Functions defined on Golang.
    static public native void NativeAuthCallback(long handler, String idToken);

    public void openGeneral(View view, String url, long handler) {
        if (connection == null) {
            connection = new CustomTabsServiceConnection() {
                @Override
                public void onCustomTabsServiceConnected(ComponentName componentName, CustomTabsClient customTabsClient) {
                    client = customTabsClient;
                }

                @Override
                public void onServiceDisconnected(ComponentName name) {
                    client = null;
                }
            };

            CustomTabsClient.bindCustomTabsService(view.getContext().getApplicationContext(), "com.android.chrome", connection);
        }

        CustomTabsIntent.Builder builder = new CustomTabsIntent.Builder();
        CustomTabsIntent customTabsIntent = builder.build();

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

        CredentialManager credentialManager = CredentialManager.create(view.getContext());
        CancellationSignal cancellationSignal = new CancellationSignal();

        GetSignInWithGoogleOption googleIdOption = new GetSignInWithGoogleOption.Builder(clientID).setNonce(nonce).build();

        GetCredentialRequest request = new GetCredentialRequest.Builder().addCredentialOption(googleIdOption).build();

        credentialManager.getCredentialAsync(
          ((Activity) view.getContext()),
          request,
          cancellationSignal,
          Executors.newSingleThreadExecutor(),
          new CredentialManagerCallback<GetCredentialResponse, GetCredentialException>() {
            @Override
            public void onResult(GetCredentialResponse result) {
              handleSignIn(result);
            }

            @Override
            public void onError(GetCredentialException e) {
              handleFailure(e);
            }
          }
        );
    }

    public void handleSignIn(GetCredentialResponse result) {
      // Handle the successfully returned credential.
      Credential credential = result.getCredential();

      if (credential instanceof CustomCredential) {
        if (GoogleIdTokenCredential.TYPE_GOOGLE_ID_TOKEN_CREDENTIAL.equals(credential.getType())) {
          try {
            GoogleIdTokenCredential googleIdTokenCredential = GoogleIdTokenCredential.createFrom(((CustomCredential) credential).getData());
            NativeAuthCallback(LastHandler, googleIdTokenCredential.getIdToken());
          } catch (Exception e) {
            Log.e("gioplugins_auth", "Received an invalid Google ID token response", e);
          }
        } else {
          // Catch any unrecognized custom credential type here.
          Log.e("gioplugins_auth", "Unexpected type of credential");
        }
      } else {
        // Catch any unrecognized credential type here.
        Log.e("gioplugins_auth", "Unexpected type of credential");
      }
    }

    public void handleFailure(GetCredentialException e) {
      // Handle the error.
      Log.e("gioplugins_auth", "Error getting credential", e);
      NativeAuthCallback(LastHandler, "");
    }
}