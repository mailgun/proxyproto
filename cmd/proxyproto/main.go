package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mailgun/proxyproto"
	"github.com/pkg/errors"
	"github.com/thrawn01/args"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	parser := args.NewParser(args.Name("proxyproto"),
		args.Desc("Proxy Protocol CLI"))

	parser.AddCommand("serve", serve)
	//parser.AddCommand("client-v1", client)
	//parser.AddCommand("client-v2", client)
	parser.AddCommand("client-plain", plainClient)

	// Run the command chosen by the user
	retCode, err := parser.ParseAndRun(nil, nil)
	if err != nil {
		fmt.Printf("[ERR] %s\n", err)
	}
	os.Exit(retCode)
}

func plainClient(p *args.ArgParser, d interface{}) (int, error) {
	var network, address, msg string

	p.AddOption("network").
		Help("The which golang network to connect on").
		StoreString(&network).
		Choices([]string{"tcp6", "tcp4", "tcp"}).
		Alias("-n").
		Default("tcp")

	p.AddOption("address").
		Help("The ip:port address to connect to").
		StoreString(&address).
		Alias("-a").
		Default("127.0.0.1:2319")

	p.AddOption("message").
		Help("The message to send to the TCP server").
		StoreString(&msg).
		Alias("-m").
		Default(" \n")

	if _, err := p.Parse(nil); err != nil {
		return 1, err
	}

	conn, err := net.Dial(network, address)
	if err != nil {
		return 1, errors.Wrapf(err, "while dialing '%s'", address)
	}

	if _, err := fmt.Fprint(conn, msg); err != nil {
		return 1, errors.Wrap(err, "while writing message to socket")
	}
	/*if _, err := io.Copy(os.Stdin, conn); err != nil {
		return 1, errors.Wrap(err, "while copying io")
	}*/
	conn.Close()
	return 0, nil
}

func serve(p *args.ArgParser, d interface{}) (int, error) {
	var network, port string

	p.AddOption("network").
		Help("The which golang network to listen on").
		StoreString(&network).
		Choices([]string{"tcp6", "tcp4", "tcp"}).
		Alias("-n").
		Default("tcp")

	p.AddOption("port").
		Help("The which port to listen on").
		StoreString(&port).
		Alias("-p").
		Default(":2319")

	if _, err := p.Parse(nil); err != nil {
		return 1, err
	}

	ln, err := net.Listen(network, port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Serve '%s' on: %s....\n", network, port)

	for {
		conn, e := ln.Accept()
		fmt.Printf("New Connection: %s\n", conn.RemoteAddr().String())
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				fmt.Printf("[ERR] Accept temporary: %s", ne)
				continue
			}

			fmt.Printf("[ERR] Accept: %s", e)
			return 1, nil
		}

		h, err := proxyproto.ReadHeader(conn)
		if err != nil {
			fmt.Printf("[ERR] Proxy protocol: %s\n", err)
			conn.Close()
			continue
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
		fmt.Printf("Destination Address: %s\n", h.Destination.String())

		// Your application is now free to read the remainder of the content
		o, err := ioutil.ReadAll(conn)
		if e != nil {
			fmt.Printf("[ERR] ReadAll: %s\n", err)
			conn.Close()
			continue
		}
		fmt.Printf("Trailing Data: %X\n", o)
		fmt.Printf("Connection Closed\n")
		conn.Close()
	}
}
