package generator

import (
	"coda-schema-generator/internal/dto"
	"coda-schema-generator/internal/templates"
	"io"
	"text/template"
)

type generator struct {
	Tables      dto.Tables
	Columns     map[string]dto.TableColumns
	PackageName string
	name        nameConverter
}

func NewGenerator(packageName string, tables dto.Tables, columns map[string]dto.TableColumns) generator {
	return generator{
		PackageName: packageName,
		Tables:      tables,
		Columns:     columns,
		name:        NewNameConverter(),
	}
}

func (d *generator) Generate(w io.Writer) error {
	tpl := template.New("top").
		Funcs(template.FuncMap{
			"fieldName": d.name.ConvertNameToGoSymbol,
			"typeName":  d.name.ConvertNameToGoType,
		})

	tpl, err := tpl.ParseFS(templates.FS, "*.tmpl")
	if err != nil {
		return err
	}

	err = tpl.ExecuteTemplate(w, "main.go.tmpl", d)

	if err != nil {
		return err
	}

	return nil
}
