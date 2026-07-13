// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Copyright 2026 github.com/coalaura. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "webp.h"
#include "webp/encode.h"
#include "webp/decode.h"
#include "webp/demux.h"
#include "webp/mux.h"

#include <assert.h>
#include <stdlib.h>
#include <string.h>

extern int goWebPWrite(const uint8_t* data, size_t data_size, void* custom_ptr);

static void webpSetDecodeTransform(
	WebPDecoderConfig* config,
	int crop_left, int crop_top, int crop_width, int crop_height,
	int scaled_width, int scaled_height
) {
	if (crop_width > 0 && crop_height > 0) {
		config->options.use_cropping = 1;
		config->options.crop_left = crop_left;
		config->options.crop_top = crop_top;
		config->options.crop_width = crop_width;
		config->options.crop_height = crop_height;
	}
	if (scaled_width > 0 && scaled_height > 0) {
		config->options.use_scaling = 1;
		config->options.scaled_width = scaled_width;
		config->options.scaled_height = scaled_height;
	}
}

int webpGetInfo(
	const uint8_t* data, size_t data_size,
	int* width, int* height,
	int* has_alpha
) {
	WebPBitstreamFeatures features;
	if (WebPGetFeatures(data, data_size, &features) != VP8_STATUS_OK) {
		return 0;
	}
	if(width != NULL) {
		*width = features.width;
	}
	if(height != NULL) {
		*height = features.height;
	}
	if(has_alpha != NULL) {
		*has_alpha = features.has_alpha;
	}
	return 1;
}

uint8_t* webpDecodeGray(
	const uint8_t* data, size_t data_size,
	int* width, int* height
) {
	int w, h;
	uint8_t *y, *u, *v;
	uint8_t *gray, *dst, *src;
	int stride, uv_stride;
	int i;

	if((y = WebPDecodeYUV(data, data_size, &w, &h, &u, &v, &stride, &uv_stride)) == NULL) {
		return NULL;
	}
	if (width != NULL) {
		*width = w;
	}
	if (height != NULL) {
		*height = h;
	}

	if(stride == w) {
		return y;
	}

	if((gray = (uint8_t*)malloc(w*h)) == NULL) {
		free(y);
		return NULL;
	}

	src = y;
	dst = gray;
	for(i = 0; i < h; ++i) {
		memmove(dst, src, w);
		src += stride;
		dst += w;
	}

	free(y);
	return gray;
}

uint8_t* webpDecodeRGB(
	const uint8_t* data, size_t data_size,
	int* width, int* height
) {
	return WebPDecodeRGB(data, data_size, width, height);
}

uint8_t* webpDecodeRGBA(
	const uint8_t* data, size_t data_size,
	int* width, int* height
) {
	return WebPDecodeRGBA(data, data_size, width, height);
}

int webpDecodeRGBIntoDefault(const uint8_t* data, size_t data_size,
	int width, int height, int outStride, uint8_t* out, int use_threads
) {
	WebPDecoderConfig config;
	if (!WebPInitDecoderConfig(&config)) {
		return -1;
	}

	config.options.use_threads = use_threads;
	config.output.colorspace = MODE_RGB;
	config.output.u.RGBA.rgba = out;
	config.output.u.RGBA.stride = outStride;
	config.output.u.RGBA.size = outStride * height;
	config.output.is_external_memory = 1;

	return WebPDecode(data, data_size, &config);
}

int webpDecodeRGBAIntoDefault(const uint8_t* data, size_t data_size,
	int width, int height, int outStride, uint8_t* out, int use_threads
) {
	WebPDecoderConfig config;
	if (!WebPInitDecoderConfig(&config)) {
		return -1;
	}

	config.options.use_threads = use_threads;
	config.output.colorspace = MODE_RGBA;
	config.output.u.RGBA.rgba = out;
	config.output.u.RGBA.stride = outStride;
	config.output.u.RGBA.size = outStride * height;
	config.output.is_external_memory = 1;

	return WebPDecode(data, data_size, &config);
}

