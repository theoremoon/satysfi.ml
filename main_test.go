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

func TestTravarseDirectory(t *testing.T) {
	const root = "testdata/travarse"
	got, err := travarseDirectory(root)
	assert.Nil(t, err)

	assert.Equal(t, root, got.Name)
	assert.Equal(t, root, got.Path)
	assert.Equal(t, 2, len(got.Children))
	assert.Equal(t, "001.txt", got.Children[0].Name)
	assert.Equal(t, "/001.txt", got.Children[0].Path)
	assert.Equal(t, "002.txt", got.Children[1].Name)
	assert.Equal(t, "/002.txt", got.Children[1].Path)

	assert.Equal(t, 1, len(got.ChildDirs))
	assert.Equal(t, "003", got.ChildDirs[0].Name)
	assert.Equal(t, "/003", got.ChildDirs[0].Path)

	assert.Equal(t, 1, len(got.ChildDirs[0].Children))
	assert.Equal(t, "001.txt", got.ChildDirs[0].Children[0].Name)
	assert.Equal(t, "/003/001.txt", got.ChildDirs[0].Children[0].Path)
}

func TestVerifyPath(t *testing.T) {
	tests := []struct {
		desc string
		in   string
		want bool
	}{
		{
			desc: "Returns true when path is normal",
			in:   "abcde/ghi",
			want: true,
		},
		{
			desc: "Returns false when path is invalid (directory travarsal)",
			in:   "abcde/../ghi",
			want: false,
		},
		{
			desc: "Returns false when path is invalid (directory travarsal)",
			in:   "../ghi",
			want: false,
		},
		{
			desc: "Returns false when path is invalid (directory travarsal)",
			in:   "../../../../../../../../../../etc/passwd",
			want: false,
		},
		{
			desc: "Returns false when path has prefix '.git'",
			in:   ".git/config",
			want: false,
		},
		{
			desc: "Returns false when path has prefix '.git'",
			in:   ".git",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := verifyPath(tt.in)
			assert.Equal(t, tt.want, got, tt.desc)
		})
	}
}

func TestVerifyID(t *testing.T) {
	tests := []struct {
		desc string
		in   string
		want bool
	}{
		{
			desc: "Returns true when ID has only hex characters",
			in:   randomName(),
			want: true,
		},
		{
			desc: "Returns false when ID contains large characters",
			in:   "0123456789ABCDEF",
			want: false,
		},
		{
			desc: "Returns false when ID contains non hex characters",
			in:   "g",
			want: false,
		},
		{
			desc: "Returns false when ID contains non hex characters",
			in:   "あ",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := verifyID(tt.in)
			assert.Equal(t, tt.want, got, tt.desc)
		})
	}
}

func TestCompile(t *testing.T) {

}
