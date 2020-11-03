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
	events, err := k.Perform(req, func(r *kvika.Recorder, buf []byte) {
		r.Record("received", nil)
	})
	if err != nil {
		os.Exit(-1)
	}
	fmt.Println("record: ", events)
}
