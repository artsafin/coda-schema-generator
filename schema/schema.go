package schema

import (
	"github.com/artsafin/coda-schema-generator/dto"
	"github.com/artsafin/coda-schema-generator/internal/api"
	"github.com/artsafin/coda-schema-generator/internal/generator"
	"io"
)

func WriteDoc(opts dto.APIOptions, pkgName string, wr io.Writer) (err error) {
	coda := api.NewClient(opts)

	var tables dto.TableList
	if tables, err = coda.LoadTables(); err != nil {
		return err
	}

	var formulas dto.EntityList
	if formulas, err = coda.LoadFormulas(); err != nil {
		return err
	}

	var controls dto.EntityList
	if controls, err = coda.LoadControls(); err != nil {
		return err
	}

	var columns map[string]dto.TableColumns
	if columns, err = coda.LoadColumns(tables); err != nil {
		return err
	}

	gen := generator.NewGenerator(pkgName, tables, columns, formulas, controls)

	return gen.Generate(wr)
}
