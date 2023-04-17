package cgroupmemlimited

import (
	"os"
	"runtime/debug"

	"pontus.dev/cgroupmemlimited/internal"
)

var LimitAfterInit int64

func init() {
	if l := internal.Limit(os.DirFS("/sys/fs/cgroup")); l == -1 {
		// SetMemoryLimit(-1) is used to query current limit.
		LimitAfterInit = debug.SetMemoryLimit(l)
	} else {
		debug.SetMemoryLimit(l)
		LimitAfterInit = l
	}
}
