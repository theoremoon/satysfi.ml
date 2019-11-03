package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomName(t *testing.T) {
	name := randomName()
	assert.Equal(t, 16, len(name))
}

func TestIsDir(t *testing.T) {
	tests := []struct {
		desc string
		in   string
		want bool
	}{
		{
			desc: "Returns true when path is directory",
			in:   "docs",
			want: true,
		},
		{
			desc: "Returns false when path is file",
			in:   "README.md",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := isDir(tt.in)
			assert.Equal(t, tt.want, got, tt.desc)
		})
	}
}

func TestIsFile(t *testing.T) {
	tests := []struct {
		desc string
		in   string
		want bool
	}{
		{
			desc: "Returns true when path is directory",
			in:   "docs",
			want: false,
		},
		{
			desc: "Returns false when path is file",
			in:   "README.md",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := isFile(tt.in)
			assert.Equal(t, tt.want, got, tt.desc)
		})
	}
}
