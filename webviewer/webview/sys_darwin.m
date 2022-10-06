#include <TargetConditionals.h>
#include <_cgo_export.h>
#include <stdint.h>

#if TARGET_OS_IPHONE
    @import UIKit;
#else
    @import AppKit;
#endif

@import WebKit;

@interface callbackHandler : NSObject<WKScriptMessageHandler>
@property (nonatomic, assign) uintptr_t handler;
@end

@implementation callbackHandler
    uintptr_t handler;

    - (void)userContentController:(WKUserContentController *)userContentController didReceiveScriptMessage:(WKScriptMessage *)message {
        javascriptManagerCallback(_handler, (char *)[[message body] UTF8String]);
    }
@end

@interface giowebview : WKWebView
@property (nonatomic, assign) uintptr_t handler;
@end

@implementation giowebview
    uintptr_t handler;

    -(void)observeValueForKeyPath:(NSString *)keyPath ofObject:(id)object change:(NSDictionary<NSKeyValueChangeKey,id> *)change context:(void *)context {
        if ([keyPath isEqual:@("URL")]) {
            reportLoadStatus(self.handler, (char *)[[[self URL] absoluteString] UTF8String]);
            return;
        }
        if ([keyPath isEqual:@("title")]) {
            reportTitleStatus(self.handler, (char *)[[self title] UTF8String]);
            return;
        }
    }
@end

CFTypeRef config() {
    WKWebViewConfiguration *conf = [[WKWebViewConfiguration alloc] init];
    return CFBridgingRetain(conf);
}

CFTypeRef create(CFTypeRef config, uintptr_t handler) {
	giowebview *webView = [[giowebview alloc] initWithFrame:CGRectMake(0,0,0,0) configuration: (__bridge WKWebViewConfiguration *)config];
    webView.handler = handler;

    NSString * watch[2] = { @("URL"), @("title") };
    for (uint64_t i = 0; i < 2; i++) {
        [webView addObserver:webView forKeyPath:watch[i] options:NSKeyValueObservingOptionNew context:NULL];
    }

	#if TARGET_OS_IPHONE

	#else
        NSColor * c = [NSColor clearColor];
        webView.layer.backgroundColor = [c CGColor];
        webView.layer.opaque = false;
	#endif

	return CFBridgingRetain(webView);
}

void resize(CFTypeRef web, CFTypeRef windowRef, float x, float y, float w, float h) {
	WKWebView *webView = (__bridge WKWebView *)web;
	#if TARGET_OS_IPHONE
	UIView *view = (__bridge UIView *)windowRef;
	#else
	NSView *view = (__bridge NSView *)windowRef;
	y = (view.bounds.size.height - h) - y;
	#endif

	[webView setFrame: CGRectMake((CGFloat)x, (CGFloat)y, (CGFloat)w, (CGFloat)h)];
}

void show(CFTypeRef web) {
	WKWebView *webView = (__bridge WKWebView *)web;
	[webView setHidden:NO];
}

void hide(CFTypeRef web) {
	WKWebView *webView = (__bridge WKWebView *)web;
	[webView setHidden:YES];
}

void seturl(CFTypeRef web, char *u) {
	WKWebView *webView = (__bridge WKWebView *)web;

	NSURL *url = [NSURL URLWithString:@(u)];
	NSURLRequest *requestObj = [NSURLRequest requestWithURL:url];
	[webView loadRequest:requestObj];
}

void run(CFTypeRef web, CFTypeRef windowRef) {
	WKWebView *webView = (__bridge WKWebView *)web;
	#if TARGET_OS_IPHONE
	UIView *view = [((__bridge UIViewController *)windowRef) view];
	#else
	NSView *view = (__bridge NSView *)windowRef;
	#endif

	[webView setBounds: view.frame];
	[view addSubview:webView];
}

