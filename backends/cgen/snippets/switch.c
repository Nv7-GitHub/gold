#include "strings.c"

// Switch-case (adapted from https://github.com/haipome/fnv)
#define FNV_32_PRIME 0x01000193
#define FNV1_32_INIT 0x811c9dc5

uint32_t str_hash(string* str) {
  uint32_t hval = FNV1_32_INIT;
  for (int i = 0; i < str->len; i++) {
    hval ^= str->data[i];
    hval *= FNV_32_PRIME;
  }
  return hval;
}