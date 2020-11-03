package main

import (
	"fmt"
	"os"

	"github.com/link-u/kvika-framework/pkg/kvika"
)

func main() {
	k := kvika.New()
	req := &kvika.Request{
		Method: "GET",
		URL:    "https://www.google.com/",
	}
	resp, err := k.Perform(req, func(r *kvika.Recorder, buf []byte) {
		r.Record("data-received", fmt.Sprintf("%d bytes", len(buf)))
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to perform request: %v", err)
		os.Exit(-1)
	}
	if resp.StatusCode != 200 {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to perform request: status=%v", resp.StatusCode)
		os.Exit(-1)
	}
	for i, ev := range resp.Events {
		if ev.Payload != nil {
			fmt.Printf("[%02d/%0.3fms] %s: %v\n", i, ev.At, ev.Name, ev.Payload)
		} else {
			fmt.Printf("[%02d/%0.3fms] %s\n", i, ev.At, ev.Name)
		}
	}
}
