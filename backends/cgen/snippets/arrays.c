#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct array {
  int len;
  int cap;
  void* data;

  int elemsize;
  int refs;
} array;

array* array_new(int elemsize, int cap) {
  array* a = malloc(sizeof(array));
  a->len = 0;
  a->cap = cap;
  a->data = malloc(elemsize * cap);
  a->elemsize = elemsize;
  a->refs = 1;
  return a;
}

// Get pointer to element at index i.
static inline void* array_get(array* a, int i) {
  return a->data + (i * a->elemsize);
}

static inline void* array_ind(array* a, int i) {
  return *((void**)(a->data + (i * a->elemsize))); // Dereference
}

void array_free(array* a) {
  if (a == NULL) {
    return;
  }

  a->refs--;
  if (a->refs == 0) {
    free(a->data);
    free(a);
  }
}

// static means it can only be used in this file, doesn't matter since there is only one file
void array_grow(array* a, int len) {
  if (a->cap < len) {
    a->cap = len;
    a->data = realloc(a->data, a->elemsize * a->cap);
  }
}

void array_append(array* a, void* val) {
  array_grow(a, a->len + 1);
  memcpy(array_get(a, a->len), &val, a->elemsize);
  a->len++;
}

void array_grow_gold(array* a, int len) {
  array_grow(a, len);
  a->len = len;
}