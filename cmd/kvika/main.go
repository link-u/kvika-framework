package main

import (
	"fmt"
	"os"

	"github.com/link-u/kvika-framework/pkg/kvika"
	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "v0.1.0"

/*****************************************************************************
  Flags
 *****************************************************************************/

var standardLog = kingpin.
	Flag("cli.standard-log", "Print logs in standard format, not in json").
	Default("false").
	Bool()

func perform() {

}

func main() {
	var err error
	var log *zap.Logger

	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	urlArg := kingpin.Arg("URL", "URL to observe").Required().URL()
	kingpin.Parse()

	// Check weather terminal or not
	if *standardLog || isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create logger: %v", err)
		os.Exit(-1)
	}
	undo := zap.ReplaceGlobals(log)
	defer undo()

	k := kvika.New()
	req := &kvika.Request{
		Method: "GET",
		URL:    *urlArg,
	}
	resp, err := k.Perform(req, func(r *kvika.Recorder, buf []byte) {
		r.Record("data-received", fmt.Sprintf("%d bytes", len(buf)))
	})
	if err != nil {
		log.Fatal("Failed to perform request", zap.String("url", req.URL.String()), zap.Error(err))
		os.Exit(-1)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Failed to perform request",
			zap.String("url", req.URL.String()), zap.Int("status-code", resp.StatusCode),
			zap.Error(err))
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
