package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type APIOptions struct {
	Verbose  bool
	Endpoint string
	Token    string
	DocID    string
}

type DumpOptions struct {
	OutputFile  string
	PackageName string
}

type Options struct {
	APIOptions
	DumpOptions
}

func ParseArgs() (Options, error) {
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
		return Options{}, errors.New("")
	}

	if *isHelp {
		flag.Usage()
		return Options{}, errors.New("")
	}

	endpoint := strings.TrimSpace(os.Getenv("CODA_API_ENDPOINT"))
	if endpoint == "" {
		endpoint = "https://coda.io/apis/v1"
	}

	return Options{
		APIOptions: APIOptions{
			Verbose:  *isVerbose, // TODO: expose this as program argument
			Endpoint: endpoint,
			Token:    flag.Arg(0),
			DocID:    flag.Arg(1),
		},
		DumpOptions: DumpOptions{
			OutputFile:  "-",          // TODO: expose this as program argument
			PackageName: "codaschema", // TODO: expose this as program argument
		},
	}, nil
}
