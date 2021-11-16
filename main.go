package main

import (
	"coda-schema-generator/dto"
	"coda-schema-generator/generator"
	"fmt"
	"io"
	"os"
)

func main() {
	opts, err := parseArgs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	coda := NewClient(opts.APIOptions)

	var tables dto.Tables
	if tables, err = coda.loadTables(); err != nil {
		panic(err)
	}

	var columns map[string]dto.TableColumns
	if columns, err = coda.loadColumns(tables); err != nil {
		panic(err)
	}

	var outputWriter io.Writer
	outputWriter = os.Stdout
	if opts.OutputFile != "-" {
		outputWriter, err = os.OpenFile(opts.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(fmt.Errorf("unable to open %s for writing", opts.OutputFile))
		}
	}

	gen := generator.NewGenerator(opts.PackageName, tables, columns)
	err = gen.Generate(outputWriter)

	if err != nil {
		panic(err)
	}
}