void getCookies(CFTypeRef config, uintptr_t handler, uintptr_t done) {
    NSISO8601DateFormatter *dateFormatting = [[NSISO8601DateFormatter alloc] init];
    WKWebViewConfiguration *configuration = (__bridge WKWebViewConfiguration *)config;
    [[[configuration websiteDataStore] httpCookieStore] getAllCookies: ^(NSArray<NSHTTPCookie *> * array) {
        int i = 0;
        while(i < [array count]) {
            NSHTTPCookie *cookie = [array objectAtIndex:i];

            int cookieType = 0;
            if([cookie isHTTPOnly]) {
                cookieType |= 1;
            }
            if([cookie isSecure]) {
                cookieType |= 2;
            }

            bool next = getCookiesCallback(
                handler,
                cookieType,
                (char *)[[cookie name] UTF8String],
                (char *)[[cookie value] UTF8String],
                (char *)[[cookie domain] UTF8String],
                (char *)[[cookie path] UTF8String],
                (int)[[cookie expiresDate] timeIntervalSince1970]
            );

            if (!next) {
                break;
            }

            i++;
        }
        reportDone(done, nil);
    }];
}

void addCookie(CFTypeRef config, uintptr_t done, char *name, char *value, char *domain, char *path, int64_t expires, uint64_t features) {
    WKWebViewConfiguration *configuration = (__bridge WKWebViewConfiguration *)config;
    NSHTTPCookie *cookie = [NSHTTPCookie cookieWithProperties:@{
       NSHTTPCookieName: @(name),
        NSHTTPCookieValue: @(value),
        NSHTTPCookieDomain: @(domain),
        NSHTTPCookiePath: @(path),
        NSHTTPCookieExpires: [NSDate dateWithTimeIntervalSince1970:expires],
        NSHTTPCookieSecure: ((features) == 2 ? @("true") : @("false")),
    }];

    [[[configuration websiteDataStore] httpCookieStore] setCookie:cookie completionHandler: ^(void) {
        reportDone(done, nil);
    }];
}

void removeCookie(CFTypeRef config, uintptr_t done, char *name, char *domain, char *path) {
    // Copy the information to prevent be freed by the GC/C.free.
    NSString *nameString = @(name);
    NSString *domainString = @(domain);
    NSString *pathString = @(path);

    WKWebViewConfiguration *configuration = (__bridge WKWebViewConfiguration *)config;
    [[[configuration websiteDataStore] httpCookieStore] getAllCookies: ^(NSArray<NSHTTPCookie *> * array) {
        int i = 0;
        while(i < [array count]) {
            NSHTTPCookie *cookie = [array objectAtIndex:i];

            if ([[cookie name] isEqualToString:nameString] && [[cookie domain] isEqualToString:domainString] && [[cookie path] isEqualToString:pathString]) {
                [[[configuration websiteDataStore] httpCookieStore] deleteCookie:cookie completionHandler: ^(void) {
                    reportDone(done, nil);
                }];

                return;
            }

            i++;
        }
        reportDone(done, nil);
    }];
}

void installJavascript(CFTypeRef config, char *js, uint64_t when) {
    WKWebViewConfiguration *configuration = (__bridge WKWebViewConfiguration *)config;
    WKUserContentController *controller = [configuration userContentController];

    WKUserScriptInjectionTime time = WKUserScriptInjectionTimeAtDocumentStart;
    if (when == 1) {
        time = WKUserScriptInjectionTimeAtDocumentEnd;
    }

    [controller addUserScript:[[WKUserScript alloc] initWithSource:@(js) injectionTime:time forMainFrameOnly:false]];
}

void runJavascript(CFTypeRef web, char *js, uintptr_t done) {
    WKWebView *webView = (__bridge WKWebView *)web;
    [webView evaluateJavaScript:@(js) completionHandler: ^(id result, NSError *error) {
        reportDone(done, (char *)[[error localizedDescription] UTF8String]);
    }];
}

void addCallbackJavascript(CFTypeRef config, char *name, uintptr_t handler) {
    WKWebViewConfiguration *configuration = (__bridge WKWebViewConfiguration *)config;
    WKUserContentController *controller = [configuration userContentController];

    callbackHandler *scriptMessageHandler = [callbackHandler alloc];
    scriptMessageHandler.handler = handler;

    [controller addScriptMessageHandler:scriptMessageHandler name:@(name)];
}
