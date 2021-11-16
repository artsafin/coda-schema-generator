package main

import (
	"coda-schema-generator/internal/config"
	"coda-schema-generator/internal/dto"
	"coda-schema-generator/internal/generator"
	"fmt"
	"io"
	"os"
)

func main() {
	opts, err := config.ParseArgs()

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
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
