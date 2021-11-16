package main

import (
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

func parseArgs() (Options, error) {
	if len(os.Args) != 3 {
		return Options{}, fmt.Errorf("usage: %s <CODA API TOKEN> <CODA DOCUMENT ID>", os.Args[0])
	}

	endpoint := strings.TrimSpace(os.Getenv("CODA_API_ENDPOINT"))
	if endpoint == "" {
		endpoint = "https://coda.io/apis/v1"
	}

	return Options{
		APIOptions: APIOptions{
			Verbose:  false, // TODO: expose this as program argument
			Endpoint: endpoint,
			Token:    os.Args[1],
			DocID:    os.Args[2],
		},
		DumpOptions: DumpOptions{
			OutputFile:  "-",          // TODO: expose this as program argument
			PackageName: "codaschema", // TODO: expose this as program argument
		},
	}, nil
}
