package main

import (
	"fmt"
	"io"

	"golang.org/x/crypto/openpgp/packet"
)

func getRecipients(r io.Reader) ([]string, error) {
	pr := packet.NewReader(r)
	var ret []string
	for {
		p, err := pr.Next()
		if err != nil {
			return nil, err
		}
		if key, ok := p.(*packet.EncryptedKey); ok {
			ret = append(ret, fmt.Sprintf("%X", key.KeyId))
		} else {
			// According to RFC 4880 section 11.3, encrypted keys must be the
			// first thing in the file. If we get a packet that isn't a key,
			// that means we've already found all of the keys.
			break
		}
	}
	return ret, nil
}

// containsAny returns true if any element of needles occurs in haystack; i.e.,
// if intersection(haystack, needles) is non-empty.
func containsAny(haystack []string, needles []string) bool {
	// haystack and needles should be smallish... just use the simpler and O(1) memory but O(n*m) time algorithm
	for _, h := range haystack {
		for _, n := range needles {
			if h == n {
				return true
			}
		}
	}
	return false
}
