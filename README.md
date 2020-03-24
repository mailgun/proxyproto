## Proxy Protocol
This is a golang implementation of the proxy protocol as described
[here](https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt)

This is a prerequisite to implementing TCP servers that need to identify the
connecting client and sit behind a load balancer. For instance, an SMTP server
must be able to identify the source ip and port of the connecting client. Load
balancers that wish to pass on this information implement the proxy protocol
when proxying requests downstream. AWS load balancers support V1 and V2 across
various products. To maximize flexability both protocols are supported but V2 is
recommended as it is the most performant.

This implementation is focused on performance and completeness.

Current benchmarks for V1 and V2 implemenation
```
BenchmarkReadHeaderV1/TCP4-Minimal-8             2323470               519 ns/op
BenchmarkReadHeaderV1/TCP4-Maximal-8             1156881              1042 ns/op
BenchmarkReadHeaderV1/TCP4-Typical-8             1293009               929 ns/op
BenchmarkReadHeaderV1/TCP6-Minimal-8             2402085               501 ns/op
BenchmarkReadHeaderV1/TCP6-Maximal-8              497880              2164 ns/op

BenchmarkReadHeaderV2/TCP6-With-TLVs-8           3514410               342 ns/op
BenchmarkReadHeaderV2/TCP6-Minimal-8             3412977               345 ns/op
BenchmarkReadHeaderV2/TCP6-Maximal-8             3550310               340 ns/op
BenchmarkReadHeaderV2/TCP4-Minimal-8             3066440               396 ns/op
```

## Installation

```bash
$ go get github.com/mailgun/proxyproto
```

## Usage
The library is designed to read the header right after the TCP connection is
made. This is unlike other implementations which place themselves directly in
the IO read path. Since `protoproxy` only reads from the connection once,
performance for the rest of the server is not impacted.

```go
func main() {
	ln, err := net.Listen(network, port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Serve '%s' on: %s....\n", network, port)

	for {
		conn, e := ln.Accept()
		fmt.Printf("New Connection: %s\n", conn.RemoteAddr().String())
		if e != nil {
			panic(err)
		}

		h, err := proxyproto.ReadHeader(conn)
		if err != nil {
			panic(err)
		}

		// True if the proxy header was UNKNOWN (v1) or if proto was set to LOCAL (v2)
		// In which case Header.Source and Header.Destination will both be nil. TLVs still
		// maybe available if v2, and Header.Unknown will be populated if v1.
		if h.IsLocal {
			// Local Proxy command for V2 may have some TLVs
			if h.Version == 2 && len(h.RawTLVs) != 0 {
				spew.Dump(h.ParseTLVs())
			}
			// UNKNOWN Proxy command for V1 might have some data we are interested in
			if h.Version == 1 && len(h.Unknown) != 0 {
				spew.Dump(h.Unknown)
			}
			conn.Close()
			continue
		}

		fmt.Printf("Protocol Version: %d\n", h.Version)
		fmt.Printf("Source Address: %s\n", h.Source.String())
		fmt.Printf("Source Address: %s\n", h.Destination.String())

		// Your application is now free to read the remainder of the content
		o, err := ioutil.ReadAll(conn)
		if e != nil {
			panic(err)
		}
		fmt.Printf("Your App Data: %X\n", o)
	}
}
```
