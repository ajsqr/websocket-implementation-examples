package main

import (
	"fmt"
	"net"
	"bufio"
	"crypto/sha1"
	"encoding/base64"

	http "net/http"

)

const (
	// websocketGUID is the GUID specified in RFC 6455
	websocketGUID =  "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

func main() {
	listener, err := net.Listen("tcp", ":8123")
	if err != nil{
		panic(err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil{
			panic(err.Error())
		}

		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)

		r, err := http.ReadRequest(reader)
		if err != nil{
			panic(err.Error())
		}


		resp := newAcceptResponse(r)
		resp.Write(writer)
		writer.Flush()
	}
}



func newAcceptResponse(r *http.Request) *http.Response {
	websocketKey := r.Header.Get("Sec-WebSocket-Key")
	acceptToken := generateWebsocketAcceptToken(websocketKey)
	resp := http.Response{
		Proto : "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		StatusCode: 101,
		Header: http.Header{},
	}

	resp.Header.Set("Upgrade", "websocket")
	resp.Header.Set("Connection", "Upgrade")
	resp.Header.Set("Sec-WebSocket-Accept", acceptToken)

	return &resp
}

func generateWebsocketAcceptToken(secWebsocketKey string) string {
	combinedKey := []byte(fmt.Sprintf("%s%s", secWebsocketKey, websocketGUID))
	hash := sha1.Sum(combinedKey)
	encodedKey := base64.StdEncoding.EncodeToString(hash[:])
	return encodedKey
}