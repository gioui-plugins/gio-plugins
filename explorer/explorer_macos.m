// SPDX-License-Identifier: Unlicense OR MIT

//go:build darwin && !ios
// +build darwin,!ios

#include "_cgo_export.h"
#import <Foundation/Foundation.h>
#import <Appkit/AppKit.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>

void gioplugins_explorer_saveFile(CFTypeRef viewRef, char *name, uintptr_t id) {
    NSView *view = (__bridge NSView *) viewRef;

    NSSavePanel *panel = [NSSavePanel savePanel];

    [panel setNameFieldStringValue:@(name)];
    [panel beginSheetModalForWindow:[view window] completionHandler:^(NSInteger result) {
        if (result == NSModalResponseOK) {
            gioplugins_explorer_exportCallback((char *) [[panel URL].absoluteString UTF8String], id);
        } else {
            gioplugins_explorer_exportCallback((char *) (""), id);
        }
    }];
}

void gioplugins_explorer_openFile(CFTypeRef viewRef, char *ext, uintptr_t id) {
    NSView *view = (__bridge NSView *) viewRef;

    NSOpenPanel *panel = [NSOpenPanel openPanel];

    NSMutableArray < NSString * > *exts = [[@(ext) componentsSeparatedByString:@","] mutableCopy];
    NSMutableArray < UTType * > *contentTypes = [[NSMutableArray alloc] init];

    int i;
    for (i = 0; i < [exts count]; i++) {
        UTType *utt = [UTType typeWithFilenameExtension:exts[i]];
        if (utt != nil) {
            [contentTypes addObject:utt];
        }
    }

    [(NSSavePanel *) panel setAllowedContentTypes:[NSArray arrayWithArray:contentTypes]];
    [panel beginSheetModalForWindow:[view window] completionHandler:^(NSInteger result) {
        if (result == NSModalResponseOK) {
            gioplugins_explorer_importCallback((char *) [[panel URL].absoluteString UTF8String], id);
        } else {
            gioplugins_explorer_importCallback((char *) (""), id);
        }
    }];
}