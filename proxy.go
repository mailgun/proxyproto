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
	Unknown bool
	Client net.Addr
}

const (
	v1Identifier   = "PROXY "
	v1UnKnownProto = "UNKNOWN"
	CRLF = "\r\n"
	v2Identifier   = "\x0D\x0A\x0D\x0A\x00\x0D\x0A\x51\x55\x49\x54\x0A"
)

func ReadHeader(r io.Reader) (*Header, error) {
	var buf [232]byte

	// Read the first 13 bytes which should contain the identifier
	if _, err := io.ReadFull(r, buf[0:13]); err != nil {
		return nil, errors.Wrap(err, "while reading proxy proto identifier")
	}

	// Look for V1 or V2 identifiers
	if bytes.HasPrefix(buf[0:13], []byte(v2Identifier)) {
		h, err := readV2Header(buf, r)
		if err != nil {
			return nil, errors.Wrap(err, "while parsing proxy proto v2 header")
		}
		return h, nil
	}

	if bytes.HasPrefix(buf[0:13], []byte(v1Identifier)) {
		h, err := readV1Header(buf, r)
		if err != nil {
			return nil, errors.Wrap(err, "while parsing proxy proto v1 header")
		}
		return h, nil
	}

	return nil, fmt.Errorf("expected proxy protocol; found '%s' instead", hex.Dump(buf[0:14]))
}

func readV2Header(buf [232]byte, r io.Reader) (*Header, error) {
	return nil, nil
}

// readV1Header assumes the passed buf contains the first 13 bytes which should look like one of
// the following. (Where XX is the start of the tcp address)
// 		"PROXY TCP4 XX"
// 		"PROXY TCP6 XX"
// 		"PROXY UNKNOWN"
func readV1Header(buf [232]byte, r io.Reader) (*Header, error) {
	// For "UNKNOWN", the rest of the line before the CRLF may be omitted by the
	// sender, and the receiver must ignore anything presented before the CRLF is found.
	if bytes.Equal(buf[6:13], []byte(v1UnKnownProto)) {
		b, err := readUntilCRLF(buf, r, 13)
		if err != nil {
			return nil, errors.Wrap(err, "while looking for CRLF after UNKNOWN proto")
		}
		// TODO: save the header in raw form
		fmt.Printf("Unknown: '%s'\n", b)
		return nil, nil
	}

	// Minimum v1 line is `PROXY TCP4 1.1.1.1 1.1.1.1 2 3\r\n` which is 32 bytes, minus the 13 we have
	// already read which leaves 18, so we optimistically read them now.
	if _, err := io.ReadFull(r, buf[13:32]); err != nil {
		return nil, errors.Wrap(err, "while reading proxy proto identifier")
	}

	// If the optimistic read ended in CRLF then no more bytes to read
	if bytes.Equal(buf[30:32], []byte(CRLF)) {
		return parseV1Header(buf[0:30])
	}

	// else we have more bytes to read until we find the CRLF
	b, err := readUntilCRLF(buf, r, 32)
	if err != nil {
		return nil, errors.Wrap(err, "while looking for CRLF after proto")
	}
	return parseV1Header(b)
}


// readUntilCRLF reads from the reader placing the bytes into `buf` starting at `idx` until
// it finds the terminating CRLF or we exceed 107 bytes which is the max length of the v1
// proxy proto header.
func readUntilCRLF(buf [232]byte, r io.Reader, idx int) ([]byte, error) {
	// Read until we find the CRLF or we hit our max possible header length
	for idx < 107 {
		if _, err := r.Read(buf[idx:idx+1]); err != nil {
			return nil, err
		}
		if bytes.Equal(buf[idx-1:idx+1], []byte(CRLF)) {
			return buf[0:idx-1], nil
		}
		idx++
	}
	return nil, errors.New("gave up after 107 bytes")
}

// parseV1Header parses the provided v1 proxy protocol header in the form
// "PROXY TCP4 1.1.1.1 1.1.1.1 2 3" into it's individual parts
func parseV1Header(buf []byte) (*Header, error) {
	fmt.Printf("Parse: '%s'\n", buf)

	/*if !bytes.Equal(buf[6:11], []byte("TCP4 ")) && !bytes.Equal(buf[6:11], []byte("TCP6 ")) {
		return nil, errors.Errorf("unrecognized protocol '%s'", buf[6:10])
	}

	ip := net.ParseIP(buf)
	if ip == nil {
		return nil, errors.Errorf("invalid ip '%s'", buf)
	}

	port, err := strconv.Atoi(buf)
	if err != nil {
		return 0, err
	}
	//client = &net.TCPAddr{IP: ip, Port: port}
	//proxy = &net.TCPAddr{IP: "", Port: ""}
	*/
	return nil, nil
}
