package main

import (
	"fmt"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
)

const (
	// websocketGUID is the GUID specified in RFC 6455
	websocketGUID =  "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

func main() {
	ws := wsh{}
	http.ListenAndServe(":8123", &ws)
}

type wsh struct {}

func (ws *wsh) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, rw, err := hijacker.Hijack()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := newAcceptResponse(r)
	resp.Write(rw)
	rw.Flush()
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