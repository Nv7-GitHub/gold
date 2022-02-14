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

void* array_get(array* a, int i) {
  return a->data + (i * a->elemsize);
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

void array_append(array* a, void* val) {
  if (a->len == a->cap) {
    a->cap *= 2;
    a->data = realloc(a->data, a->elemsize * a->cap);
  }

  memcpy(array_get(a, a->len), val, a->elemsize);
  a->len++;
}

void array_grow(array* a, int len) {
  if (a->cap < len) {
    a->cap = len;
    a->data = realloc(a->data, a->elemsize * a->cap);
  }

  a->len = len;
}