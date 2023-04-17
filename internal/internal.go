package internal

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	cgroupV1File = "memory/memory.limit_in_bytes"
	cgroupV2File = "memory.max"
)

// Limit returns a new memory limit based on maximum available cgroup memory.
// If no cgroup environment is detected, or if the env var GOMEMLIMIT is set, it
// returns -1.
// Errors reading the cgroup fs, or parsing the value will cause a panic.
func Limit(root fs.FS) int64 {
	// Do not override explicit limit set via env.
	if os.Getenv("GOMEMLIMIT") != "" {
		return -1
	}

	// First try cgroup v2 style.
	f, err := root.Open(cgroupV2File)
	if errors.Is(err, fs.ErrNotExist) {
		// Fall back to cgroup v1
		f, err = root.Open(cgroupV1File)
		if errors.Is(err, fs.ErrNotExist) {
			// Not running with cgroup.
			return -1
		}
	}
	if err != nil {
		panic(fmt.Sprintf("failed to read cgroup memory limit: %v", err))
	}
	defer f.Close()
	v, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("failed to read cgroup memory file: %v", err))
	}
	cgl := strings.TrimSpace(string(v))

	// No cgroup limit set.
	if cgl == "max" {
		return -1
	}

	limit, err := strconv.ParseInt(cgl, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("failed to parse cgroup memory limit as uint: %v", err))
	}

	// Use approx 90% of available memory as limit, as per https://go.dev/doc/gc-guide#Suggested_uses:
	//
	//   Do take advantage of the memory limit when the execution environment of your Go program is
	//   entirely within your control, and the Go program is the only program with access to some
	//   set of resources (i.e. some kind of memory reservation, like a container memory limit).
	//
	//   A good example is the deployment of a web service into containers with a fixed amount
	//   of available memory.
	//
	//   In this case, a good rule of thumb is to leave an additional 5-10% of headroom to
	//   account for memory sources the Go runtime is unaware of.
	return int64(math.Ceil(float64(limit) * 0.9))
}
