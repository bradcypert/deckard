package cmd

import (
	"testing"
	"time"
)

func TestMakeTimestamp(t *testing.T) {
	timestamp := makeTimestamp(time.Date(2000, 4, 17, 4, 44, 44, 0, time.UTC))
	if timestamp != "955946684000" {
		t.Error("Expected to find:", 0, "instead got:", timestamp)
	}
}