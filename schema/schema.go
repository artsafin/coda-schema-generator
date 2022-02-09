package schema

import (
	"github.com/artsafin/coda-schema-generator/dto"
	"github.com/artsafin/coda-schema-generator/internal/api"
)

func Get(opts dto.APIOptions) (s *dto.Schema, err error) {
	coda := api.NewHTTPClient(opts)

	var tables dto.TableList
	if tables, err = coda.LoadTables(); err != nil {
		return nil, err
	}

	var formulas dto.EntityList
	if formulas, err = coda.LoadFormulas(); err != nil {
		return nil, err
	}

	var controls dto.EntityList
	if controls, err = coda.LoadControls(); err != nil {
		return nil, err
	}

	var columns map[string]dto.TableColumns
	if columns, err = coda.LoadColumns(tables); err != nil {
		return nil, err
	}

	return &dto.Schema{
		Tables:   tables,
		Columns:  columns,
		Formulas: formulas,
		Controls: controls,
	}, nil
}
