## Proxy Protocol
This is a golang implementation of the proxy protocol as described
[here](https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt)

This is a prerequisite to implementing an SMTP server in golang, as the SMTP
server must be able to identify the source ip and port of the connecting
client. Load balancers that which to pass on this information implement the
proxy protocol when proxying requests downstream. AWS load balancers support
V1 and V2 across various products to maximize flexability both protocols are
supported but V2 is recommended.

This implementation is focused on performance and completeness.

Current benchmarks for V1 implemenation
```
pkg: github.com/mailgun/proxyproto
BenchmarkReadHeaderV1/TCP4-Minimal-8             2916885               412 ns/op
BenchmarkReadHeaderV1/TCP4-Maximal-8             1306275               920 ns/op
BenchmarkReadHeaderV1/TCP4-Typical-8             1477530               813 ns/op
BenchmarkReadHeaderV1/TCP6-Minimal-8             3085143               383 ns/op
BenchmarkReadHeaderV1/TCP6-Maximal-8              569586              1993 ns/op
```