int webpDecodeRGBInto(const uint8_t* data, size_t data_size,
	int width, int height, int outStride, uint8_t* out, int use_threads,
	int crop_left, int crop_top, int crop_width, int crop_height,
	int scaled_width, int scaled_height
) {
	WebPDecoderConfig config;
	if (!WebPInitDecoderConfig(&config)) {
		return -1;
	}

	config.options.use_threads = use_threads;
	webpSetDecodeTransform(&config, crop_left, crop_top, crop_width, crop_height, scaled_width, scaled_height);
	config.output.colorspace = MODE_RGB;
	config.output.u.RGBA.rgba = out;
	config.output.u.RGBA.stride = outStride;
	config.output.u.RGBA.size = outStride * height;
	config.output.is_external_memory = 1;

	return WebPDecode(data, data_size, &config);
}

int webpDecodeRGBAInto(const uint8_t* data, size_t data_size,
	int width, int height, int outStride, uint8_t* out, int use_threads,
	int crop_left, int crop_top, int crop_width, int crop_height,
	int scaled_width, int scaled_height
) {
	WebPDecoderConfig config;
	if (!WebPInitDecoderConfig(&config)) {
		return -1;
	}

	config.options.use_threads = use_threads;
	webpSetDecodeTransform(&config, crop_left, crop_top, crop_width, crop_height, scaled_width, scaled_height);
	config.output.colorspace = MODE_RGBA;
	config.output.u.RGBA.rgba = out;
	config.output.u.RGBA.stride = outStride;
	config.output.u.RGBA.size = outStride * height;
	config.output.is_external_memory = 1;

	return WebPDecode(data, data_size, &config);
}

int webpDecodeGrayToSize(const uint8_t* data, size_t data_size,
	int width, int height, int outStride, uint8_t* out, int use_threads,
	int crop_left, int crop_top, int crop_width, int crop_height
) {
	WebPDecoderConfig config;
	if(!WebPInitDecoderConfig(&config)) {
		return -1;
	}

	config.options.use_threads = use_threads;
	webpSetDecodeTransform(&config, crop_left, crop_top, crop_width, crop_height, width, height);
	config.output.colorspace = MODE_YUV;

	int status = WebPDecode(data, data_size, &config);
	if(status != VP8_STATUS_OK) {
		return status;
	}

	int yStride = config.output.u.YUVA.y_stride;
	uint8_t* src = config.output.u.YUVA.y;
	uint8_t* dst = out;
	int i;

	for(i = 0; i < height; ++i) {
		memmove(dst, src, width);
		src += yStride;
		dst += outStride;
	}

	WebPFreeDecBuffer(&config.output);
	return status;
}

int webpDecodeRGBToSize(const uint8_t* data, size_t data_size,
	int width, int height, int outStride, uint8_t* out, int use_threads,
	int crop_left, int crop_top, int crop_width, int crop_height
) {
	WebPDecoderConfig config;
	if(!WebPInitDecoderConfig(&config)) {
		return -1;
	}

	config.options.use_threads = use_threads;
	webpSetDecodeTransform(&config, crop_left, crop_top, crop_width, crop_height, width, height);
	config.output.colorspace = MODE_RGB;
	config.output.u.RGBA.rgba = out;
	config.output.u.RGBA.stride = outStride;
	config.output.u.RGBA.size = outStride * height;
	config.output.is_external_memory = 1;

	return WebPDecode(data, data_size, &config);
}

int webpDecodeRGBAToSize(const uint8_t* data, size_t data_size,
	int width, int height, int outStride, uint8_t* out, int use_threads,
	int crop_left, int crop_top, int crop_width, int crop_height
) {
	WebPDecoderConfig config;
	if(!WebPInitDecoderConfig(&config)) {
		return -1;
	}

	config.options.use_threads = use_threads;
	webpSetDecodeTransform(&config, crop_left, crop_top, crop_width, crop_height, width, height);
	config.output.colorspace = MODE_RGBA;
	config.output.u.RGBA.rgba = out;
	config.output.u.RGBA.stride = outStride;
	config.output.u.RGBA.size = outStride * height;
	config.output.is_external_memory = 1;

	return WebPDecode(data, data_size, &config);
}

