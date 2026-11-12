#ifndef GO_MEDIAINFO_WRAPPER_H
#define GO_MEDIAINFO_WRAPPER_H

#include <stddef.h>

int mediainfo_c_load(void);
void *mediainfo_c_open(const char *filename);
const char *mediainfo_c_get(void *handle, int kind, size_t stream,
                            const char *key);
size_t mediainfo_c_count(void *handle, int kind);
const char *mediainfo_c_option(void *handle, const char *option,
                               const char *value);
const char *mediainfo_c_inform(void *handle);
void mediainfo_c_close(void *handle);

#endif
