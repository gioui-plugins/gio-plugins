package com.inkeliz.webview;

import android.os.Bundle;
import android.view.ViewGroup;
import android.app.Activity;
import android.view.View;
import android.view.ViewGroup;
import android.widget.FrameLayout;
import android.view.KeyEvent;
import android.webkit.WebSettings;
import android.content.Context;
import android.webkit.WebViewClient;
import android.widget.Toast;
import android.webkit.WebView;
import android.webkit.WebChromeClient;
import android.util.Log;
import android.os.Build;
import android.os.Parcelable;
import android.net.Proxy;
import java.lang.reflect.*;
import android.util.ArrayMap;
import android.content.Intent;
import java.util.concurrent.Semaphore;
import android.net.http.SslError;
import android.webkit.SslErrorHandler;
import java.security.cert.Certificate;
import android.net.http.SslCertificate;
import java.security.PublicKey;
import java.security.cert.X509Certificate;
import java.security.MessageDigest;
import java.security.cert.CertificateFactory;
import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.util.ArrayList;
import android.os.Build.VERSION;
import android.os.Build.VERSION_CODES;
import android.util.Base64;
import android.webkit.URLUtil;
import android.webkit.WebResourceRequest;
import android.webkit.JavascriptInterface;
import android.graphics.Bitmap;
import android.webkit.ValueCallback;
import android.webkit.CookieManager;
import java.lang.Boolean;

public class sys_android {
    private ViewGroup primaryView;
    private WebView webBrowser;
    private PublicKey[] additionalCerts;
    private ArrayList<String> onLoadAdditionalScripts;
    private ArrayList<String> onFinishAdditionalScripts;

    public class gowebview_javascript {
        public long handler;

        gowebview_javascript(long handler) {
            this.handler = handler;
        }

        @JavascriptInterface
        public String callback(String msg) {
            sendCallback(handler, msg);
            return "";
        }
    }

    public class gowebview_webbrowser extends WebViewClient {
        public long handler;

        @Override
        public void onPageStarted(WebView v, String url, Bitmap favicon) {
           super.onPageStarted(v, url, favicon);
           if (onLoadAdditionalScripts != null) {
                for (int i = 0; i < onLoadAdditionalScripts.size(); i++) {
                    v.evaluateJavascript(onLoadAdditionalScripts.get(i), null);
                }
            }
            reportLoadStatus(handler, url);
        }

        @Override
        public void onPageFinished(WebView v, String url) {
            super.onPageFinished(v, url);
            if (onFinishAdditionalScripts != null) {
                for (int i = 0; i < onFinishAdditionalScripts.size(); i++) {
                    v.evaluateJavascript(onFinishAdditionalScripts.get(i), null);
                }
            }
        }

        @Override public boolean shouldOverrideUrlLoading(WebView v, WebResourceRequest request) {
            String url = request.getUrl().toString();
            if (url.isEmpty()) {
                return false;
            }
            if (URLUtil.isNetworkUrl(url)) {
                return false;
            }
            return true;
        }

