/*
 * MediaInfoDLL.h exposes the library through function pointers that are
 * resolved at runtime with dlopen, which cgo cannot call directly. These
 * thin wrappers give Go plain C functions to call instead.
 *
 * MediaInfoDLL.h defines static state, so it must only be included from
 * this single translation unit.
 */

#include <stddef.h>
#include <MediaInfoDLL/MediaInfoDLL.h>

#include "mediainfo_wrapper.h"

int mediainfo_c_load(void)
{
    return MediaInfoDLL_Load() == (size_t)1;
}

void *mediainfo_c_open(const char *filename)
{
    void *handle = MediaInfo_New();
    if (!handle)
        return NULL;

    if (!MediaInfo_Open(handle, filename)) {
        MediaInfo_Delete(handle);
        return NULL;
    }

    return handle;
}

const char *mediainfo_c_get(void *handle, int kind, size_t stream,
                            const char *key)
{
    return MediaInfo_Get(handle, (MediaInfo_stream_C)kind, stream, key,
                         MediaInfo_Info_Text, MediaInfo_Info_Name);
}

size_t mediainfo_c_count(void *handle, int kind)
{
    return MediaInfo_Count_Get(handle, (MediaInfo_stream_C)kind, (size_t)-1);
}

const char *mediainfo_c_option(void *handle, const char *option,
                               const char *value)
{
    return MediaInfo_Option(handle, option, value);
}

const char *mediainfo_c_inform(void *handle)
{
    return MediaInfo_Inform(handle, 0);
}

void mediainfo_c_close(void *handle)
{
    MediaInfo_Close(handle);
    MediaInfo_Delete(handle);
}
