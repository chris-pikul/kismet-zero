package llama

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ebitengine/purego"
)

type NUMASetting uint8

const (
	GGML_NUMA_STRATEGY_DISABLED   NUMASetting = 0
	GGML_NUMA_STRATEGY_DISTRIBUTE NUMASetting = 1
	GGML_NUMA_STRATEGY_ISOLATE    NUMASetting = 2
	GGML_NUMA_STRATEGY_NUMACTL    NUMASetting = 3
	GGML_NUMA_STRATEGY_MIRROR     NUMASetting = 4
	GGML_NUMA_STRATEGY_COUNT      NUMASetting = 5
)

type LogLevel uint8

const (
	GGML_LOG_LEVEL_NONE  LogLevel = 0
	GGML_LOG_LEVEL_DEBUG LogLevel = 1
	GGML_LOG_LEVEL_INFO  LogLevel = 2
	GGML_LOG_LEVEL_WARN  LogLevel = 3
	GGML_LOG_LEVEL_ERROR LogLevel = 4
	GGML_LOG_LEVEL_CONT  LogLevel = 5
)

// libptr is a pointer to the loaded dynamic library.
var libptr uintptr

var did_error func() bool
var get_last_error func() string
var clear_error func()
var get_last_output func() string
var clear_output func()
var set_log_level func(n LogLevel)

var set_predict func(n int32)
var set_context_size func(n int32)
var set_batch func(n int32)
var set_ubatch func(n int32)
var set_draft_size func(n int32)
var set_parallel func(n int32)
var set_gpu_layers func(n int32)

var set_seed func(n int32)
var set_min_keep func(n int32)
var set_top_k func(n int32)
var set_top_p func(n float32)
var set_typical_p func(n float32)
var set_min_p func(n float32)
var set_temp func(n float32)

var init_library func(numa NUMASetting)
var free_library func()
var load_model func(path string) bool
var free_model func()

var infer_sync func(prompt string) bool

func init() {
	libpath, err := findLlama()
	if err != nil {
		panic(err)
	}
	if libptr, err = load(libpath); err != nil {
		panic(err)
	}

	// Load the library functions
	purego.RegisterLibFunc(&did_error, libptr, "did_error")
	purego.RegisterLibFunc(&get_last_error, libptr, "get_last_error")
	purego.RegisterLibFunc(&clear_error, libptr, "clear_error")
	purego.RegisterLibFunc(&get_last_output, libptr, "get_last_output")
	purego.RegisterLibFunc(&clear_output, libptr, "clear_output")
	purego.RegisterLibFunc(&set_log_level, libptr, "set_log_level")

	purego.RegisterLibFunc(&set_predict, libptr, "set_predict")
	purego.RegisterLibFunc(&set_context_size, libptr, "set_context_size")
	purego.RegisterLibFunc(&set_batch, libptr, "set_batch")
	purego.RegisterLibFunc(&set_ubatch, libptr, "set_ubatch")
	purego.RegisterLibFunc(&set_draft_size, libptr, "set_draft_size")
	purego.RegisterLibFunc(&set_parallel, libptr, "set_parallel")
	purego.RegisterLibFunc(&set_gpu_layers, libptr, "set_gpu_layers")

	purego.RegisterLibFunc(&set_seed, libptr, "set_seed")
	purego.RegisterLibFunc(&set_min_keep, libptr, "set_min_keep")
	purego.RegisterLibFunc(&set_top_k, libptr, "set_top_k")
	purego.RegisterLibFunc(&set_top_p, libptr, "set_top_p")
	purego.RegisterLibFunc(&set_typical_p, libptr, "set_typical_p")
	purego.RegisterLibFunc(&set_min_p, libptr, "set_min_p")
	purego.RegisterLibFunc(&set_temp, libptr, "set_temp")

	purego.RegisterLibFunc(&init_library, libptr, "init_library")
	purego.RegisterLibFunc(&free_library, libptr, "free_library")
	purego.RegisterLibFunc(&load_model, libptr, "load_model")
	purego.RegisterLibFunc(&free_model, libptr, "free_model")

	purego.RegisterLibFunc(&infer_sync, libptr, "infer_sync")
}

// findLlama searches for the dynamic library in standard system paths.
func findLlama() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return findLibrary("llama_go.dll", runtime.GOOS)
	case "darwin":
		return findLibrary("libllama_go.dylib", runtime.GOOS)
	default:
		return findLibrary("libllama_go.so", runtime.GOOS)
	}
}

// findLibrary searches for a dynamic library by name across standard system paths.
// It returns the full path to the library if found, or an error listing all searched paths.
func findLibrary(libName, goos string, dirs ...string) (string, error) {
	libExt, commonPaths := findLibDirs(goos)
	dirs = append(dirs, commonPaths...)

	// Append the correct extension if missing
	if !strings.HasSuffix(libName, libExt) {
		libName += libExt
	}

	// Include current working directory
	if cwd, err := os.Getwd(); err == nil {
		dirs = append(dirs, cwd)

		// Add llama package folder
		dirs = append(dirs, cwd+"/llama")
	}

	// Iterate through directories and search for the library
	searched := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		filename := filepath.Join(dir, libName)
		searched = append(searched, filename)
		if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
			return filename, nil // Library found
		}
	}

	// Construct error message listing all searched paths
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Library '%s' not found, checked following paths:\n", libName))
	for _, path := range searched {
		sb.WriteString(fmt.Sprintf(" - %s\n", path))
	}

	return "", errors.New(sb.String())
}

// findLibDirs returns the library extension, relevant environment path, and common library directories based on the OS.
func findLibDirs(goos string) (string, []string) {
	switch goos {
	case "windows":
		systemRoot := os.Getenv("SystemRoot")
		return ".dll", append(
			filepath.SplitList(os.Getenv("PATH")),
			filepath.Join(systemRoot, "System32"),
			filepath.Join(systemRoot, "SysWOW64"),
		)
	case "darwin":
		return ".dylib", append(
			filepath.SplitList(os.Getenv("DYLD_LIBRARY_PATH")),
			"/usr/lib",
			"/usr/local/lib",
		)
	default: // Unix/Linux
		return ".so", append(
			filepath.SplitList(os.Getenv("LD_LIBRARY_PATH")),
			"/lib",
			"/usr/lib",
			"/usr/local/lib",
		)
	}
}