uint8_t* webpEncodeGray(
	const uint8_t* gray, int width, int height, int stride, float quality_factor,
	int method, int target_size, int alpha_quality, int autofilter, int thread_level, size_t* output_size
) {
	uint8_t* output;
	WebPConfig config;
	WebPPicture pic;
	WebPMemoryWriter wrt;
	int ok;

	if(!WebPConfigPreset(&config, WEBP_PRESET_DEFAULT, quality_factor)) {
		return NULL;
	}
	config.method = method;
	config.target_size = target_size;
	config.alpha_quality = alpha_quality;
	config.autofilter = autofilter;
	config.thread_level = thread_level;

	if(!WebPPictureInit(&pic)) {
		return NULL;
	}
	pic.width = width;
	pic.height = height;
	pic.writer = WebPMemoryWrite;
	pic.custom_ptr = &wrt;
	WebPMemoryWriterInit(&wrt);
	if(!WebPPictureAlloc(&pic)) {
		WebPPictureFree(&pic);
		return NULL;
	}

	{
		uint8_t* dst = pic.y;
		int y;
		for(y = 0; y < height; ++y) {
			const uint8_t* src = gray + y*stride;
			memcpy(dst, src, width);
			dst += pic.y_stride;
		}

		if(pic.u != NULL && pic.v != NULL) {
			const int uv_width = (width + 1) >> 1;
			const int uv_height = (height + 1) >> 1;
			uint8_t* u = pic.u;
			uint8_t* v = pic.v;
			int i;
			for(i = 0; i < uv_height; ++i) {
				memset(u, 128, uv_width);
				memset(v, 128, uv_width);
				u += pic.uv_stride;
				v += pic.uv_stride;
			}
		}
	}

	ok = WebPEncode(&config, &pic);
	WebPPictureFree(&pic);

	if(!ok) {
		WebPMemoryWriterClear(&wrt);
		return NULL;
	}
	*output_size = wrt.size;
	return wrt.mem;
}

uint8_t* webpEncodeRGB(
	const uint8_t* rgb, int width, int height, int stride, float quality_factor,
	int method, int target_size, int alpha_quality, int autofilter, int thread_level, size_t* output_size
) {
	WebPConfig config;
	WebPPicture pic;
	WebPMemoryWriter wrt;
	int ok;

	if(!WebPConfigPreset(&config, WEBP_PRESET_DEFAULT, quality_factor)) {
		return NULL;
	}
	config.method = method;
	config.target_size = target_size;
	config.alpha_quality = alpha_quality;
	config.autofilter = autofilter;
	config.thread_level = thread_level;

	if(!WebPPictureInit(&pic)) {
		return NULL;
	}
	pic.width = width;
	pic.height = height;
	pic.writer = WebPMemoryWrite;
	pic.custom_ptr = &wrt;
	WebPMemoryWriterInit(&wrt);

	ok = WebPPictureImportRGB(&pic, rgb, stride) && WebPEncode(&config, &pic);
	WebPPictureFree(&pic);

	if(!ok) {
		WebPMemoryWriterClear(&wrt);
		return NULL;
	}
	*output_size = wrt.size;
	return wrt.mem;
}

uint8_t* webpEncodeRGBA(
	const uint8_t* rgba, int width, int height, int stride, float quality_factor,
	int method, int target_size, int alpha_quality, int autofilter, int thread_level, size_t* output_size
) {
	WebPConfig config;
	WebPPicture pic;
	WebPMemoryWriter wrt;
	int ok;

	if(!WebPConfigPreset(&config, WEBP_PRESET_DEFAULT, quality_factor)) {
		return NULL;
	}
	config.method = method;
	config.target_size = target_size;
	config.alpha_quality = alpha_quality;
	config.autofilter = autofilter;
	config.thread_level = thread_level;

	if(!WebPPictureInit(&pic)) {
		return NULL;
	}
	pic.width = width;
	pic.height = height;
	pic.writer = WebPMemoryWrite;
	pic.custom_ptr = &wrt;
	WebPMemoryWriterInit(&wrt);

	ok = WebPPictureImportRGBA(&pic, rgba, stride) && WebPEncode(&config, &pic);
	WebPPictureFree(&pic);

	if(!ok) {
		WebPMemoryWriterClear(&wrt);
		return NULL;
	}
	*output_size = wrt.size;
	return wrt.mem;
}


