package generator

import "github.com/artsafin/coda-schema-generator/internal/api"

const (
	lookupTypeSuffix    = "Lookup"
	lookupRefTypeSuffix = "RowRef"
)

type Field struct {
	Name, Type, ConvertFn string
}

type LookupField struct {
	LookupValuesType string
	LookupRowRefType string
	DTOType          string
}

type FieldMapper struct {
	namer        nameConverter
	Fields       map[string][]Field
	lookupFields map[string]LookupField
}

func NewFieldMapper(namer nameConverter) *FieldMapper {
	return &FieldMapper{
		namer:        namer,
		Fields:       make(map[string][]Field),
		lookupFields: make(map[string]LookupField),
	}
}

func (m *FieldMapper) GetLookupFields() (r []LookupField) {
	for _, v := range m.lookupFields {
		r = append(r, v)
	}
	return
}

func (m *FieldMapper) registerField(tableID string, c api.Column) {
	typeData := mapColumnFormatToGoType(m.namer, c)

	m.Fields[tableID] = append(m.Fields[tableID], Field{
		Name:      m.namer.ConvertNameToGoSymbol(c.Name),
		Type:      typeData.LiteralType,
		ConvertFn: typeData.ValuerFn,
	})

	if c.Format.Type == api.ColumnFormatTypeLookup {
		m.registerLookup(c)
	}
}

func (m *FieldMapper) registerLookup(c api.Column) {
	if _, ok := m.lookupFields[c.Format.Table.ID]; ok {
		return
	}

	referencedTableSym := m.namer.ConvertNameToGoSymbol(c.Format.Table.Name)

	m.lookupFields[c.Format.Table.ID] = LookupField{
		LookupValuesType: referencedTableSym + lookupTypeSuffix,
		LookupRowRefType: referencedTableSym + lookupRefTypeSuffix,
		DTOType:          referencedTableSym,
	}
}

type ColumnTypeData struct {
	LiteralType string
	ValuerFn    string
}

func mapColumnFormatToGoType(namer nameConverter, c api.Column) ColumnTypeData {
	switch c.Format.Type {
	case api.ColumnFormatTypeDateTime:
		return ColumnTypeData{LiteralType: "time.Time", ValuerFn: "ToDateTime"}
	case api.ColumnFormatTypeTime:
		return ColumnTypeData{LiteralType: "time.Time", ValuerFn: "ToTime"}
	case api.ColumnFormatTypeDate:
		return ColumnTypeData{LiteralType: "time.Time", ValuerFn: "ToDate"}
	case api.ColumnFormatTypeScale:
		return ColumnTypeData{LiteralType: "uint8", ValuerFn: "ToUint8"}
	case api.ColumnFormatTypeNumber, api.ColumnFormatTypeSlider:
		return ColumnTypeData{LiteralType: "float64", ValuerFn: "ToFloat64"}
	case api.ColumnFormatTypeCheckbox:
		return ColumnTypeData{LiteralType: "bool", ValuerFn: "ToBool"}
	case api.ColumnFormatTypeLookup:
		t := namer.ConvertNameToGoSymbol(c.Format.Table.Name) + lookupTypeSuffix
		return ColumnTypeData{
			LiteralType: t,
			ValuerFn:    "To" + t,
		}
	case api.ColumnFormatTypePerson, api.ColumnFormatTypeReaction:
		return ColumnTypeData{LiteralType: "[]Person", ValuerFn: "ToPersons"}
	case api.ColumnFormatTypeImage, api.ColumnFormatTypeAttachments:
		return ColumnTypeData{LiteralType: "[]Attachment", ValuerFn: "ToAttachments"}
	case api.ColumnFormatTypeCurrency:
		return ColumnTypeData{LiteralType: "MonetaryAmount", ValuerFn: "ToMonetaryAmount"}
	}

	return ColumnTypeData{LiteralType: "string", ValuerFn: "ToString"}
}
