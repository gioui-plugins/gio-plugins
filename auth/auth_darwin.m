#include "_cgo_export.h"
#import <Foundation/Foundation.h>

#if TARGET_OS_IOS
#import <UIKit/UIKit.h>
#else
#import <Appkit/AppKit.h>
#endif

#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>
#import <AuthenticationServices/AuthenticationServices.h>

@interface GioPluginsAuthContextProvider : NSObject <ASAuthorizationControllerPresentationContextProviding, ASAuthorizationControllerDelegate, ASWebAuthenticationPresentationContextProviding>
@property (nonatomic, assign) CFTypeRef window;
@property (nonatomic, assign) uintptr_t authHandler;
@end

@implementation GioPluginsAuthContextProvider
    CFTypeRef window;
    uintptr_t authHandler;

- (instancetype)initWithWindow:(CFTypeRef)window initWithHandler:(uintptr_t)h {
	self = [super init];
	if (self) {
		self.window = window;
		self.authHandler = h;
	}
	return self;
}

- (void)authorizationController:(ASAuthorizationController *)controller didCompleteWithAuthorization:(ASAuthorization *)authorization {
    if ([authorization.credential isKindOfClass:[ASAuthorizationAppleIDCredential class]]) {
        ASAuthorizationAppleIDCredential * credential = (ASAuthorizationAppleIDCredential *)authorization.credential;

        NSData * identityToken = credential.identityToken;
        NSData * authorizationCode = credential.authorizationCode;

        auth_apple_callback((char *)[authorizationCode bytes], (char *)[identityToken bytes], self.authHandler);
    }
}

- (void)authorizationController:(ASAuthorizationController *)controller didCompleteWithError:(NSError *)error {
    if (error == 0) {
        auth_apple_callback(0, 0, self.authHandler);
    } else {
        if (error.domain == ASAuthorizationErrorDomain) {
            auth_apple_report_error(error.code, self.authHandler);
        } else {
            auth_apple_report_error(0, self.authHandler);
        }
    }
}

- (ASPresentationAnchor)presentationAnchorForWebAuthenticationSession:(ASWebAuthenticationSession *)session {
	#if TARGET_OS_IOS
		UIViewController * v = (__bridge UIViewController *)self.window;
		return (ASPresentationAnchor)([[v view] window]);
	#else
		NSView * v = (__bridge NSView *)self.window;
		return (ASPresentationAnchor)([v window]);
	#endif
}
- (ASPresentationAnchor)presentationAnchorForAuthorizationController:(ASAuthorizationController *)controller {
	#if TARGET_OS_IOS
		UIViewController * v = (__bridge UIViewController *)self.window;
		return (ASPresentationAnchor)([[v view] window]);
	#else
		NSView * v = (__bridge NSView *)self.window;
		return (ASPresentationAnchor)([v window]);
	#endif
}
@end

CFTypeRef gioplugins_auth_createContextProvider(CFTypeRef view, uintptr_t id) {
	return CFBridgingRetain([[GioPluginsAuthContextProvider alloc] initWithWindow:view initWithHandler:(uintptr_t)(id)]);
}

uintptr_t gioplugins_auth_general_open(CFTypeRef contextProvider, char * url, char * scheme, uintptr_t id) {
	NSURL * u = [NSURL URLWithString: [NSString stringWithUTF8String: url]];
	NSString * s = [NSString stringWithUTF8String: scheme];

	ASWebAuthenticationSession * session = [[ASWebAuthenticationSession alloc] initWithURL:u callbackURLScheme:s completionHandler:^(NSURL *callbackURL, NSError *error) {
		free(url);
		free(scheme);

        if (error == 0) {
            auth_general_callback((char *) [callbackURL.absoluteString UTF8String], id);
        } else {
            if (error.domain == ASWebAuthenticationSessionErrorDomain) {
                auth_report_error(error.code, id);
            } else {
                auth_report_error(0, id);
            }
        }
    }];

    session.presentationContextProvider = (__bridge GioPluginsAuthContextProvider *)contextProvider;
    if (![session canStart]) {
        return 1;
    }
	[session start];

	return 0;
}

CFTypeRef gioplugins_auth_apple_open(CFTypeRef contextProvider) {
    ASAuthorizationAppleIDProvider * provider = [[ASAuthorizationAppleIDProvider alloc] init];
    ASAuthorizationAppleIDRequest * request = [provider createRequest];
    request.requestedScopes = @[ASAuthorizationScopeFullName, ASAuthorizationScopeEmail];

	ASAuthorizationController * session = [[ASAuthorizationController alloc] initWithAuthorizationRequests:@[request]];
	session.presentationContextProvider = (__bridge GioPluginsAuthContextProvider *)contextProvider;
	session.delegate = (__bridge GioPluginsAuthContextProvider *)contextProvider;

	[session performRequests];

	return 0;
}