uint8_t* webpEncodeLosslessGray(
	const uint8_t* gray, int width, int height, int stride,
	int method, int target_size, int alpha_quality, int autofilter, int thread_level, size_t* output_size
) {
	uint8_t* output;
	WebPConfig config;
	WebPPicture pic;
	WebPMemoryWriter wrt;
	int ok;

	if(!WebPConfigPreset(&config, WEBP_PRESET_DEFAULT, 100)) {
		return NULL;
	}
	config.lossless = 1;
	config.method = method;
	config.target_size = target_size;
	config.alpha_quality = alpha_quality;
	config.autofilter = autofilter;
	config.thread_level = thread_level;

	if(!WebPPictureInit(&pic)) {
		return NULL;
	}
	pic.use_argb = 1;
	pic.width = width;
	pic.height = height;
	pic.writer = WebPMemoryWrite;
	pic.custom_ptr = &wrt;
	WebPMemoryWriterInit(&wrt);
	if(!WebPPictureAlloc(&pic)) {
		WebPPictureFree(&pic);
		return NULL;
	}

	{
		uint32_t* dst = pic.argb;
		int y;
		for(y = 0; y < height; ++y) {
			const uint8_t* src = gray + y*stride;
			uint32_t* row = dst + y*pic.argb_stride;
			int x;
			for(x = 0; x < width; ++x) {
				const uint32_t v = (uint32_t)(*src++);
				row[x] = 0xff000000u | (v << 16) | (v << 8) | v;
			}
		}
	}

	ok = WebPEncode(&config, &pic);
	WebPPictureFree(&pic);

	if(!ok) {
		WebPMemoryWriterClear(&wrt);
		return NULL;
	}
	*output_size = wrt.size;
	return wrt.mem;
}


uint8_t* webpEncodeLosslessRGB(
	const uint8_t* rgb, int width, int height, int stride,
	int method, int target_size, int alpha_quality, int autofilter, int thread_level, size_t* output_size
) {
	WebPConfig config;
	WebPPicture pic;
	WebPMemoryWriter wrt;
	int ok;

	if(!WebPConfigPreset(&config, WEBP_PRESET_DEFAULT, 100)) {
		return NULL;
	}
	config.lossless = 1;
	config.method = method;
	config.target_size = target_size;
	config.alpha_quality = alpha_quality;
	config.autofilter = autofilter;
	config.thread_level = thread_level;

	if(!WebPPictureInit(&pic)) {
		return NULL;
	}
	pic.use_argb = 1;
	pic.width = width;
	pic.height = height;
	pic.writer = WebPMemoryWrite;
	pic.custom_ptr = &wrt;
	WebPMemoryWriterInit(&wrt);

	ok = WebPPictureImportRGB(&pic, rgb, stride) && WebPEncode(&config, &pic);
	WebPPictureFree(&pic);

	if(!ok) {
		WebPMemoryWriterClear(&wrt);
		return NULL;
	}
	*output_size = wrt.size;
	return wrt.mem;
}


