//go:build !js && !wasm

// Stubs for non-wasm builds so these functions always exist when go builds
// despite the build target. All functions are noops.

package wasm

func RegisterWasmFuncs() {}
