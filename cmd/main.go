package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/artsafin/coda-schema-generator/dto"
	"github.com/artsafin/coda-schema-generator/generator"
	"github.com/artsafin/coda-schema-generator/schema"
	"io"
	"os"
	"strings"
	"time"
)

func fatal(format string, vals ...interface{}) {
	fmt.Fprintf(os.Stderr, format, vals...)
	os.Exit(1)
}

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
			fatal("error: unable to open %s for writing\n", opts.OutputFile)
		}
	}

	sch, err := schema.Get(opts.APIOptions)
	if err != nil {
		fatal("error: %v\n", err)
	}

	err = generator.Generate(opts.PackageName, sch, outputWriter)
	if err != nil {
		fatal("error: %v\n", err)
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
