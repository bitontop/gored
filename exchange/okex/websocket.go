package okex

import (
	"bytes"
	"compress/flate"
	"io/ioutil"

	"github.com/gorilla/websocket"

	"log"
	"net/url"
)

func main() {
	u := url.URL{Scheme: "wss", Host: "real.okex.com:10441", Path: "/websocket", RawQuery: "compress=true"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	c.WriteMessage(websocket.TextMessage, []byte(`{"channel":"ok_sub_futureusd_btc_depth_quarter","event":"addChannel"}`))

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			switch messageType {
			case websocket.TextMessage:
				// no need uncompressed
				log.Printf("recv: %s", message)
			case websocket.BinaryMessage:
				// uncompressed
				text, err := GzipDecode(message)
				if err != nil {
					log.Println("err", err)
				} else {
					log.Printf("recv: %s", text)
				}
			}
			if err != nil {
				log.Println("read:", err)
				return
			}

		}
	}()

	select {}

}

func GzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)

}
