package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
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

var output = kingpin.
	Flag("output", "output filename").
	Short('o').
	Default("").
	String()

func perform(u *url.URL) {
	log := zap.L()
	k := kvika.New()
	req := &kvika.Request{
		Method: "GET",
		URL:    u,
	}
	resp, err := k.Perform(req, func(r *kvika.Recorder, buf []byte) {
		r.Record("data-received", buf)
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
		switch p := ev.Payload.(type) {
		case []byte:
			fmt.Printf("[% 3d/%0.3fms] %s: %v\n", i, ev.At, ev.Name, fmt.Sprintf("%d bytes", len(p)))
			if len(*output) > 0 {
				filename := fmt.Sprintf("%s.%03d", *output, i)
				err := ioutil.WriteFile(filename, p, 0644)
				if err != nil {
					log.Fatal("Failed to write file", zap.String("file-name", filename), zap.Error(err))
					os.Exit(-1)
				}
			}
		default:
			fmt.Printf("[% 3d/%0.3fms] %s\n", i, ev.At, ev.Name)
			break
		}
	}
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

	perform(*urlArg)
}
