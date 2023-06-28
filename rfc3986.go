// RFC 3986 URI Query Escape/Unescape inspired from "net/url" written in Go
package rfc3986

import (
	"strconv"
	"strings"
)

const upperhex = "0123456789ABCDEF"

type (
	EscapeError      string
	InvalidHostError string
)

func (e EscapeError) Error() string {
	return "invalid URL escape " + strconv.Quote(string(e))
}

func (e InvalidHostError) Error() string {
	return "invalid character " + strconv.Quote(string(e)) + " in host name"
}

// QueryUnescape does the inverse transformation of QueryEscape,
// converting each 3-byte encoded substring of the form "%AB" into the
// hex-decoded byte 0xAB.
// It returns an error if any % is not followed by two hexadecimal
// digits.
func QueryUnescape(s string) (string, error) {
	return unescape(s)
}

// QueryEscape escapes the string so it can be safely placed
// inside a URL query.
func QueryEscape(s string) string {
	return escape(s)
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func shouldEscape(c byte) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false
	default:
		// Everything else must be escaped.
		return true
	}
}

// unescape unescapes a string
func unescape(s string) (string, error) {
	var (
		// Count %, check that they're well-formed.
		n       int  = 0
		lenS    int  = len(s)
		hasPlus bool = false
	)
	for i := 0; i < lenS; {
		switch s[i] {
		case '%':
			n++
			if i+2 >= lenS || !ishex(s[i+1]) || !ishex(s[i+2]) {
				s = s[i:]
				if lenS > 3 {
					s = s[:3]
				}
				return "", EscapeError(s)
			}
			i += 3
		case '+':
			hasPlus = true
			i++
		default:
			i++
		}
	}

	if n == 0 && !hasPlus {
		return s, nil
	}

	var t strings.Builder
	defer t.Reset()
	t.Grow(lenS - 2*n)
	for i := 0; i < lenS; i++ {
		switch s[i] {
		case '%':
			t.WriteByte(unhex(s[i+1])<<4 | unhex(s[i+2]))
			i += 2
		case '+':
			t.WriteByte(' ')
		default:
			t.WriteByte(s[i])
		}
	}
	return t.String(), nil
}

func escape(s string) string {
	lenS, hexCount := len(s), 0
	for i := 0; i < lenS; i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	var (
		buf [64]byte
		t   []byte
	)

	required := lenS + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	j := 0
	for i := 0; i < lenS; i++ {
		switch c := s[i]; {
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