uint8_t* webpEncodeLosslessRGBA(
	int exact, const uint8_t* rgba, int width, int height, int stride,
	int method, int target_size, int alpha_quality, int autofilter, int thread_level, size_t* output_size
) {
	WebPPicture pic;
	WebPMemoryWriter wrt;
	WebPConfig config;
	int ok;

	if (!WebPConfigPreset(&config, WEBP_PRESET_DEFAULT, 100) || !WebPPictureInit(&pic)) {
		return 0;
	}

	config.lossless = 1;
	config.exact = exact;
	config.method = method;
	config.target_size = target_size;
	config.alpha_quality = alpha_quality;
	config.autofilter = autofilter;
	config.thread_level = thread_level;

	pic.use_argb = 1;
	pic.width = width;
	pic.height = height;

	pic.writer = WebPMemoryWrite;
	pic.custom_ptr = &wrt;
	WebPMemoryWriterInit(&wrt);

	ok = WebPPictureImportRGBA(&pic, rgba, stride) && WebPEncode(&config, &pic);

	WebPPictureFree(&pic);
	if (!ok) {
		WebPMemoryWriterClear(&wrt);
		return 0;
	}
	*output_size = wrt.size;

	return wrt.mem;
}

static int webpGoWriter(const uint8_t* data, size_t data_size, const WebPPicture* picture) {
	return goWebPWrite(data, data_size, picture->custom_ptr);
}

int webpEncodeToWriter(
	const uint8_t* pixels, int width, int height, int stride, int mode,
	int lossless, int exact, int lossless_level, float quality_factor,
	int method, int target_size, int alpha_quality, int autofilter,
	int thread_level, void* custom_ptr
) {
	WebPConfig config;
	WebPPicture pic;
	int ok;

	if (lossless) {
		if (lossless_level >= 0) {
			if (!WebPConfigLosslessPreset(&config, lossless_level)) return 0;
		} else if (!WebPConfigPreset(&config, WEBP_PRESET_DEFAULT, 100)) {
			return 0;
		}
		config.lossless = 1;
	} else if (!WebPConfigPreset(&config, WEBP_PRESET_DEFAULT, quality_factor)) {
		return 0;
	}
	config.exact = exact;
	config.method = method;
	config.target_size = target_size;
	config.alpha_quality = alpha_quality;
	config.autofilter = autofilter;
	config.thread_level = thread_level;
	if (!WebPValidateConfig(&config) || !WebPPictureInit(&pic)) return 0;
	pic.width = width;
	pic.height = height;
	pic.use_argb = lossless;
	pic.writer = webpGoWriter;
	pic.custom_ptr = custom_ptr;

	if (mode == 1) {
		if (!WebPPictureAlloc(&pic)) {
			WebPPictureFree(&pic);
			return 0;
		}
		if (lossless) {
			for (int y = 0; y < height; ++y) {
				const uint8_t* src = pixels + y * stride;
				uint32_t* dst = pic.argb + y * pic.argb_stride;
				for (int x = 0; x < width; ++x) {
					const uint32_t value = src[x];
					dst[x] = 0xff000000u | (value << 16) | (value << 8) | value;
				}
			}
		} else {
			for (int y = 0; y < height; ++y) memcpy(pic.y + y * pic.y_stride, pixels + y * stride, width);
			for (int y = 0; y < (height + 1) / 2; ++y) {
				memset(pic.u + y * pic.uv_stride, 128, (width + 1) / 2);
				memset(pic.v + y * pic.uv_stride, 128, (width + 1) / 2);
			}
		}
		ok = WebPEncode(&config, &pic);
	} else if (mode == 3) {
		ok = WebPPictureImportRGB(&pic, pixels, stride) && WebPEncode(&config, &pic);
	} else {
		ok = WebPPictureImportRGBA(&pic, pixels, stride) && WebPEncode(&config, &pic);
	}
	WebPPictureFree(&pic);
	return ok;
}

