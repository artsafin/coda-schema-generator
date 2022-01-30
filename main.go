package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/artsafin/coda-schema-generator/dto"
	"github.com/artsafin/coda-schema-generator/schema"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	opts, err := parseArgs()

	if err != nil {
		if err.Error() != "" {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		os.Exit(1)
	}

	var outputWriter io.Writer = os.Stdout
	if opts.OutputFile != "-" {
		outputWriter, err = os.OpenFile(opts.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to open %s for writing\n", opts.OutputFile)
			os.Exit(1)
		}
	}

	err = schema.WriteDoc(opts.APIOptions, opts.PackageName, outputWriter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() (dto.Options, error) {
	isVerbose := flag.Bool("verbose", false, "Verbose or not")
	isHelp := flag.Bool("help", false, "Print help")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "%s [OPTIONS] <CODA API TOKEN> <CODA DOCUMENT ID>\n\nOptions:\n", os.Args[0])

		flag.PrintDefaults()
	}

	if flag.NArg() != 2 {
		flag.Usage()
		return dto.Options{}, errors.New("")
	}

	if *isHelp {
		flag.Usage()
		return dto.Options{}, errors.New("")
	}

	endpoint := strings.TrimSpace(os.Getenv("CODA_API_ENDPOINT"))
	if endpoint == "" {
		endpoint = "https://coda.io/apis/v1"
	}

	return dto.Options{
		APIOptions: dto.APIOptions{
			Verbose:        *isVerbose, // TODO: expose this as program argument
			Endpoint:       endpoint,
			Token:          flag.Arg(0),
			DocID:          flag.Arg(1),
			RequestTimeout: time.Second * 15, // TODO: expose this as program argument
		},
		DumpOptions: dto.DumpOptions{
			OutputFile:  "-",          // TODO: expose this as program argument
			PackageName: "codaschema", // TODO: expose this as program argument
		},
	}, nil
}
