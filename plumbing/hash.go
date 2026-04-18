// Package plumbing implements the core data structures and operations
// for the go-git library.
package plumbing

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

// Hash represents a SHA-1 hash of a git object.
type Hash [20]byte

// ZeroHash is the zero-value Hash, representing the absence of a hash.
var ZeroHash Hash

// NewHash creates a new Hash from a hex string.
// Returns ZeroHash if the string is invalid.
func NewHash(s string) Hash {
	b, err := hex.DecodeString(s)
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
