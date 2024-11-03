//go:build windows
// +build windows

// Copied from github.com/kelindar/search, original copyright is as follows:
//
// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
package llama

import "syscall"

func load(name string) (uintptr, error) {
	// Use [syscall.LoadLibrary] here to avoid external dependencies (#270).
	// For actual use cases, [golang.org/x/sys/windows.NewLazySystemDLL] is recommended.
	handle, err := syscall.LoadLibrary(name)
	return uintptr(handle), err
}
