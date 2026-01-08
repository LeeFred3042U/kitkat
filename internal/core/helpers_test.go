package core

import (
	"runtime"
	"testing"
)

func TestIsSafePath(t *testing.T) {

	absPath := "/etc/passwd"
	if runtime.GOOS == "windows" {
		absPath = "C:\\Windows\\System32"
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "relative safe path",
			path: "folder/file.txt",
			want: true,
		},
		{
			name: "absolute path",
			path: absPath,
			want: false,
		},
		{
			name: "parent directory traversal",
			path: "../secret.txt",
			want: false,
		},
		{
			name: "nested traversal",
			path: "a/../../b",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSafePath(tt.path)
			if got != tt.want {
				t.Errorf("IsSafePath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
