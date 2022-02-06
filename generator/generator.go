package generator

import (
	"github.com/artsafin/coda-schema-generator/dto"
	"github.com/artsafin/coda-schema-generator/internal/templates"
	"io"
	"text/template"
)

type templateData struct {
	dto.Schema
	DTO               *fieldMapper
	SchemaPackageName string
	APIPackageName    string
	name              nameConverter
}

func Generate(schemaPackageName, apiPackageName string, data *dto.Schema, w io.Writer) error {
	nc := newNameConverter()
	fm := newFieldMapper(nc)

	for _, t := range data.Tables.Items {
		for _, c := range data.Columns[t.ID].Items {
			fm.registerField(t.ID, c)
		}
	}

	tpldata := templateData{
		Schema:            *data,
		DTO:               fm,
		SchemaPackageName: schemaPackageName,
		APIPackageName:    apiPackageName,
		name:              nc,
	}

	tpl := template.New("main.go.tmpl").
		Funcs(template.FuncMap{
			"fieldName": tpldata.name.ConvertNameToGoSymbol,
			"typeName":  tpldata.name.ConvertNameToGoType,
		})

	tpl, err := tpl.ParseFS(templates.FS, "*.tmpl", "types/*.tmpl")
	if err != nil {
		return err
	}

	err = tpl.ExecuteTemplate(w, "main.go.tmpl", tpldata)
	if err != nil {
		return err
	}

	var additionalFiles []string

	dirs, err := templates.FS.ReadDir("types")
	if err != nil {
		return err
	}
	for _, typeFile := range dirs {
		if typeFile.IsDir() {
			continue
		}
		additionalFiles = append(additionalFiles, typeFile.Name())
	}
	additionalFiles = append(additionalFiles, "dto.go.tmpl", "doc.go.tmpl", "doc_load_shallow.go.tmpl", "doc_load_deep.go.tmpl")

	for _, f := range additionalFiles {
		err = tpl.ExecuteTemplate(w, f, tpldata)
		if err != nil {
			return err
		}
	}

	return nil
}
