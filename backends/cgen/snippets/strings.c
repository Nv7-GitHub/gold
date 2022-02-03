#include <strings.h>
#include <stdbool.h>

typedef struct string {
  int refs;
  int len;
  char* data;

  bool is_static;
} string;

string* string_new(char* data, int len) {
  string* s = malloc(sizeof(string));
  s->refs = 1;
  s->len = len;
  s->data = data;
  s->is_static = true;
  return s;
}

void string_free(string* s) {
  s->refs--;
  if (s->refs == 0) {
    if (!s->is_static) {
      free(s->data);
    }
    free(s);
  }
}

string* string_concat(string* a, string* b) {
  char* data = malloc(a->len + b->len);
  memcpy(data, a->data, a->len);
  memcpy(data + a->len, b->data, b->len);

  string* c = string_new("", 0);
  c->is_static = false;
  c->data = data;
  return c;
}