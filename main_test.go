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
	// TODO
}

func TestVerifyID(t *testing.T) {
	// TODO
}

func TestCompile(t *testing.T) {

}
