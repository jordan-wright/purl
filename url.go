package main

import (
	"bytes"
	"net/url"
	"strings"
)

// URL is an overridden instance of a url.URL. We override the String()
// function and provide a helper to split the non-empty attributes by a
// separator
type URL struct {
	url.URL
}

// Join returns a string with the non-empty parts of the URL joined with a
// separator
func (u *URL) Join(sep string) string {
	parts := []string{}
	if u.Scheme != "" {
		parts = append(parts, u.Scheme)
	}
	if u.Opaque != "" {
		parts = append(parts, u.Opaque)
	}
	if ui := u.User; ui != nil {
		parts = append(parts, ui.String())
	}
	if u.Host != "" {
		parts = append(parts, u.Host)
	}
	if u.Path != "" {
		parts = append(parts, u.Path)
	}
	if u.RawQuery != "" {
		parts = append(parts, u.RawQuery)
	}
	if u.Fragment != "" {
		parts = append(parts, u.Fragment)
	}
	return strings.Join(parts, sep)
}

func (u *URL) String() string {
	var buf bytes.Buffer
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
	}
	if u.Opaque != "" {
		// jordan-wright - Only print the separator if the scheme and opaque
		// are provided
		if u.Scheme != "" {
			buf.WriteByte(':')
		}
		buf.WriteString(u.Opaque)
	} else {
		if u.Scheme != "" || u.Host != "" || u.User != nil {
			if u.Scheme != "" && u.Host != "" {
				// jordan-wright - Only print the separator if the scheme and
				// host are provided
				buf.WriteString("://")
			}
			// jordan-wright - Only print the @ separator if both the user and
			// hostname are provided
			if ui := u.User; ui != nil {
				buf.WriteString(ui.String())
				if u.Host != "" {
					buf.WriteByte('@')
				}
			}
			if h := u.Host; h != "" {
				buf.WriteString(escape(h, encodeHost))
			}
		}
		path := u.EscapedPath()
		if path != "" && path[0] != '/' && u.Host != "" {
			buf.WriteByte('/')
		}
		if buf.Len() == 0 {
			// RFC 3986 §4.2
			// A path segment that contains a colon character (e.g., "this:that")
			// cannot be used as the first segment of a relative-path reference, as
			// it would be mistaken for a scheme name. Such a segment must be
			// preceded by a dot-segment (e.g., "./this:that") to make a relative-
			// path reference.
			if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i], '/') == -1 {
				buf.WriteString("./")
			}
		}
		buf.WriteString(path)
	}
	if u.ForceQuery || u.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.RawQuery)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(escape(u.Fragment, encodeFragment))
	}
	return buf.String()
}

type encoding int

const (
	encodePath encoding = 1 + iota
	encodePathSegment
	encodeHost
	encodeZone
	encodeUserPassword
	encodeQueryComponent
	encodeFragment
)

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
//
// Please be informed that for now shouldEscape does not check all
// reserved characters correctly. See golang.org/issue/5684.
func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == encodeHost || mode == encodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last three as well. That leaves only ? to escape.
			return c == '?'

		case encodePathSegment: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments.
			return c == '/' || c == ';' || c == ',' || c == '?'

		case encodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	// Everything else must be escaped.
	return true
}

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
