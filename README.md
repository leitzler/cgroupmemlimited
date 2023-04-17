# cgroupmemlimited

This is an experimental module that automatically sets the Go soft memory limit, introduced in Go 1.19, to 90% of the memory limit defined by cgroup.

It is intended to be used when the execution environment of your Go program is entirely within your control, and the Go program is the only program with access to some set of resources (see https://go.dev/doc/gc-guide#Suggested_uses).

## Usage
Add the module to your `go.mod`:
```sh
go get pontus.dev/cgroupmemlimited
```
Then import the package to your Go application:

```go
import _ "pontus.dev/cgroupmemlimited"
```

That's it.

The imported package's `init()` read memory limits from cgroup file(s) that are assumed to exist under `/sys/fs/cgroup`. Older versions of Kubernetes and Docker uses cgroup v1, while newer uses v2. Both versions are supported, and v1 is only used if the v2 file doesn't exist.

The Go soft memory limit isn't changed if no cgroup files are found, or if the env var `GOMEMLIMIT` is set.

Any errors that occurs will cause a panic.

There is one variable defined in the package, `LimitAfterInit`. It will hold the Go soft memory limit as it was defined during `init()`, or the current value if it wasn't changed. E.g.:

```go
package main

import "pontus.dev/cgroupmemlimited"

func main() {
	println("Limit: ", cgroupmemlimited.LimitAfterInit)
}
```
