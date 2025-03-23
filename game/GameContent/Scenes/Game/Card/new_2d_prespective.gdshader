shader_type canvas_item;

uniform vec4 rect_info;

vec2 uv_to_region_uv(vec2 uv, vec2 texture_pixel_size) {
	return (uv-(rect_info.xy*texture_pixel_size))/(rect_info.zw*texture_pixel_size);
}

vec2 region_uv_to_uv(vec2 region_uv, vec2 texture_pixel_size) {
	return (rect_info.xy*texture_pixel_size)+region_uv*(rect_info.zw*texture_pixel_size);
}

void fragment() {
	vec2 region_uv = uv_to_region_uv(UV, TEXTURE_PIXEL_SIZE);
}