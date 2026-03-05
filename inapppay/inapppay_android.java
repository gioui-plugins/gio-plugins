package com.inkeliz.inapppay_android;

import android.app.Activity;
import android.content.Context;
import android.os.Handler;
import android.os.Looper;
import java.util.ArrayList;
import java.util.List;
import android.view.View;

// Google Play Billing
import com.android.billingclient.api.BillingClient;
import com.android.billingclient.api.BillingClientStateListener;
import com.android.billingclient.api.BillingFlowParams;
import com.android.billingclient.api.BillingResult;
import com.android.billingclient.api.ProductDetails;
import com.android.billingclient.api.ProductDetailsResponseListener;
import com.android.billingclient.api.Purchase;
import com.android.billingclient.api.PurchasesUpdatedListener;
import com.android.billingclient.api.QueryProductDetailsParams;

public class inapppay_android implements PurchasesUpdatedListener {

    private final Handler mainHandler = new Handler(Looper.getMainLooper());
    private BillingClient googleBillingClient;

    static public native void NativeOnProductDetails(long handle, String[] ids, String[] titles, String[] descs, String[] prices, String[] codes, int count);
    static public native void NativeOnPurchaseResult(long handle, String productId, String orderId, int purchaseState, String purchaseData, String signature);
    static public native void NativeOnReportError(long handle, String error);

    private long currentHandle;

    public inapppay_android() {
    }

    private synchronized void setHandle(long handle) {
        this.currentHandle = handle;
    }

    private synchronized long getHandle() {
        return this.currentHandle;
    }

    public void listProductsGoogle(View view, String[] productIDs, long handle) {
        Activity activity = (Activity) view.getContext();

        setHandle(handle);
        initGoogleClient(activity, handle, new Runnable() {
            @Override
            public void run() {
                queryGoogleProducts(productIDs, handle);
            }
        });
    }

    public void purchaseGoogle(View view, String productID, String payload, int isPersonalized, long handle) {
        Activity activity = (Activity) view.getContext();

        setHandle(handle);
        initGoogleClient(activity, handle, new Runnable() {
            @Override
            public void run() {
                launchGooglePurchase(activity, productID, payload, isPersonalized, handle);
            }
        });
    }

    private void initGoogleClient(Context context, long handle, Runnable onReady) {
        if (googleBillingClient != null) {
            if (googleBillingClient.isReady()) {
                onReady.run();
                return;
            }
        } else {
            googleBillingClient = BillingClient.newBuilder(context)
                    .setListener(this)
                    .enablePendingPurchases()
                    .build();
        }

        googleBillingClient.startConnection(new BillingClientStateListener() {
            @Override
            public void onBillingSetupFinished(BillingResult billingResult) {
                if (billingResult.getResponseCode() == BillingClient.BillingResponseCode.OK) {
                    onReady.run();
                } else {
                    NativeOnReportError(handle, "Billing setup failed: " + billingResult.getDebugMessage());
                }
            }

            @Override
            public void onBillingServiceDisconnected() {
            }
        });
    }

    private void queryGoogleProducts(String[] productIDs, long handle) {
        List<QueryProductDetailsParams.Product> productList = new ArrayList<>();
        for (String pid : productIDs) {
            productList.add(
                    QueryProductDetailsParams.Product.newBuilder()
                            .setProductId(pid)
                            .setProductType(BillingClient.ProductType.INAPP)
                            .build());
        }

        QueryProductDetailsParams params = QueryProductDetailsParams.newBuilder()
                .setProductList(productList)
                .build();

        googleBillingClient.queryProductDetailsAsync(params, new ProductDetailsResponseListener() {
            @Override
            public void onProductDetailsResponse(BillingResult billingResult, List<ProductDetails> list) {
                if (billingResult.getResponseCode() == BillingClient.BillingResponseCode.OK) {
                    // Send Arrays
                    int count = list != null ? list.size() : 0;
                    String[] ids = new String[count];
                    String[] titles = new String[count];
                    String[] descs = new String[count];
                    String[] prices = new String[count];
                    String[] codes = new String[count];

                    if (list != null) {
                        for (int i = 0; i < count; i++) {
                            ProductDetails pd = list.get(i);
                            ids[i] = pd.getProductId();
                            titles[i] = pd.getTitle();
                            descs[i] = pd.getDescription();
                            if (pd.getOneTimePurchaseOfferDetails() != null) {
                                prices[i] = pd.getOneTimePurchaseOfferDetails().getFormattedPrice();
                                codes[i] = pd.getOneTimePurchaseOfferDetails().getPriceCurrencyCode();
                            } else {
                                prices[i] = "";
                                codes[i] = "";
                            }
                        }
                    }
                    NativeOnProductDetails(handle, ids, titles, descs, prices, codes, count);
                } else {
                    NativeOnReportError(handle, "Query Failed: " + billingResult.getDebugMessage());
                }
            }
        });
    }

    private void launchGooglePurchase(Activity activity, String productID, String payload, int isPersonalized, long handle) {
        List<QueryProductDetailsParams.Product> productList = new ArrayList<>();
        productList.add(
                QueryProductDetailsParams.Product.newBuilder()
                        .setProductId(productID)
                        .setProductType(BillingClient.ProductType.INAPP)
                        .build());

        QueryProductDetailsParams params = QueryProductDetailsParams.newBuilder()
                .setProductList(productList)
                .build();

        googleBillingClient.queryProductDetailsAsync(params, new ProductDetailsResponseListener() {
            @Override
            public void onProductDetailsResponse(BillingResult billingResult, List<ProductDetails> list) {
                if (billingResult.getResponseCode() == BillingClient.BillingResponseCode.OK && !list.isEmpty()) {
                    ProductDetails pd = list.get(0);
                    BillingFlowParams flowParams = BillingFlowParams.newBuilder()
                            .setProductDetailsParamsList(
                                    List.of(BillingFlowParams.ProductDetailsParams.newBuilder()
                                            .setProductDetails(pd)
                                            .build()))
                            .setObfuscatedAccountId(payload)
                            .setIsOfferPersonalized(isPersonalized == 1)
                            .build();

                    googleBillingClient.launchBillingFlow(activity, flowParams);
                } else {
                    NativeOnReportError(handle, "Product not found for purchase");
                }
            }
        });
    }

    @Override
    public void onPurchasesUpdated(BillingResult billingResult, List<Purchase> purchases) {
        long handle = getHandle();
        if (billingResult.getResponseCode() == BillingClient.BillingResponseCode.OK && purchases != null) {
            for (Purchase purchase : purchases) {
                NativeOnPurchaseResult(handle,
                        !purchase.getProducts().isEmpty() ? purchase.getProducts().get(0) : "",
                        purchase.getOrderId(),
                        1, // 1 = Purchased (Simplified enum)
                        purchase.getOriginalJson(),
                        purchase.getSignature());
            }
        } else if (billingResult.getResponseCode() == BillingClient.BillingResponseCode.USER_CANCELED) {
            NativeOnReportError(handle, "cancelled");
        } else {
            NativeOnReportError(handle, "Purchase failed: " + billingResult.getDebugMessage());
        }
    }
}
