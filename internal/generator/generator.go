package generator

import (
	"github.com/artsafin/coda-schema-generator/internal/api"
	"github.com/artsafin/coda-schema-generator/internal/templates"
	"io"
	"text/template"
)

type generator struct {
	Tables      api.TableList
	Formulas    api.EntityList
	Controls    api.EntityList
	Columns     map[string]api.TableColumns
	DTO         *FieldMapper
	PackageName string
	name        nameConverter
}

func NewGenerator(
	packageName string,
	tables api.TableList,
	columns map[string]api.TableColumns,
	formulas api.EntityList,
	controls api.EntityList,
) generator {
	nc := NewNameConverter()
	fm := NewFieldMapper(nc)

	for _, t := range tables.Items {
		for _, c := range columns[t.ID].Items {
			fm.registerField(t.ID, c)
		}
	}

	return generator{
		PackageName: packageName,
		Tables:      tables,
		Columns:     columns,
		DTO:         fm,
		Formulas:    formulas,
		Controls:    controls,
		name:        nc,
	}
}

func (d *generator) Generate(w io.Writer) error {
	tpl := template.New("main.go.tmpl").
		Funcs(template.FuncMap{
			"fieldName": d.name.ConvertNameToGoSymbol,
			"typeName":  d.name.ConvertNameToGoType,
		})

	tpl, err := tpl.ParseFS(templates.FS, "*.tmpl", "types/*.tmpl")
	if err != nil {
		return err
	}

	err = tpl.ExecuteTemplate(w, "main.go.tmpl", d)
	if err != nil {
		return err
	}

	dirs, err := templates.FS.ReadDir("types")
	if err != nil {
		return err
	}
	for _, typeFile := range dirs {
		if typeFile.IsDir() {
			continue
		}
		err = tpl.ExecuteTemplate(w, typeFile.Name(), d)
		if err != nil {
			return err
		}
	}

	err = tpl.ExecuteTemplate(w, "dto.go.tmpl", d)
	if err != nil {
		return err
	}

	return nil
}
