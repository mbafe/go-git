// Package plumbing implements the core data structures and operations
// for the go-git library.
package plumbing

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"strings"
)

// Hash represents a SHA-1 hash of a git object.
type Hash [20]byte

// ZeroHash is the zero-value Hash, representing the absence of a hash.
var ZeroHash Hash

// NewHash creates a new Hash from a hex string.
// Returns ZeroHash if the string is invalid.
// Handles both uppercase and lowercase hex strings.
// Note: partial hashes (fewer than 40 hex chars) are not supported and will return ZeroHash.
func NewHash(s string) Hash {
	// Trim whitespace to be more forgiving of input (e.g. trailing newlines from shell output)
	s = strings.TrimSpace(s)
	b, err := hex.DecodeString(strings.ToLower(s))
	if err != nil || len(b) != 20 {
		return ZeroHash
	}

	var h Hash
	copy(h[:], b)
	return h
}

// String returns the hex string representation of the Hash.
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// IsZero returns true if the Hash is the zero-value Hash.
func (h Hash) IsZero() bool {
	return h == ZeroHash
}

// ComputeHash computes the SHA-1 hash of the given object type and content.
// Git hashes are computed as: sha1("<type> <size>\0<content>").
func ComputeHash(t ObjectType, content []byte) Hash {
	h := sha1.New()
	_, _ = fmt.Fprintf(h, "%s %d\x00", t, len(content))
	_, _ = h.Write(content)

	var hash Hash
	copy(hash[:], h.Sum(nil))
	return hash
}

// NewHasher returns a new hasher that computes the SHA-1 hash of a git object.
// The returned hash.Hash already has the git object header written into it.
func NewHasher(t ObjectType, size int64) (hash.Hash, error) {
	h := sha1.New()
	_, err := fmt.Fprintf(h, "%s %d\x00", t, size)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// HashReader computes the hash of the content read from r, treating it as
// an object of the given type with the given size.
// The caller is responsible for ensuring that r contains exactly size bytes;
// a mismatch will produce a hash that does not match the stored object.
//
// Note: io.Copy reads in 32 KiB chunks by default, which should be fine for
// most objects. For very large blobs this may be worth revisiting.
// TODO: consider adding a strict mode that returns an error on size mismatch.
func HashReader(t ObjectType, size int64, r io.Reader) (Hash, error) {
	h, err := NewHasher(t, size)
	if err != nil {
		return ZeroHash, err
	}

	if _, err := io.Copy(h, r); err != nil {
		return ZeroHash, err
	}

	var result Hash
	copy(result[:], h.Sum(nil))
	return result, nil
}
