package heartbleed

import (
	"net"
	"testing"
)

func TestBleedIIS(t *testing.T) {
	// IIS (?)
	tgt := Target{"twitch.tv", "https"}
	_, err := Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	if err != Closed {
		t.Errorf("twitch.tv: %v", err)
	}
}

func TestBleedELB(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping ELB test to save time.")
	}

	// ELB
	tgt := Target{"www.theneeds.com", "https"}
	_, err := Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	if err != Safe {
		t.Errorf("www.theneeds.com: %v", err)
	}
}

func TestBleedSafe(t *testing.T) {
	// SAFE
	tgt := Target{"gmail.com", "https"}
	_, err := Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	if err != Safe {
		t.Errorf("gmail.com: %v", err)
	}
}

func TestBleedVulnerable(t *testing.T) {
	// VULNERABLE
	tgt := Target{"www.cloudflarechallenge.com", "https"}
	_, err := Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	if err != nil {
		t.Errorf("www.cloudflarechallenge.com: %v", err)
	}
}

func TestBleedTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test to save time.")
	}

	// TIMEOUT
	tgt := Target{"www.cloudflarechallenge.com:4242", "https"}
	_, err := Heartbleed(&tgt, []byte("FiloSottile/Heartbleed"), false)
	nerr, ok := err.(*net.OpError)
	if !ok || nerr.Err.Error() != "i/o timeout" {
		t.Errorf("www.cloudflarechallenge.com:4242: %v", err)
	}
}
