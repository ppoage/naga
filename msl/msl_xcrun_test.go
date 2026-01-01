package msl

import (
	"testing"

	"github.com/gogpu/naga"
)

func TestMSLCompilesWithXcrun(t *testing.T) {
	const wgslSource = `
struct VertexOutput {
	@builtin(position) position: vec4<f32>,
	@location(0) color: vec3<f32>
}

@vertex
fn vs_main(@builtin(vertex_index) vertex_index: u32) -> VertexOutput {
	let positions = array<vec2<f32>, 3>(
		vec2<f32>(0.0, 0.5),
		vec2<f32>(-0.5, -0.5),
		vec2<f32>(0.5, -0.5)
	);
	let colors = array<vec3<f32>, 3>(
		vec3<f32>(1.0, 0.0, 0.0),
		vec3<f32>(0.0, 1.0, 0.0),
		vec3<f32>(0.0, 0.0, 1.0)
	);
	var out: VertexOutput;
	out.position = vec4<f32>(positions[vertex_index], 0.0, 1.0);
	out.color = colors[vertex_index];
	return out;
}

@fragment
fn fs_main(input: VertexOutput) -> @location(0) vec4<f32> {
	return vec4<f32>(input.color, 1.0);
}
`

	ast, err := naga.Parse(wgslSource)
	if err != nil {
		t.Fatalf("naga.Parse failed: %v", err)
	}

	module, err := naga.LowerWithSource(ast, wgslSource)
	if err != nil {
		t.Fatalf("naga.LowerWithSource failed: %v", err)
	}

	mslSource, _, err := Compile(module, DefaultOptions())
	if err != nil {
		t.Fatalf("msl.Compile failed: %v", err)
	}
	verifyMSLWithXcrun(t, mslSource)
}

func TestMSLCompilesWithXcrun_TexturedStructInputs(t *testing.T) {
	const wgslSource = `
struct VertexInput {
	@location(0) position: vec2<f32>,
	@location(1) color: vec3<f32>,
}

struct VSOut {
	@builtin(position) position: vec4<f32>,
	@location(0) color: vec3<f32>,
}

@group(0) @binding(0) var<uniform> uTint: vec4<f32>;
@group(0) @binding(1) var uSampler: sampler;
@group(0) @binding(2) var uTexture: texture_2d<f32>;

@vertex
fn vs_main(input: VertexInput) -> VSOut {
	var out: VSOut;
	out.position = vec4<f32>(input.position, 0.0, 1.0);
	out.color = input.color;
	return out;
}

@fragment
fn fs_main(input: VSOut) -> @location(0) vec4<f32> {
	let texColor = textureSample(uTexture, uSampler, vec2<f32>(0.5, 0.5));
	return vec4<f32>(input.color, 1.0) * uTint * texColor;
}
`

	ast, err := naga.Parse(wgslSource)
	if err != nil {
		t.Fatalf("naga.Parse failed: %v", err)
	}

	module, err := naga.LowerWithSource(ast, wgslSource)
	if err != nil {
		t.Fatalf("naga.LowerWithSource failed: %v", err)
	}

	mslSource, _, err := Compile(module, DefaultOptions())
	if err != nil {
		t.Fatalf("msl.Compile failed: %v", err)
	}
	verifyMSLWithXcrun(t, mslSource)
}
