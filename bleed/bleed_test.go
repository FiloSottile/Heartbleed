package heartbleed

import (
	"net"
	"testing"
)

func TestBleed(t *testing.T) {
	// IIS (?)
	tgt := Target{"twitch.tv", "https"}
	_, err := Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	if err != Closed {
		t.Errorf("twitch.tv: %v", err)
	}

	// ELB
	tgt = Target{"www.theneeds.com", "https"}
	_, err = Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	if err != Safe {
		t.Errorf("www.theneeds.com: %v", err)
	}

	// SAFE
	tgt = Target{"gmail.com", "https"}
	_, err = Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	if err != Safe {
		t.Errorf("gmail.com: %v", err)
	}

	// VULNERABLE
	tgt = Target{"www.cloudflarechallenge.com", "https"}
	_, err = Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	if err != nil {
		t.Errorf("www.cloudflarechallenge.com: %v", err)
	}

	// TIMEOUT
	tgt = Target{"www.cloudflarechallenge.com:4242", "https"}
	_, err = Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	nerr, ok := err.(*net.OpError)
	if !ok || nerr.Err.Error() != "i/o timeout" {
		t.Errorf("www.cloudflarechallenge.com:4242: %v", err)
	}
}
