package internal

import (
	"io"
	"io/fs"
	"testing"
)

func TestSetLimit(t *testing.T) {
	t.Setenv("GOMEMLIMIT", "")
	tt := []struct {
		desc     string
		root     fs.FS
		expected int64
	}{
		{"1k_v1", fakeFS{cgroupV1File: []byte("1000\n")}, 900},
		{"1k_v2", fakeFS{cgroupV2File: []byte("1000\n")}, 900},
		{"no_cgroup", fakeFS{}, -1},
		{"no_limit_v1", fakeFS{cgroupV1File: []byte("max\n")}, -1},
		{"no_limit_v2", fakeFS{cgroupV2File: []byte("max\n")}, -1},
		{"1k_v1_2k_v2", fakeFS{cgroupV1File: []byte("1000\n"), cgroupV2File: []byte("2000\n")}, 1800},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			got := Limit(tc.root)
			if got != tc.expected {
				t.Fatalf("got %d, expected: %d", got, tc.expected)
			}
		})
	}
}

func TestWithEnv(t *testing.T) {
	t.Setenv("GOMEMLIMIT", "3000")
	tt := []struct {
		desc     string
		root     fs.FS
		expected int64
	}{
		{"1k_v1", fakeFS{cgroupV1File: []byte("1000\n")}, -1},
		{"1k_v2", fakeFS{cgroupV2File: []byte("1000\n")}, -1},
		{"no_cgroup", fakeFS{}, -1},
		{"no_limit_v1", fakeFS{cgroupV1File: []byte("max\n")}, -1},
		{"no_limit_v2", fakeFS{cgroupV2File: []byte("max\n")}, -1},
		{"1k_v1_2k_v2", fakeFS{cgroupV1File: []byte("1000\n"), cgroupV2File: []byte("2000\n")}, -1},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			got := Limit(tc.root)
			if got != tc.expected {
				t.Fatalf("got %d, expected: %d", got, tc.expected)
			}
		})
	}
}

type fakeFS map[string]fakeFile

func (ffs fakeFS) Open(name string) (fs.File, error) {
	f, ok := ffs[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return f, nil
}

type fakeFile []byte

func (ff fakeFile) Stat() (fs.FileInfo, error) { panic("N/A") }

func (ff fakeFile) Read(b []byte) (int, error) {
	// Let's assume that the buffer passed in is large enough
	// to hold the entire string, and just panic if it isn't.
	if cap(b) < len(ff) {
		panic("short read")
	}
	n := copy(b, ff)
	return n, io.EOF
}
func (ff fakeFile) Close() error { return nil }