char* webpGetEXIF(const uint8_t* data, size_t data_size, size_t* metadata_size) {
	char* metadata = NULL;
	WebPData webp_data = {data, data_size};
	WebPDemuxer* demux = WebPDemux(&webp_data);
	uint32_t flags = WebPDemuxGetI(demux, WEBP_FF_FORMAT_FLAGS);
	*metadata_size = 0;
	if(flags & EXIF_FLAG) {
		WebPChunkIterator it;
		memset(&it, 0, sizeof(it));
		if(WebPDemuxGetChunk(demux, "EXIF", 1, &it)) {
			if(it.chunk.bytes != NULL && it.chunk.size > 0) {
				metadata = (char*)malloc(it.chunk.size);
				memcpy(metadata, it.chunk.bytes, it.chunk.size);
				*metadata_size = it.chunk.size;
			}
		}
		WebPDemuxReleaseChunkIterator(&it);
	}
	WebPDemuxDelete(demux);
	return metadata;
}
char* webpGetICCP(const uint8_t* data, size_t data_size, size_t* metadata_size) {
	char* metadata = NULL;
	WebPData webp_data = {data, data_size};
	WebPDemuxer* demux = WebPDemux(&webp_data);
	uint32_t flags = WebPDemuxGetI(demux, WEBP_FF_FORMAT_FLAGS);
	*metadata_size = 0;
	if(flags & ICCP_FLAG) {
		WebPChunkIterator it;
		memset(&it, 0, sizeof(it));
		if(WebPDemuxGetChunk(demux, "ICCP", 1, &it)) {
			if(it.chunk.bytes != NULL && it.chunk.size > 0) {
				metadata = (char*)malloc(it.chunk.size);
				memcpy(metadata, it.chunk.bytes, it.chunk.size);
				*metadata_size = it.chunk.size;
			}
		}
		WebPDemuxReleaseChunkIterator(&it);
	}
	WebPDemuxDelete(demux);
	return metadata;
}
char* webpGetXMP(const uint8_t* data, size_t data_size, size_t* metadata_size) {
	char* metadata = NULL;
	WebPData webp_data = {data, data_size};
	WebPDemuxer* demux = WebPDemux(&webp_data);
	uint32_t flags = WebPDemuxGetI(demux, WEBP_FF_FORMAT_FLAGS);
	*metadata_size = 0;
	if(flags & XMP_FLAG) {
		WebPChunkIterator it;
		memset(&it, 0, sizeof(it));
		if(WebPDemuxGetChunk(demux, "XMP ", 1, &it)) {
			if(it.chunk.bytes != NULL && it.chunk.size > 0) {
				metadata = (char*)malloc(it.chunk.size);
				memcpy(metadata, it.chunk.bytes, it.chunk.size);
				*metadata_size = it.chunk.size;
			}
		}
		WebPDemuxReleaseChunkIterator(&it);
	}
	WebPDemuxDelete(demux);
	return metadata;
}

static const char* webpMetadataTag(int type) {
	switch (type) {
		case 1: return "EXIF";
		case 2: return "ICCP";
		case 3: return "XMP ";
		default: return NULL;
	}
}

int webpGetMetadataSize(const uint8_t* data, size_t data_size, int type, size_t* metadata_size) {
	const char* tag = webpMetadataTag(type);
	WebPData webp_data = {data, data_size};
	WebPDemuxer* demux;
	WebPChunkIterator it;
	if (metadata_size == NULL || tag == NULL) return 0;
	*metadata_size = 0;
	demux = WebPDemux(&webp_data);
	if (demux == NULL) return 0;
	memset(&it, 0, sizeof(it));
	if (WebPDemuxGetChunk(demux, tag, 1, &it) && it.chunk.bytes != NULL && it.chunk.size > 0) {
		*metadata_size = it.chunk.size;
		WebPDemuxReleaseChunkIterator(&it);
		WebPDemuxDelete(demux);
		return 1;
	}
	WebPDemuxReleaseChunkIterator(&it);
	WebPDemuxDelete(demux);
	return 0;
}

int webpCopyMetadata(const uint8_t* data, size_t data_size, int type, uint8_t* output, size_t output_size) {
	const char* tag = webpMetadataTag(type);
	WebPData webp_data = {data, data_size};
	WebPDemuxer* demux;
	WebPChunkIterator it;
	if (tag == NULL || output == NULL) return 0;
	demux = WebPDemux(&webp_data);
	if (demux == NULL) return 0;
	memset(&it, 0, sizeof(it));
	if (WebPDemuxGetChunk(demux, tag, 1, &it) && it.chunk.bytes != NULL && it.chunk.size == output_size) {
		memcpy(output, it.chunk.bytes, output_size);
		WebPDemuxReleaseChunkIterator(&it);
		WebPDemuxDelete(demux);
		return 1;
	}
	WebPDemuxReleaseChunkIterator(&it);
	WebPDemuxDelete(demux);
	return 0;
}

