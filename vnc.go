package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/scr34m/go-kvm-ui/domain"
)

var upgrader websocket.Upgrader

func wsh(w http.ResponseWriter, r *http.Request) {
	// try upgrade connection
	useBinary := false
	responseHeader := make(http.Header)
	requestedSubprotocols := websocket.Subprotocols(r)

	if len(requestedSubprotocols) > 0 {
		// pick base64 or binary subprotocol if available
		// otherwise arbitrarily pick the first one and hope for the best
		pickedSubprotocol := ""
		for _, subprotocol := range requestedSubprotocols {
			if subprotocol == "base64" {
				pickedSubprotocol = "base64"
				break
			} else if subprotocol == "binary" {
				pickedSubprotocol = "binary"
				useBinary = true
				break
			}
		}
		if pickedSubprotocol == "" {
			pickedSubprotocol = requestedSubprotocols[0]

			log.Printf("Warning: client %s did not offer base64 or binary subprotocols, falling back to %s", r.RemoteAddr, pickedSubprotocol)
		}
		responseHeader.Set("Sec-Websocket-Protocol", pickedSubprotocol)
	}

	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		log.Printf("websockify error (%s): %s", r.RemoteAddr, err.Error())
		return
	}
	defer conn.Close()

	vmname := r.URL.Path[len("/websockify/"):]

	// TODO check already exsist
	domain := domain.Load(vmname)

	target := fmt.Sprintf("127.0.0.1:%d", domain.GetVNC())

	log.Printf("Initializing connection from %s to %s", r.RemoteAddr, target)

	sock, err := net.Dial("tcp", target)
	if err != nil {
		log.Print(err)
		return
	}
	defer sock.Close()

	done := make(chan bool, 2)
	go func() {
		defer func() {
			done <- true
		}()

		wbuf := make([]byte, 32*1024)
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var fwdbuf []byte

			if messageType == websocket.TextMessage {
				n, _ := base64.StdEncoding.Decode(wbuf, p)
				fwdbuf = wbuf[:n]
			} else if messageType == websocket.BinaryMessage {
				fwdbuf = p
			}

			if fwdbuf != nil {
				_, err = sock.Write(fwdbuf)
				if err != nil {
					return
				}
			}
		}
	}()
	go func() {
		defer func() {
			done <- true
		}()

		rbuf := make([]byte, 8192)
		wbuf := make([]byte, len(rbuf)*2)
		for {
			n, err := sock.Read(rbuf)
			if err != nil {
				return
			}

			if n > 0 {
				var err error

				if useBinary {
					err = conn.WriteMessage(websocket.BinaryMessage, rbuf[:n])
				} else {
					base64.StdEncoding.Encode(wbuf, rbuf[:n])
					err = conn.WriteMessage(websocket.TextMessage, wbuf[:base64.StdEncoding.EncodedLen(n)])
				}

				if err != nil {
					return
				}
			}
		}
	}()
	<-done
}
