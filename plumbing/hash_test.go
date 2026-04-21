package plumbing_test

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
)

func TestNewHash(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid sha1 hash",
			input: "e94f826af816c4c6a0f36e4a2b0d3e8b6c1e2f3a",
			want:  "e94f826af816c4c6a0f36e4a2b0d3e8b6c1e2f3a",
		},
		{
			name:  "zero hash",
			input: "0000000000000000000000000000000000000000",
			want:  "0000000000000000000000000000000000000000",
		},
		{
			name:  "invalid hex string returns zero hash",
			input: "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
			want:  "0000000000000000000000000000000000000000",
		},
		{
			// Partial/short hex strings should also return the zero hash
			name:  "short hex string returns zero hash",
			input: "e94f826af8",
			want:  "0000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := plumbing.NewHash(tt.input)
			if got := h.String(); got != tt.want {
				t.Errorf("NewHash(%q).String() = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestComputeHash(t *testing.T) {
	tests := []struct {
		name    string
		objType plumbing.ObjectType
		content []byte
	}{
		{
			name:    "blob object",
			objType: plumbing.BlobObject,
			content: []byte("hello world\n"),
		},
		{
			name:    "empty blob",
			objType: plumbing.BlobObject,
			content: []byte{},
		},
		{
			name:    "tree object",
			objType: plumbing.TreeObject,
			content: []byte("tree content"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := plumbing.ComputeHash(tt.objType, tt.content)
			if h.IsZero() && len(tt.content) > 0 {
				t.Errorf("ComputeHash returned zero hash for non-empty content")
			}
			// Verify determinism: same input should produce same hash
			h2 := plumbing.ComputeHash(tt.objType, tt.content)
			if h != h2 {
				t.Errorf("ComputeHash is not deterministic: got %v and %v", h, h2)
			}
		})
	}
}

func TestHashIsZero(t *testing.T) {
	var zero plumbing.Hash
	if !zero.IsZero() {
		t.Error("default Hash should be zero")
	}

	nonZero := plumbing.NewHash("e94f826af816c4c6a0f36e4a2b0d3e8b6c1e2f3a")
	if nonZero.IsZero() {
		t.Error("non-zero Hash should not be zero")
	}
}

func TestNewHasher(t *testing.T) {
	hasher := plumbing.NewHasher(plumbing.BlobObject, 11)
	if hasher == nil {
		t.Fatal("NewHasher returned nil")
	}

	_, err := hasher.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("hasher.Write failed: %v", err)
	}

	h := hasher.Sum()
	if h.IsZero() {
		t.Error("hasher.Sum() returned zero hash")
	}
}
