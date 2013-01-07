//
// Copyright © 2012 Guy M. Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

/*
Connect and Disconnect from a STOMP 1.1 broker with a TLS connection.
*/
package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gmallard/stompngo"
	. "github.com/gmallard/stompngo_examples/sngecomm"
	"log"
	"net"
)

var (
	exampid    = "connect_ssl_11:"
	testConfig *tls.Config
)

// Connect to a STOMP 1.1 broker using TLS and disconnect.
func main() {
	fmt.Println(exampid, "starts ...")

	// TLS Configuration.  This configuration assumes that:
	// a) The server used does *not* require client certificates
	// b) This client has no need to authenticate the server
	// Note that the tls.Config structure can be modified to support any
	// authentication scheme, including two-way/mutual authentication
	testConfig = new(tls.Config)
	testConfig.InsecureSkipVerify = true // Do *not* check the server's certificate

	// Open a net connection
	h, p := HostAndPort11()

	// Alternative connect logic: use net.Dial and tls.Client
	t, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "dial complete ...")
	n := tls.Client(t, testConfig)
	e = n.Handshake()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}

	// All stomp API methods require 'Headers'.  Stomp headers are key/value 
	// pairs.  The stompngo package implements them using a string slice.
	//
	ch := stompngo.Headers{"accept-version", "1.1",
		"host", Vhost()}

	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the (empty) Headers
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "stomp connect complete ...")

	// *NOTE* your application functionaltiy goes here!

	// Polite Stomp disconnects are not required, but highly recommended.
	// Empty headers again here.
	eh := stompngo.Headers{}
	e = conn.Disconnect(eh)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "stomp disconnect complete ...")

	// Close the net connection.
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "network close complete ...")

	fmt.Println(exampid, "ends ...")
}
