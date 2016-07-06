//
// Copyright © 2013-2016 Guy M. Allard
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
Subscribe and receive messages from a STOMP broker.

	Examples:

		# Subscribe to a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.1
		go run subscribe.go

		# Subscribe to a broker using STOMP protocol level 1.2:
		STOMP_PROTOCOL=1.2 go run subscribe.go

		# Subscribe to a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run subscribe.go

		# Subscribe to a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run subscribe.go

		# Subscribe to a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run subscribe.go

*/
package main

import (
	"log"
	"net"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var exampid = "subscribe: "

// Connect to a STOMP broker, subscribe and receive some messages and disconnect.
func main() {
	log.Println(exampid + "starts ...")

	// Set up the connection.
	h, p := sngecomm.HostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "dial complete ...")
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid+"stomp connect complete ...", conn.Protocol())

	pbc := sngecomm.Pbc() // Print byte count

	log.Println(exampid+"connected headers", conn.ConnectResponse.Headers)
	// *NOTE* your application functionaltiy goes here!
	// With Stomp, you must SUBSCRIBE to a destination in order to receive.
	// Subscribe returns a channel of MessageData struct.
	// Here we use a common utility routine to handle the differing subscribe
	// requirements of each protocol level.
	d := sngecomm.Dest()
	i := stompngo.Uuid()
	r := sngecomm.Subscribe(conn, d, i, "auto")
	log.Println(exampid + "stomp subscribe complete ...")
	// Read data from the returned channel
	for i := 1; i <= sngecomm.Nmsgs(); i++ {
		m := <-r
		log.Println(exampid + "channel read complete ...")
		log.Println("Message Number:", i)
		// MessageData has two components:
		// a) a Message struct
		// b) an Error value.  Check the error value as usual
		if m.Error != nil {
			log.Fatalln(m.Error) // Handle this
		}
		//
		log.Printf("Frame Type: %s\n", m.Message.Command) // Will be MESSAGE or ERROR!
		if m.Message.Command != stompngo.MESSAGE {
			log.Fatalln(m) // Handle this ...
		}
		h := m.Message.Headers
		for j := 0; j < len(h)-1; j += 2 {
			log.Printf("Header: %s:%s\n", h[j], h[j+1])
		}
		if pbc > 0 {
			maxlen := pbc
			if len(m.Message.Body) < maxlen {
				maxlen = len(m.Message.Body)
			}
			ss := string(m.Message.Body[0:maxlen])
			log.Printf("Payload: %s\n", ss) // Data payload
		}
	}
	// It is polite to unsubscribe, although unnecessary if a disconnect follows.
	// Again we use a utility routine to handle the different protocol level
	// requirements.
	sngecomm.Unsubscribe(conn, d, i)
	log.Println(exampid + "unsubscribe complete")

	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "network close complete ...")

	log.Println(exampid + "ends ...")
}
