#include "strings.c"
#include <stdio.h>

string* string_itoa(int val) {
  int length = snprintf(NULL, 0, "%d", val);
  char* data = malloc(length);
  snprintf(data, length + 1, "%d", val);
  string* s = string_new(data, length);
  s->is_static = false;
  return s;
}

string* string_ftoa(float val) {
  int length = snprintf(NULL, 0, "%f", val);
  char* data = malloc(length);
  snprintf(data, length + 1, "%f", val);
  string* s = string_new(data, length);
  s->is_static = false;
  return s;
}