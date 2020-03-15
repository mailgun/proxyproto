package proxyproto

import "io"

// readV2Header assumes the passed buf contains the first 13 bytes which should look like
// the following. (Where X is the proto proxy version and command)
// 		"\r\n\r\n\x00\r\nQUIT\nX"
func readV2Header(buf []byte, r io.Reader) (*Header, error) {
	// Read the next 3 bytes which contain the proto and the length of the rest of the header

	// Read the rest of the header as dictated by the length (15-16th bytes)
	//length := binary.BigEndian.Uint16([]byte())
	// POSTFIX thinks the max length of the header could be 1220, but limits the length to 536. We do the same here.

	// Translate the addresses according to the spec
	return nil, nil
}
