//go:build ios
// +build ios

#include <UIKit/UIKit.h>
#include <stdint.h>
#include <UniformTypeIdentifiers/UniformTypeIdentifiers.h>
#include "_cgo_export.h"

@implementation explorer_picker
- (void)documentPicker:(UIDocumentPickerViewController *)controller didPickDocumentsAtURLs:(NSArray<NSURL *> *)urls {
    CFTypeRef url = nil;
    if ([urls count] > 0) {
        url = (__bridge_retained CFTypeRef)([urls objectAtIndex:0]);
    }
    [self.picker removeFromParentViewController];
    pickerCallback(url, self.callback);
}

- (void)documentPickerWasCancelled:(UIDocumentPickerViewController *)controller {
    [self.picker removeFromParentViewController];
    pickerCallback(0, self.callback);
}
@end

CFTypeRef saveFile(CFTypeRef view, char * name, uintptr_t callback, CFTypeRef pooled) {
  explorer_picker * explorer = nil;
  if (pooled == 0) {
    explorer = [[explorer_picker alloc] init];
    explorer.picker = [UIDocumentPickerViewController alloc];
    pooled = (__bridge_retained CFTypeRef)explorer;
  } else {
    explorer = (__bridge explorer_picker*) pooled;
  }

   if (@available(iOS 14, *)) {
        explorer.controller = (__bridge UIViewController *)view;
        explorer.callback = callback;
        explorer.picker = [explorer.picker initForExportingURLs:@[[NSURL URLWithString:@(name)]] asCopy:true];
        explorer.picker.delegate = explorer;

        [explorer.controller presentViewController:explorer.picker animated:YES completion:nil];
        return pooled;
    }
    return 0;
}

CFTypeRef openFile(CFTypeRef view, char * ext, uintptr_t callback, CFTypeRef pooled) {
  explorer_picker * explorer = nil;
  if (pooled == 0) {
    explorer = [[explorer_picker alloc] init];
    explorer.picker = [UIDocumentPickerViewController alloc];
    pooled = (__bridge_retained CFTypeRef)explorer;
  } else {
    explorer = (__bridge explorer_picker*) pooled;
  }

  if (@available(iOS 14, *)) {
        explorer.controller = (__bridge UIViewController *)view;
        explorer.callback = callback;

        NSMutableArray<NSString*> *exts = [[@(ext) componentsSeparatedByString:@","] mutableCopy];
        NSMutableArray<UTType*> *contentTypes = [[NSMutableArray alloc]init];

        int i;
        for (i = 0; i < [exts count]; i++) {
            UTType *utt = [UTType typeWithFilenameExtension:exts[i]];
            if (utt != nil) {
                [contentTypes addObject:utt];
            }
        }

        explorer.picker = [explorer.picker initForOpeningContentTypes:contentTypes asCopy:true];
        explorer.picker.delegate = explorer;

        CFTypeRef ref = (__bridge_retained CFTypeRef)explorer;
        [explorer.controller presentViewController:explorer.picker animated:YES completion:nil];
        return pooled;
    }
    return 0;
}