#include "_cgo_export.h"
#import <Foundation/Foundation.h>
#import <StoreKit/StoreKit.h>
#include <stdint.h>

// Suppress deprecation warnings for StoreKit 1 APIs
// These are necessary for compatibility with macOS versions before 15.0
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"

@interface InAppPayManager
    : NSObject <SKProductsRequestDelegate, SKPaymentTransactionObserver>
@property(nonatomic, assign) uintptr_t goHandle;
@property(nonatomic, strong)
    NSMutableDictionary<NSString *, SKProduct *> *cachedProducts;
@end

@implementation InAppPayManager

- (instancetype)initWithHandle:(uintptr_t)handle {
  self = [super init];
  if (self) {
    self.goHandle = handle;
    self.cachedProducts = [NSMutableDictionary dictionary];
    [[SKPaymentQueue defaultQueue] addTransactionObserver:self];
  }
  return self;
}

- (void)dealloc {
  [[SKPaymentQueue defaultQueue] removeTransactionObserver:self];
}

- (void)fetchProducts:(NSArray<NSString *> *)productIDs {
  NSSet *productIdentifiers = [NSSet setWithArray:productIDs];
  SKProductsRequest *request = [[SKProductsRequest alloc] initWithProductIdentifiers:productIdentifiers];
  request.delegate = self;
  [request start];
}

#pragma mark - SKProductsRequestDelegate

- (void)productsRequest:(SKProductsRequest *)request didReceiveResponse:(SKProductsResponse *)response {
  NSInteger count = response.products.count;
  if (count == 0) {
    // Check for invalid product identifiers (errors from Apple)
    if (response.invalidProductIdentifiers.count > 0) {
      NSString *invalidIDs = [response.invalidProductIdentifiers componentsJoinedByString:@", "];
      NSString *errorMessage = [NSString stringWithFormat:
          @"Failed to retrieve products. Invalid identifiers (%lu): %@",
          (unsigned long)response.invalidProductIdentifiers.count, invalidIDs];
      gioplugins_inapppay_report_error(self.goHandle, (char *)[errorMessage UTF8String]);
    } else {
      // No products and no invalid identifiers - truly empty response
      gioplugins_inapppay_on_product_details(self.goHandle, NULL, 0);
    }
    return;
  }

  gioplugins_inapppay_product_t *cProducts =
      malloc(sizeof(gioplugins_inapppay_product_t) * count);
  if (!cProducts) {
    gioplugins_inapppay_report_error(self.goHandle, "Memory allocation failed");
    return;
  }

  NSNumberFormatter *formatter = [[NSNumberFormatter alloc] init];
  [formatter setFormatterBehavior:NSNumberFormatterBehavior10_4];
  [formatter setNumberStyle:NSNumberFormatterCurrencyStyle];

  for (NSInteger i = 0; i < count; i++) {
    SKProduct *product = response.products[i];
    [self.cachedProducts setObject:product forKey:product.productIdentifier];

    [formatter setLocale:product.priceLocale];
    NSString *priceStr = [formatter stringFromNumber:product.price];

    cProducts[i].id = (char *)[product.productIdentifier UTF8String];
    cProducts[i].title = (char *)[product.localizedTitle UTF8String];
    cProducts[i].description =
        (char *)[product.localizedDescription UTF8String];
    cProducts[i].price = (char *)[priceStr UTF8String];
    cProducts[i].currencyCode = (char *)[[product.priceLocale
        objectForKey:NSLocaleCurrencyCode] UTF8String];
  }

  gioplugins_inapppay_on_product_details(self.goHandle, cProducts, (int)count);

  free(cProducts);
}

- (void)request:(SKRequest *)request didFailWithError:(NSError *)error {
  gioplugins_inapppay_report_error(
      self.goHandle, (char *)[[error localizedDescription] UTF8String]);
}

#pragma mark - SKPaymentTransactionObserver

- (void)paymentQueue:(SKPaymentQueue *)queue
    updatedTransactions:(NSArray<SKPaymentTransaction *> *)transactions {
  for (SKPaymentTransaction *transaction in transactions) {
    switch (transaction.transactionState) {
    case SKPaymentTransactionStatePurchased:
    case SKPaymentTransactionStateRestored:
      [self handleTransaction:transaction status:1];
      [[SKPaymentQueue defaultQueue] finishTransaction:transaction];
      break;
    case SKPaymentTransactionStateFailed:
      if (transaction.error.code == SKErrorPaymentCancelled) {
        [self handleTransaction:transaction status:2];
      } else {
        gioplugins_inapppay_report_error(
            self.goHandle,
            (char *)[[transaction.error localizedDescription] UTF8String]);
      }
      [[SKPaymentQueue defaultQueue] finishTransaction:transaction];
      break;
    case SKPaymentTransactionStateDeferred:
    case SKPaymentTransactionStatePurchasing:
      break;
    }
  }
}

- (void)handleTransaction:(SKPaymentTransaction *)transaction
                   status:(int)status {
  if (status == 2) {
    gioplugins_inapppay_report_error(self.goHandle, "user cancelled");
    return;
  }

  NSURL *receiptURL = [[NSBundle mainBundle] appStoreReceiptURL];
  NSData *receiptData = [NSData dataWithContentsOfURL:receiptURL];
  NSString *receiptString = [receiptData base64EncodedStringWithOptions:0];

  gioplugins_inapppay_purchase_result_t res;
  res.productID = (char *)[transaction.payment.productIdentifier UTF8String];
  res.purchaseID = (char *)[transaction.transactionIdentifier UTF8String];
  res.status = status;
  res.developerPayload = (char *)[transaction.payment.applicationUsername UTF8String];
  res.originalJSON = (char *)[receiptString UTF8String];
  res.signature = "";

  gioplugins_inapppay_on_purchase_result(self.goHandle, res);
}

@end

static InAppPayManager *manager = nil;

void gioplugins_inapppay_create(uintptr_t data) {
  if (manager == nil) {
    manager = [[InAppPayManager alloc] initWithHandle:data];
  } else {
    manager.goHandle = data;
  }
}

void gioplugins_inapppay_list_products(uintptr_t data, char **productIDs, int count) {
  if (manager && count > 0) {
    NSMutableArray *ids = [NSMutableArray arrayWithCapacity:count];
    for (int i = 0; i < count; i++) {
      if (productIDs[i]) {
        [ids addObject:[NSString stringWithUTF8String:productIDs[i]]];
      }
    }
    [manager fetchProducts:ids];
  }
}

void gioplugins_inapppay_purchase(uintptr_t data, char *productID, char *developerPayload) {
  if (manager) {
    NSString *pid = [NSString stringWithUTF8String:productID];
    SKProduct *product = [manager.cachedProducts objectForKey:pid];

    if (product) {
      SKMutablePayment *payment = [SKMutablePayment paymentWithProduct:product];

      // Set the application username (developer payload / server user ID)
      // This allows the server to identify which user made the purchase
      if (developerPayload && strlen(developerPayload) > 0) {
        payment.applicationUsername = [NSString stringWithUTF8String:developerPayload];
      }

      [[SKPaymentQueue defaultQueue] addPayment:payment];
    } else {
      gioplugins_inapppay_report_error(
          data, "Product not found in cache. Call ListProducts first.");
    }
  }
}

#pragma clang diagnostic pop

