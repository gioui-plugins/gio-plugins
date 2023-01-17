#import <Foundation/Foundation.h>
#include <TargetConditionals.h>
#include <_cgo_export.h>
#include <stdint.h>

uint8_t setSecret(char * identifier, char * desc, uint8_t * value, uint64_t value_len) {
   NSString * i = @(identifier);
   NSString * d = @(desc);
   NSData * v = [NSData dataWithBytesNoCopy:value length:value_len freeWhenDone:NO];

   NSDictionary *query = @{
    (__bridge id)kSecClass:(__bridge id)kSecClassGenericPassword,
    (__bridge id)kSecAttrDescription: d,
    (__bridge id)kSecAttrAccount: i,
    (__bridge id)kSecValueData: v,
    (__bridge id)kSecAttrAccessible: (__bridge id)kSecAttrAccessibleAfterFirstUnlock,
   };

   OSStatus status = SecItemAdd((__bridge CFDictionaryRef)query, NULL);
   if (status != 0) {
    return 1;
   }
   return 0;
}

uint8_t updateSecret(char * identifier, char * desc, uint8_t * value, uint64_t value_len) {
   NSString * i = @(identifier);
   NSString * d = @(desc);
   NSData * v = [NSData dataWithBytesNoCopy:value length:value_len freeWhenDone:NO];

   NSDictionary *query = @{
    (__bridge id)kSecClass:(__bridge id)kSecClassGenericPassword,
    (__bridge id)kSecAttrDescription: d,
    (__bridge id)kSecAttrAccount: i,
    (__bridge id)kSecValueData: v,
    (__bridge id)kSecAttrAccessible: (__bridge id)kSecAttrAccessibleAfterFirstUnlock,
   };

   NSDictionary *search = @{
    (__bridge id)kSecClass:(__bridge id)kSecClassGenericPassword,
    (__bridge id)kSecAttrAccount: i,
   };

   OSStatus status = SecItemUpdate((__bridge CFDictionaryRef)search, (__bridge CFDictionaryRef)query);
   if (status != 0) {
    return 1;
   }
   return 0;
}

CFTypeRef getSecret(char * identifier, uint32_t * retLength) {
   NSMutableDictionary *query = [[NSMutableDictionary alloc] init];
   [query setObject:(__bridge id)kSecClassGenericPassword forKey:(__bridge id)kSecClass];
   [query setObject:(__bridge id)kCFBooleanTrue forKey:(__bridge id)kSecReturnAttributes];
   if (identifier == NULL) {
      [query setObject:(__bridge id)kSecMatchLimitAll forKey:(__bridge id)kSecMatchLimit];
      [query setObject:(__bridge id)kCFBooleanFalse forKey:(__bridge id)kSecReturnData];
   } else {
      NSString *i = [NSString stringWithUTF8String:identifier];
      [query setObject:i forKey:(__bridge id)kSecAttrAccount];
      [query setObject:(__bridge id)kSecMatchLimitOne forKey:(__bridge id)kSecMatchLimit];
      [query setObject:(__bridge id)kCFBooleanTrue forKey:(__bridge id)kSecReturnData];
   }

   NSArray *array = NULL;
   OSStatus status = 0;
   if (identifier == NULL) {
      CFArrayRef res = NULL;
      status = SecItemCopyMatching((__bridge CFDictionaryRef)query, (CFTypeRef *)&res);
      if (status != 0 || res == NULL) {
        return NULL;
      }
      array = (__bridge NSArray*)res;
   } else {
      CFDictionaryRef res = NULL;
      status = SecItemCopyMatching((__bridge CFDictionaryRef)query, (CFTypeRef *)&res);
      if (status != 0 || res == NULL) {
        return NULL;
      }
      array = @[(__bridge NSDictionary*)res];
   }

   *retLength = [array count];
   if ([array count] == 0) {
       return 0;
   }

   return CFBridgingRetain(array);
}

uint8_t getSecretAt(CFTypeRef array, uint32_t index, char ** retId, char ** retDesc, uint8_t ** retData, uint32_t * sizeData) {
   NSDictionary *result = [(__bridge NSArray*)array objectAtIndex:index];
   if (result == NULL) {
    return 1;
   }

   if (retData != NULL) {
       NSData * data = (NSData*)[result valueForKey:(__bridge id)kSecValueData];
       *retData = (uint8_t *)[data bytes];
       *sizeData = [data length];
   }

   if (retDesc != NULL) {
       *retDesc = (char *)[[result valueForKey:(__bridge id)kSecAttrDescription] UTF8String];
   }

   if (retId != NULL) {
       *retId = (char *)[[result valueForKey:(__bridge id)kSecAttrAccount] UTF8String];
   }

   return 0;
}

uint8_t removeSecret(char * identifier) {
   NSString *i = [NSString stringWithUTF8String:identifier];

   NSDictionary *query = @{
    (__bridge id)kSecClass: (__bridge id)kSecClassGenericPassword,
    (__bridge id)kSecAttrAccount: i,
   };

   OSStatus status = SecItemDelete((__bridge CFDictionaryRef)query);
   if (status != 0) {
      return 1;
   }
   return 0;
}