#include "hashmap.c"

typedef uint64_t (*hashfn)(const void *item, uint64_t seed0, uint64_t seed1);
typedef int (*comparefn)(const void *a, const void *b, void *udata);
typedef void (*freefn)(void *item);

typedef struct map {
  struct hashmap* map;
  int refs;
} map;

map* map_new(int elemsize, comparefn comparefn, hashfn hashfn, freefn freefn) {
  map* m = malloc(sizeof(map));
  m->map = hashmap_new(elemsize, 0, 0, 0, hashfn, comparefn, freefn, NULL);
  m->refs = 1;
  return m;
}

void map_free(map* map) {
  map->refs--;
  if (map->refs == 0) {
    hashmap_free(map->map);
    free(map);
  }
}
