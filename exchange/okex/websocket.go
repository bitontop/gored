package okex

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/exchange/okex"
	"github.com/bitontop/gored/pair"
	"github.com/bradfitz/slice"
	"github.com/gorilla/websocket"
	"github.com/tony0408/Coinbene_json/coin"
	"github.com/tony0408/Coinbene_json/test/conf"

	"log"
	"net/url"
)

func InitOkex() exchange.Exchange {
	coin.Init()
	pair.Init()

	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	conf.Exchange(exchange.OKEX, config)

	ex := okex.CreateOkex(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}

// Socket get orderbook from websocket
func Socket(pair *pair.Pair) {
	u := url.URL{Scheme: "wss", Host: "real.okex.com:10442", Path: "/ws/v3", RawQuery: "compress=true"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	msg := fmt.Sprintf(`{"op": "subscribe", "args": ["spot/depth:%v-%v"]}`, pair.Target.Code, pair.Base.Code)

	// e := InitOkex()
	// e.GetSymbolByPair(pair)

	// c.WriteMessage(websocket.TextMessage, []byte(`{"channel":"ok_sub_futureusd_btc_depth_quarter","event":"addChannel"}`))
	// c.WriteMessage(websocket.TextMessage, []byte(`{"op": "subscribe", "args": ["spot/ticker:ETH-USDT"]}`))
	// c.WriteMessage(websocket.TextMessage, []byte(`{"op": "subscribe", "args": ["spot/depth:ETH-BTC"]}`))
	c.WriteMessage(websocket.TextMessage, []byte(msg))

	done := make(chan struct{})

	okexMaker := &exchange.Maker{}

	go func() {
		defer close(done)
		// i := 0
		for {
			messageType, message, err := c.ReadMessage()
			switch messageType {
			case websocket.TextMessage:
				// no need uncompressed
				// log.Printf("recv Text: %s", message)
				// log.Printf("Orderbook: %v", writeOrderBook(message, okexMaker))
				log.Printf("Orderbook: %v", writeOrderBook(message, okexMaker))
				maker := writeOrderBook(message, okexMaker)
				if maker != nil {
					log.Printf("Orderbook Bids: %v", maker.Bids[:min(len(maker.Bids), 5)])
					log.Printf("Orderbook Asks: %v", maker.Asks[:min(len(maker.Asks), 5)])
				}
			case websocket.BinaryMessage:
				// uncompressed
				text, err := GzipDecode(message)
				if err != nil {
					log.Println("err:", err)
				} else {
					// log.Printf("recv Bin: %s", text)
					// log.Printf("Orderbook: %v", writeOrderBook(text, okexMaker))
					maker := writeOrderBook(text, okexMaker)
					if maker != nil {
						fmt.Printf("Orderbook Bids: %v\n", maker.Bids[:min(len(maker.Bids), 5)])
						fmt.Printf("Orderbook Asks: %v\n", maker.Asks[:min(len(maker.Asks), 5)])
						log.Println("")
					}

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

func writeOrderBook(text []byte, maker *exchange.Maker) *exchange.Maker {
	orderBook := WSOrderBook{}

	// maker = &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.Timestamp = float64(time.Now().UnixNano() / 1e6)

	if err := json.Unmarshal(text, &orderBook); err != nil {
		log.Printf("Okex WebSocket Get Orderbook Json Unmarshal Err: %v \n%v", err, text)
		return nil
	} else if len(orderBook.Data) == 0 {
		return nil
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Data[0].Bids {
		var buydata exchange.Order
		buydata.Rate, _ = strconv.ParseFloat(bid[0], 64)
		buydata.Quantity, _ = strconv.ParseFloat(bid[1], 64)

		updated, order := update(maker.Bids, buydata.Rate, buydata.Quantity)
		if updated {
			maker.Bids = order
		} else {
			maker.Bids = append(maker.Bids, buydata)
		}
	}
	for _, ask := range orderBook.Data[0].Asks {
		var selldata exchange.Order
		selldata.Rate, _ = strconv.ParseFloat(ask[0], 64)
		selldata.Quantity, _ = strconv.ParseFloat(ask[1], 64)

		updated, order := update(maker.Asks, selldata.Rate, selldata.Quantity)
		if updated {
			maker.Asks = order
		} else {
			maker.Asks = append(maker.Asks, selldata)
		}
	}

	if maker == nil {
		return nil
	}
	slice.Sort(maker.Bids[:], func(i, j int) bool {
		return maker.Bids[i].Rate > maker.Bids[j].Rate
	})
	slice.Sort(maker.Asks[:], func(i, j int) bool {
		return maker.Asks[i].Rate < maker.Asks[j].Rate
	})

	return maker
}

func update(orders []exchange.Order, rate, quantity float64) (bool, []exchange.Order) {
	for i, order := range orders {
		if order.Rate == rate {
			if quantity == 0 {
				orders = remove(orders, i)
				return true, orders
			}
			orders[i].Quantity = quantity
			return true, orders
		}
	}

	return false, orders
}

func remove(array []exchange.Order, i int) []exchange.Order {
	if len(array) == 1 {
		return array[:0]
	} else if len(array) < 1 {
		return array
	}

	array[len(array)-1], array[i] = array[i], array[len(array)-1]
	return array[:len(array)-1]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
