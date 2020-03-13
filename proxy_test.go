package proxyproto_test

import (
	"github.com/mailgun/proxyproto"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseV1Header(t *testing.T) {
	//	r := strings.NewReader("PROXY TCP4 1.1.1.1 1.1.1.1 2 3\r\n")
	//	r := strings.NewReader("PROXY TCP4 192.168.1.1 127.0.0.1 65535 65535\r\n")
	//r := strings.NewReader("PROXY UNKNOWN 0000:0000:0000:0000:0000:0000:0000:0001 0000:0000:0000:0000:0000:0000:0000:0001 65535 65535\r\n")
	//r := strings.NewReader("PROXY UNKNOWN 0000:0000:0000:0000:0000:0000:0000:0001 0000:0000:0000:0000:0000:0000:0000:0001 65535 65535\r\n")
	r := strings.NewReader("PROXY UNKNOWN\r\n")

	// Should error
	//r := strings.NewReader("PROXY UNKNOWN 0000:0000:0000:0000:0000:0000:0000:0001 0000:0000:0000:0000:0000:0000:0000:0001 65535 65535X\r\n")
	//r := strings.NewReader("PROXY TCP6 0000:0000:0000:0000:0000:0000:0000:0001 0000:0000:0000:0000:0000:0000:0000:0001 65535 65535XXXX\r\n")
	_, err := proxyproto.ReadHeader(r)
	assert.NoError(t, err)
}
