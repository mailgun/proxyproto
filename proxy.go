package proxyproto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
)

type Header struct {
	// True if the proxy header was not UNKNOWN or proto (v2) was not set to LOCAL
	HasProxy bool
	// Contains the complete header minus the CRLF if the proto was UNKNOWN
	Unknown []byte

	// Source is the ip address of the party that initiated the connection
	Source net.TCPAddr
	// Destination is the ip address the remote party connected to; aka the address
	// the proxy was listening for connections on.
	Destination net.TCPAddr
}

const (
	v1Identifier   = "PROXY "
	v1UnKnownProto = "UNKNOWN"
	CRLF           = "\r\n"
	v2Identifier   = "\r\n\r\n\x00\r\nQUIT\n"
)

func ReadHeader(r io.Reader) (*Header, error) {
	var buf [232]byte

	// Read the first 13 bytes which should contain the identifier
	if _, err := io.ReadFull(r, buf[0:13]); err != nil {
		return nil, errors.Wrap(err, "while reading proxy proto identifier")
	}

	// Look for V1 or V2 identifiers
	if bytes.HasPrefix(buf[0:13], []byte(v2Identifier)) {
		h, err := readV2Header(buf[0:], r)
		if err != nil {
			return nil, errors.Wrap(err, "while parsing proxy proto v2 header")
		}
		return h, nil
	}

	if bytes.HasPrefix(buf[0:13], []byte(v1Identifier)) {
		h, err := readV1Header(buf[0:], r)
		if err != nil {
			return nil, errors.Wrap(err, "while parsing proxy proto v1 header")
		}
		return h, nil
	}

	return nil, fmt.Errorf("expected proxy protocol; found '%s' instead", hex.Dump(buf[0:14]))
}
