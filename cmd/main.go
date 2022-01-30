package main

import (
	"fmt"
	"github.com/artsafin/coda-schema-generator/internal/api"
	"github.com/artsafin/coda-schema-generator/internal/config"
	"github.com/artsafin/coda-schema-generator/internal/generator"
	"io"
	"os"
)

func main() {
	opts, err := config.ParseArgs()

	if err != nil {
		if err.Error() != "" {
			os.Stderr.WriteString(err.Error())
			os.Stderr.WriteString("\n")
		}
		os.Exit(1)
	}

	coda := api.NewClient(opts.APIOptions)

	var tables api.TableList
	if tables, err = coda.LoadTables(); err != nil {
		panic(err)
	}

	var formulas api.EntityList
	if formulas, err = coda.LoadFormulas(); err != nil {
		panic(err)
	}

	var controls api.EntityList
	if controls, err = coda.LoadControls(); err != nil {
		panic(err)
	}

	var columns map[string]api.TableColumns
	if columns, err = coda.LoadColumns(tables); err != nil {
		panic(err)
	}

	var outputWriter io.Writer = os.Stdout
	if opts.OutputFile != "-" {
		outputWriter, err = os.OpenFile(opts.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(fmt.Errorf("unable to open %s for writing", opts.OutputFile))
		}
	}

	gen := generator.NewGenerator(opts.PackageName, tables, columns, formulas, controls)
	err = gen.Generate(outputWriter)

	if err != nil {
		panic(err)
	}
}
