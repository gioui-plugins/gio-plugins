#include "_cgo_export.h"

#import <Foundation/Foundation.h>
#import <UserNotifications/UserNotifications.h>
#import <objc/runtime.h>
#import <TargetConditionals.h>

#if TARGET_OS_IPHONE
#import <UIKit/UIKit.h>
#else
#import <AppKit/AppKit.h>
#endif

@interface PushDelegate : NSObject <UNUserNotificationCenterDelegate>
@property (nonatomic, assign) void* goHandler;
@end

@implementation PushDelegate

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
       willPresentNotification:(UNNotification *)notification
         withCompletionHandler:(void (^)(UNNotificationPresentationOptions options))completionHandler
{
#if TARGET_OS_IPHONE
    completionHandler(UNNotificationPresentationOptionAlert |
                      UNNotificationPresentationOptionSound |
                      UNNotificationPresentationOptionBadge);
#else
    completionHandler(UNNotificationPresentationOptionBanner |
                      UNNotificationPresentationOptionSound |
                      UNNotificationPresentationOptionBadge);
#endif
}

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
didReceiveNotificationResponse:(UNNotificationResponse *)response
         withCompletionHandler:(void(^)(void))completionHandler
{
    completionHandler();
}

@end

static PushDelegate* delegate;

static IMP originalDidRegister = NULL;
static IMP originalDidFail = NULL;

static void sendTokenHex(void* handler, NSData* deviceToken) {
    const unsigned char *data = (const unsigned char *)[deviceToken bytes];
    NSMutableString *token = [NSMutableString string];
    for (NSUInteger i = 0; i < [deviceToken length]; i++) {
        [token appendFormat:@"%02.2hhx", data[i]];
    }
    gioplugins_pushnotification_on_push_token_received(handler, [token UTF8String], NULL);
}

#if TARGET_OS_IPHONE
static void myDidRegister(id self, SEL _cmd, UIApplication* application, NSData* deviceToken)
#else
static void myDidRegister(id self, SEL _cmd, NSApplication* application, NSData* deviceToken)
#endif
{
    if (delegate != nil && delegate.goHandler != NULL) {
        sendTokenHex(delegate.goHandler, deviceToken);
    }

#if TARGET_OS_IPHONE
    if (originalDidRegister != NULL) {
        ((void(*)(id, SEL, UIApplication*, NSData*))originalDidRegister)(self, _cmd, application, deviceToken);
    }
#else
    if (originalDidRegister != NULL) {
        ((void(*)(id, SEL, NSApplication*, NSData*))originalDidRegister)(self, _cmd, application, deviceToken);
    }
#endif
}

#if TARGET_OS_IPHONE
static void myDidFail(id self, SEL _cmd, UIApplication* application, NSError* error)
#else
static void myDidFail(id self, SEL _cmd, NSApplication* application, NSError* error)
#endif
{
    if (delegate != nil && delegate.goHandler != NULL) {
        gioplugins_pushnotification_on_push_token_received(delegate.goHandler, NULL, [[error localizedDescription] UTF8String]);
    }

#if TARGET_OS_IPHONE
    if (originalDidFail != NULL) {
        ((void(*)(id, SEL, UIApplication*, NSError*))originalDidFail)(self, _cmd, application, error);
    }
#else
    if (originalDidFail != NULL) {
        ((void(*)(id, SEL, NSApplication*, NSError*))originalDidFail)(self, _cmd, application, error);
    }
#endif
}

void setupSwizzling(void) {
#if TARGET_OS_IPHONE
    id appDelegate = [UIApplication sharedApplication].delegate;
#else
    id appDelegate = [NSApplication sharedApplication].delegate;
#endif
    if (appDelegate == nil) return;

    Class cls = [appDelegate class];

    SEL registerSel = @selector(application:didRegisterForRemoteNotificationsWithDeviceToken:);
    Method registerMethod = class_getInstanceMethod(cls, registerSel);

    if (registerMethod == NULL) {
        class_addMethod(cls, registerSel, (IMP)myDidRegister, "v@:@@");
    } else {
        originalDidRegister = method_setImplementation(registerMethod, (IMP)myDidRegister);
    }

    SEL failSel = @selector(application:didFailToRegisterForRemoteNotificationsWithError:);
    Method failMethod = class_getInstanceMethod(cls, failSel);

    if (failMethod == NULL) {
        class_addMethod(cls, failSel, (IMP)myDidFail, "v@:@@");
    } else {
        originalDidFail = method_setImplementation(failMethod, (IMP)myDidFail);
    }
}

void requestPushToken(void* handler) {
    if (delegate == nil) {
        delegate = [[PushDelegate alloc] init];
        [UNUserNotificationCenter currentNotificationCenter].delegate = delegate;
    }
    delegate.goHandler = handler;

    UNAuthorizationOptions options =
        (UNAuthorizationOptionAlert |
         UNAuthorizationOptionSound |
         UNAuthorizationOptionBadge);

    [[UNUserNotificationCenter currentNotificationCenter]
        requestAuthorizationWithOptions:options
                      completionHandler:^(BOOL granted, NSError * _Nullable error) {

        if (error != nil) {
            gioplugins_pushnotification_on_push_token_received(handler, NULL, [[error localizedDescription] UTF8String]);
            return;
        }

        if (!granted) {
            gioplugins_pushnotification_on_push_token_received(handler, NULL, "permission denied");
            return;
        }

        dispatch_async(dispatch_get_main_queue(), ^{
#if TARGET_OS_IPHONE
            [[UIApplication sharedApplication] registerForRemoteNotifications];
#else
            [[NSApplication sharedApplication] registerForRemoteNotifications];
#endif
        });
    }];
}