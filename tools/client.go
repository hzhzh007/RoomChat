// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

var concurency = flag.Int("c", 30, "http client cocurency")
var print_ = flag.Bool("p", false, "print ")
var path = flag.String("path", "/ws", "websocket path")

var dialer = websocket.Dialer{} // use default options

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: *path}
	log.Printf("connecting to %s", u.String())

	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	for i := 0; i < *concurency; i++ {
		log.Println(" processs client :", i)
		go func(i int) {
			c, _, err := dialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal("dial:", err)
			}
			defer c.Close()
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break
				}
				if *print_ {
					log.Printf("%d, recv: %s", i, message)
				}
			}
		}(i)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
