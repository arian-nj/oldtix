shader_type canvas_item;

uniform vec4 fill_color: source_color = vec4(1.0);
uniform float min_val;
uniform float current_val;
uniform float max_val;

void fragment() {
	float percent = current_val / (max_val - min_val);
	if (UV.x <= percent) {
		COLOR.rgb += fill_color.rgb;
	}
}