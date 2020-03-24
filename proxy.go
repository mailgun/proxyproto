package proxyproto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
)

type Header struct {
	// Source is the ip address of the party that initiated the connection
	Source net.Addr
	// Destination is the ip address the remote party connected to; aka the address
	// the proxy was listening for connections on.
	Destination net.Addr
	// True if the proxy header was UNKNOWN (v1) or if proto was set to LOCAL (v2)
	// In which case Header.Source and Header.Destination will both be nil. TLVs still
	// maybe available if v2, and Header.Unknown will be populated if v1.
	IsLocal bool
	// The version of the proxy protocol parsed
	Version int
	// The unparsed TLVs (Type-Length-Value) that were appended to the end of
	// the v2 proto proxy header.
	RawTLVs []byte
	// Contains the complete header minus the cRLF if the proto was UNKNOWN
	Unknown []byte
}

const (
	v1Identifier   = "PROXY "
	v1UnKnownProto = "UNKNOWN"
	cRLF           = "\r\n"
	v2Identifier   = "\r\n\r\n\x00\r\nQUIT\n"
	tlvHeaderLen   = 3
)

// ParseTLVs parses the Header.RawTLVS byte string into a TLV map
func (h Header) ParseTLVs() (map[byte][]byte, error) {
	tlv := make(map[byte][]byte)

	var offset int
	for offset+tlvHeaderLen < len(h.RawTLVs) {
		length := int(binary.BigEndian.Uint16(h.RawTLVs[offset+1 : offset+3]))

		// Begin points to the beginning of the value
		begin := offset + tlvHeaderLen
		// End points to the end of the value
		end := begin + length
		if end > len(h.RawTLVs) {
			return nil, fmt.Errorf("TLV '0x%X' length '%d' is larger than trailing header", h.RawTLVs[offset], length)
		}

		tlv[h.RawTLVs[offset]] = h.RawTLVs[begin:end]
		offset = offset + end
	}
	return tlv, nil
}

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
