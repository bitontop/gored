package okex

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
	"github.com/gorilla/websocket"

	"log"
	"net/url"
)

// Socket get orderbook from websocket
func Socket() { //pair *pair.Pair
	u := url.URL{Scheme: "wss", Host: "real.okex.com:10442", Path: "/ws/v3", RawQuery: "compress=true"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// eOkex :=
	// e.GetSymbolByPair(pair.Base)

	// c.WriteMessage(websocket.TextMessage, []byte(`{"channel":"ok_sub_futureusd_btc_depth_quarter","event":"addChannel"}`))
	// c.WriteMessage(websocket.TextMessage, []byte(`{"op": "subscribe", "args": ["spot/ticker:ETH-USDT"]}`))
	c.WriteMessage(websocket.TextMessage, []byte(`{"op": "subscribe", "args": ["spot/depth:ETH-BTC"]}`))
	// c.WriteMessage(websocket.TextMessage, []byte(`{"op": "subscribe", "args": ["spot/depth5:ETH-USDT"]}`))

	done := make(chan struct{})

	go func() {
		defer close(done)
		// i := 0
		for {
			messageType, message, err := c.ReadMessage()
			switch messageType {
			case websocket.TextMessage:
				// no need uncompressed
				// log.Printf("recv Text: %s", message)
				log.Printf("Orderbook: %v", writeOrderBook(message))
			case websocket.BinaryMessage:
				// uncompressed
				text, err := GzipDecode(message)
				if err != nil {
					log.Println("err", err)
				} else {
					// log.Printf("recv Bin: %s", text)
					log.Printf("Orderbook: %v", writeOrderBook(text))
				}
			}
			if err != nil {
				log.Println("read:", err)
				return
			}
			/* if i == 10 {
				c.WriteMessage(websocket.TextMessage, []byte(`{"op": "unsubscribe", "args": ["spot/depth:ETH-BTC"]}`))
				c.WriteMessage(websocket.TextMessage, []byte(`{"op": "subscribe", "args": ["spot/depth:ETH-BTC"]}`))
			}
			i++ */

		}
	}()

	select {}

}

// GzipDecode decode the message from socket
func GzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)

}

func writeOrderBook(text []byte) *exchange.Maker {
	orderBook := WSOrderBook{}
	// symbol := e.GetSymbolByPair(pair)

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.Timestamp = float64(time.Now().UnixNano() / 1e6)

	if err := json.Unmarshal(text, &orderBook); err != nil {
		return nil //fmt.Sprintf("Okex WebSocket Get Orderbook Json Unmarshal Err: %v \n%v", err, text)
	} else if len(orderBook.Data) == 0 {
		return nil
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Data[0].Bids {
		var buydata exchange.Order
		buydata.Rate, _ = strconv.ParseFloat(bid[0], 64)
		buydata.Quantity, _ = strconv.ParseFloat(bid[1], 64)

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Data[0].Asks {
		var selldata exchange.Order
		selldata.Rate, _ = strconv.ParseFloat(ask[0], 64)
		selldata.Quantity, _ = strconv.ParseFloat(ask[1], 64)

		maker.Asks = append(maker.Asks, selldata)
	}

	return maker
}
