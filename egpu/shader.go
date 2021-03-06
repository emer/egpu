// Copyright (c) 2022, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is initially adapted from https://github.com/vulkan-go/asche
// Copyright © 2017 Maxim Kupriianov <max@kc.vc>, under the MIT License

package egpu

import (
	"io/ioutil"
	"unsafe"

	"github.com/goki/gi/oswin/gpu"
	vk "github.com/vulkan-go/vulkan"
)

// LoadShader compiles given SPIR-V ".spv" code and returns a shader -- not typically saved!
func LoadShader(src []byte) (vk.ShaderModule, error) {
	var module vk.ShaderModule
	ret := vk.CreateShaderModule(TheGPU.Device.Device, &vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(src)),
		PCode:    SliceUint32(src),
	}, nil, &module)
	if IsError(ret) {
		return nil, NewError(ret)
	}
	return module, nil
}

// Shader manages a single Shader program
type Shader struct {
	init   bool
	Shader vk.ShaderModule
	name   string
	typ    gpu.ShaderTypes
	src    string
	orgSrc string // original source as provided by user -- program adds extra source..
}

// Name returns the unique name of this Shader
func (sh *Shader) Name() string {
	return sh.name
}

// Type returns the type of the Shader
func (sh *Shader) Type() gpu.ShaderTypes {
	return sh.typ
}

// OpenFile returns bytes from file
func OpenFile(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}

// Compile compiles given src as string (gl version -- not relevant)
func (sh *Shader) Compile(src string) error {
	return sh.CompileSPV([]byte(src))
}

// CompileSPV compiles given SPIR-V ".spv" code for the Shader, of given type and unique name.
func (sh *Shader) CompileSPV(src []byte) error {
	var module vk.ShaderModule
	ret := vk.CreateShaderModule(TheGPU.Device.Device, &vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(src)),
		PCode:    SliceUint32(src),
	}, nil, &module)
	if IsError(ret) {
		return NewError(ret)
	}
	sh.Shader = module
	sh.src = string(src)
	sh.init = true
	return nil
}

// Source returns the actual final source code for the Shader
// excluding the null terminator (for display purposes).
// This includes extra auto-generated code from the Program.
func (sh *Shader) Source() string {
	return sh.src
}

// OrigSource returns the original user-supplied source code
// excluding the null terminator (for display purposes)
func (sh *Shader) OrigSource() string {
	return sh.orgSrc
}

// Delete deletes the Shader
func (sh *Shader) Delete() {
	if !sh.init {
		return
	}
	vk.DestroyShaderModule(TheGPU.Device.Device, sh.Shader, nil)
	sh.Shader = nil
	sh.init = false
}

// todo: use 1.17 unsafe.Slice function
// https://stackoverflow.com/questions/11924196/convert-between-slices-of-different-types

func SliceUint32(data []byte) []uint32 {
	const m = 0x7fffffff
	return (*[m / 4]uint32)(unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&data)).Data))[:len(data)/4]
}

type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