        @Override public void onReceivedSslError(WebView v, final SslErrorHandler sslHandler, SslError err){
            if (additionalCerts == null || additionalCerts.length == 0) {
                super.onReceivedSslError(v, sslHandler, err);
                return;
            }

            Certificate certificate = null;
            try{
                if (android.os.Build.VERSION.SDK_INT > android.os.Build.VERSION_CODES.Q) {
                      certificate = err.getCertificate().getX509Certificate();
                } else {
                    // Old APIs doesn't have such .getX509Certificate()
                    Bundle bundle = SslCertificate.saveState(err.getCertificate());
                    byte[] certificateBytes = bundle.getByteArray("x509-certificate");
                    if (certificateBytes != null) {
                        CertificateFactory certFactory = CertificateFactory.getInstance("X.509");
                        certificate = certFactory.generateCertificate(new ByteArrayInputStream(certificateBytes));
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            }

            if (certificate == null) {
                super.onReceivedSslError(v, sslHandler, err);
                return;
            }

            for (int i = 0; i < additionalCerts.length; i++) {
                try{
                    certificate.verify(additionalCerts[i]);
                    sslHandler.proceed();
                    return;
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }

            super.onReceivedSslError(v, sslHandler, err);
        }
    }

    public class gowebview_chrome extends WebChromeClient {
        public long handler;

        @Override
        public void onReceivedTitle(WebView view, String title) {
            super.onReceivedTitle(view, title);
            reportTitleStatus(handler, title);
        }
    }

    public void webview_set_callback(long handler) {
        webBrowser.addJavascriptInterface(new gowebview_javascript(handler), "_callback");
    }

    public void webview_install_javascript(String js, long when, long done) {
        if (when == 0) {
            onLoadAdditionalScripts.add(js);
        } else {
            onFinishAdditionalScripts.add(js);
        }
        reportDone(done, "");
    }

    public void webview_run_javascript(String js, long done) {
        webBrowser.evaluateJavascript(js, new ValueCallback<String>() {
            @Override
            public void onReceiveValue(String value) {
                reportDone(done, "");
            }
        });
    }

    public void webview_create(View v, long handler) {
        if (primaryView == null) {
            if (v instanceof ViewGroup) {
                primaryView = (ViewGroup) v;
            } else {
                primaryView = (ViewGroup) v.getParent();
            }
        }

        onLoadAdditionalScripts = new ArrayList<String>();
        onFinishAdditionalScripts = new ArrayList<String>();

        webBrowser = new WebView(v.getContext());
        WebSettings webSettings = webBrowser.getSettings();
        webSettings.setJavaScriptEnabled(true);
        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O_MR1) {
            webSettings.setSafeBrowsingEnabled(false);
        }
        webSettings.setMixedContentMode(WebSettings.MIXED_CONTENT_COMPATIBILITY_MODE);
        webSettings.setUseWideViewPort(true);
        webSettings.setLoadWithOverviewMode(true);
        webSettings.setDomStorageEnabled(true);
        webSettings.setDatabaseEnabled(true);

        webBrowser.setBackgroundColor(0xFFFFFFFF);

        gowebview_webbrowser wb = new gowebview_webbrowser();
        wb.handler = handler;
        webBrowser.setWebViewClient(wb);

        gowebview_chrome wc = new gowebview_chrome();
        wc.handler = handler;
        webBrowser.setWebChromeClient(wc);

        webBrowser.setVisibility(View.VISIBLE);

        primaryView.addView(webBrowser);
        primaryView.bringChildToFront(webBrowser);
    }

    public void webview_resize(int x, int y, int width, int height) {
        webBrowser.setLayoutParams(new FrameLayout.LayoutParams(width, height));
        webBrowser.setX(x);
        webBrowser.setY(y);
    }

    public void webview_getCookies(long handler, long done) {
        CookieManager cookieManager = CookieManager.getInstance();
        String cookie = cookieManager.getCookie(webBrowser.getUrl());

        getCookiesCallback(handler, cookie);
        reportDone(done, "");
    }

    public void webview_addCookie(String domain, String cookie, long done) {
        CookieManager cookieManager = CookieManager.getInstance();
        cookieManager.setCookie(domain, cookie, new ValueCallback<Boolean>() {
            @Override
            public void onReceiveValue(Boolean value) {
                cookieManager.flush();
                reportDone(done, "");
            }
        });
    }

    public void webview_navigate(String url) {
        webBrowser.loadUrl(url);
    }

    public void webview_destroy() {
        webBrowser.onPause();
        webBrowser.removeAllViews();
        webBrowser.pauseTimers();
        webBrowser.destroy();
    }

    public void webview_show() {
        webBrowser.setVisibility(View.VISIBLE);
    }

    public void webview_hide() {
        webBrowser.setVisibility(View.GONE);
    }

    public boolean webview_proxy(String host, String port) {
        final Semaphore mutex = new Semaphore(0);

        Context app = webBrowser.getContext().getApplicationContext();

        System.setProperty("http.proxyHost", host);
        System.setProperty("http.proxyPort", port);
        System.setProperty("https.proxyHost", host);
        System.setProperty("https.proxyPort", port);

        try {
            Field apk = app.getClass().getDeclaredField("mLoadedApk");
            apk.setAccessible(true);

            Field receivers = Class.forName("android.app.LoadedApk").getDeclaredField("mReceivers");
            receivers.setAccessible(true);

            for (Object map : ((ArrayMap) receivers.get(apk.get(app))).values()) {

                for (Object receiver : ((ArrayMap) map).keySet()) {

                    Class<?> cls = receiver.getClass();
                    if (cls.getName().contains("ProxyChangeListener")) {

                        String proxyInfoName = "android.net.ProxyInfo";
                        if (Build.VERSION.SDK_INT <= Build.VERSION_CODES.KITKAT) {
                            proxyInfoName = "android.net.ProxyProperties";
                        }

                        Intent intent = new Intent(Proxy.PROXY_CHANGE_ACTION);

                        Class proxyInfoClass = Class.forName(proxyInfoName);
                        if (proxyInfoClass != null) {
                            Constructor proxyInfo = proxyInfoClass.getConstructor(String.class, Integer.TYPE, String.class);
                            proxyInfo.setAccessible(true);
                            intent.putExtra("proxy", (Parcelable) ((Object) proxyInfo.newInstance(host, Integer.parseInt(port), null)));
                        }

                        cls.getDeclaredMethod("onReceive", Context.class, Intent.class).invoke(receiver, app, intent);
                    }
                }

            }

            return true;
        } catch(Exception e) {
            return false;
        }
    }

    public boolean webview_certs(String der) {
        String[] sCerts = der.split(";");

        additionalCerts = new PublicKey[sCerts.length];

        for (int i = 0; i < sCerts.length; i++) {
            InputStream streamCert = new ByteArrayInputStream(Base64.decode(sCerts[i], android.util.Base64.DEFAULT));

            try {
                CertificateFactory factory = CertificateFactory.getInstance("X.509");
                 X509Certificate cert = (X509Certificate)factory.generateCertificate(streamCert);

                 additionalCerts[i] = cert.getPublicKey();
            } catch(Exception e) {
                e.printStackTrace();
                return false;
            }
        }

        return true;
    }


    static private native void reportDone(long handler, String error);
    static private native void sendCallback(long handler, String msg);
    static private native void getCookiesCallback(long handler, String cookies);
    static private native void reportTitleStatus(long handler, String title);
    static private native void reportLoadStatus(long handler, String status);
}