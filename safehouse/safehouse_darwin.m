#import <Foundation/Foundation.h>
#include <TargetConditionals.h>
#include <_cgo_export.h>
#include <stdint.h>

uint8_t setSecret(char * identifier, char * desc, uint8_t * value, uint64_t value_len) {
   NSString * i = @(identifier);
   NSString * d = @(desc);
   NSData * v = [NSData dataWithBytesNoCopy:value length:value_len freeWhenDone:NO];

   NSMutableDictionary *query = [[NSMutableDictionary alloc] init];
   [query setObject:(__bridge id)kSecClassGenericPassword forKey:(__bridge id)kSecClass];
   [query setObject:d forKey:(__bridge id)kSecAttrDescription];
   [query setObject:i forKey:(__bridge id)kSecAttrAccount];
   [query setObject:v forKey:(__bridge id)kSecValueData];

   OSStatus status = SecItemAdd((__bridge CFDictionaryRef)query, NULL);
   if (status != 0) {
       return 1;
   }
   return 0;
}

/*
CFTypeRef listSecret(char * identifier, uint8_t ** ret, uint32_t * size) {
    NSString * i = [NSString stringWithUTF8String:identifier];
}
*/

CFTypeRef getSecret(char * identifier, uint8_t ** retData, uint32_t * sizeData, char ** retDescription, int8_t multiple) {
   NSMutableDictionary *query = [[NSMutableDictionary alloc] init];
   [query setObject:(__bridge id)kSecClassGenericPassword forKey:(__bridge id)kSecClass];
   if (identifier != NULL) {
       NSString * i = [NSString stringWithUTF8String:identifier];
       [query setObject:i forKey:(__bridge id)kSecAttrAccount];
   }
   [query setObject:(__bridge id)kCFBooleanTrue forKey:(__bridge id)kSecReturnAttributes];
   if (multiple > 0) {
      [query setObject:(__bridge id)kSecMatchLimitAll forKey:(__bridge id)kSecMatchLimit];
      [query setObject:(__bridge id)kCFBooleanFalse forKey:(__bridge id)kSecReturnData];
      // [query setObject:(__bridge id) forKey:(__bridge id)SecMatchPolicy];
   } else {
      [query setObject:(__bridge id)kSecMatchLimitOne forKey:(__bridge id)kSecMatchLimit];
      [query setObject:(__bridge id)kCFBooleanTrue forKey:(__bridge id)kSecReturnData];
   }

   CFDictionaryRef res = NULL;
   OSStatus status = SecItemCopyMatching((__bridge CFDictionaryRef)query, (CFTypeRef *)&res);
   if (status != 0 || res == NULL) {
      return NULL;
   }

   NSDictionary *result = (__bridge NSDictionary*)res;

   NSData * data = (NSData*)[result valueForKey:(__bridge id)kSecValueData];
   *retData = (uint8_t *)[data bytes];
   *sizeData = [data length];

   NSString * desc = [result valueForKey:(__bridge id)kSecAttrDescription];
   *retDescription = (char *)[desc UTF8String];

   return CFBridgingRetain(result);
}

uint8_t removeSecret(char * identifier) {
   NSString * i = [NSString stringWithUTF8String:identifier];

   NSMutableDictionary *query = [[NSMutableDictionary alloc] init];
   [query setObject:(__bridge id)kSecClassGenericPassword forKey:(__bridge id)kSecClass];
   [query setObject:i forKey:(__bridge id)kSecAttrAccount];

   OSStatus status = SecItemDelete((__bridge CFDictionaryRef)query);
   if (status != 0) {
      return 1;
   }
   return 0;
}