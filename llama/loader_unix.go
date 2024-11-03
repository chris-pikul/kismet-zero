//go:build !windows
// +build !windows

// Copied from github.com/kelindar/search, original copyright is as follows:
//
// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
package llama

import "github.com/ebitengine/purego"

func load(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}