uint8_t* webpSetEXIF(
	const uint8_t* data, size_t data_size,
	const char* metadata, size_t metadata_size,
	size_t* new_data_size
) {
	WebPData image = {data, data_size};
	WebPData profile = {metadata, metadata_size};
	WebPData output_data = {NULL, 0};
	WebPMux* mux = WebPMuxCreate(&image, 0);
	if(WebPMuxSetChunk(mux, "EXIF", &profile, 0) == WEBP_MUX_OK) {
		WebPMuxAssemble(mux, &output_data);
	}
	WebPMuxDelete(mux);
	*new_data_size = output_data.size;
	return (uint8_t*)(output_data.bytes);
}
uint8_t* webpSetICCP(
	const uint8_t* data, size_t data_size,
	const char* metadata, size_t metadata_size,
	size_t* new_data_size
) {
	WebPData image = {data, data_size};
	WebPData profile = {metadata, metadata_size};
	WebPData output_data = {NULL, 0};
	WebPMux* mux = WebPMuxCreate(&image, 0);
	if(WebPMuxSetChunk(mux, "ICCP", &profile, 0) == WEBP_MUX_OK) {
		WebPMuxAssemble(mux, &output_data);
	}
	WebPMuxDelete(mux);
	*new_data_size = output_data.size;
	return (uint8_t*)(output_data.bytes);
}
uint8_t* webpSetXMP(
	const uint8_t* data, size_t data_size,
	const char* metadata, size_t metadata_size,
	size_t* new_data_size
) {
	WebPData image = {data, data_size};
	WebPData profile = {metadata, metadata_size};
	WebPData output_data = {NULL, 0};
	WebPMux* mux = WebPMuxCreate(&image, 0);
	if(WebPMuxSetChunk(mux, "XMP ", &profile, 0) == WEBP_MUX_OK) {
		WebPMuxAssemble(mux, &output_data);
	}
	WebPMuxDelete(mux);
	*new_data_size = output_data.size;
	return (uint8_t*)(output_data.bytes);
}

uint8_t* webpDelEXIF(const uint8_t* data, size_t data_size, size_t* new_data_size) {
	WebPData image = {data, data_size};
	WebPData output_data = {NULL, 0};
	WebPMux* mux = WebPMuxCreate(&image, 0);
	if(WebPMuxDeleteChunk(mux, "EXIF") == WEBP_MUX_OK) {
		WebPMuxAssemble(mux, &output_data);
	}
	WebPMuxDelete(mux);
	*new_data_size = output_data.size;
	return (uint8_t*)(output_data.bytes);
}
uint8_t* webpDelICCP(const uint8_t* data, size_t data_size, size_t* new_data_size) {
	WebPData image = {data, data_size};
	WebPData output_data = {NULL, 0};
	WebPMux* mux = WebPMuxCreate(&image, 0);
	if(WebPMuxDeleteChunk(mux, "ICCP") == WEBP_MUX_OK) {
		WebPMuxAssemble(mux, &output_data);
	}
	WebPMuxDelete(mux);
	*new_data_size = output_data.size;
	return (uint8_t*)(output_data.bytes);
}
uint8_t* webpDelXMP(const uint8_t* data, size_t data_size, size_t* new_data_size) {
	WebPData image = {data, data_size};
	WebPData output_data = {NULL, 0};
	WebPMux* mux = WebPMuxCreate(&image, 0);
	if(WebPMuxDeleteChunk(mux, "XMP ") == WEBP_MUX_OK) {
		WebPMuxAssemble(mux, &output_data);
	}
	WebPMuxDelete(mux);
	*new_data_size = output_data.size;
	return (uint8_t*)(output_data.bytes);
}

void* webpMalloc(size_t size) {
	return malloc(size);
}

void webpFree(void* p) {
	free(p);
}
