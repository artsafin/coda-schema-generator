package generator

import (
	"github.com/artsafin/coda-schema-generator/dto"
)

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

func (m *FieldMapper) registerField(tableID string, c dto.Column) {
	typeData := mapColumnFormatToGoType(m.namer, c)

	m.Fields[tableID] = append(m.Fields[tableID], Field{
		Name:      m.namer.ConvertNameToGoSymbol(c.Name),
		Type:      typeData.LiteralType,
		ConvertFn: typeData.ValuerFn,
	})

	if c.Format.Type == dto.ColumnFormatTypeLookup {
		m.registerLookup(c)
	}
}

func (m *FieldMapper) registerLookup(c dto.Column) {
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

func mapColumnFormatToGoType(namer nameConverter, c dto.Column) ColumnTypeData {
	switch c.Format.Type {
	case dto.ColumnFormatTypeDateTime:
		return ColumnTypeData{LiteralType: "time.Time", ValuerFn: "ToDateTime"}
	case dto.ColumnFormatTypeTime:
		return ColumnTypeData{LiteralType: "time.Time", ValuerFn: "ToTime"}
	case dto.ColumnFormatTypeDate:
		return ColumnTypeData{LiteralType: "time.Time", ValuerFn: "ToDate"}
	case dto.ColumnFormatTypeScale:
		return ColumnTypeData{LiteralType: "uint8", ValuerFn: "ToUint8"}
	case dto.ColumnFormatTypeNumber, dto.ColumnFormatTypeSlider:
		return ColumnTypeData{LiteralType: "float64", ValuerFn: "ToFloat64"}
	case dto.ColumnFormatTypeCheckbox:
		return ColumnTypeData{LiteralType: "bool", ValuerFn: "ToBool"}
	case dto.ColumnFormatTypeLookup:
		t := namer.ConvertNameToGoSymbol(c.Format.Table.Name) + lookupTypeSuffix
		return ColumnTypeData{
			LiteralType: t,
			ValuerFn:    "To" + t,
		}
	case dto.ColumnFormatTypePerson, dto.ColumnFormatTypeReaction:
		return ColumnTypeData{LiteralType: "[]Person", ValuerFn: "ToPersons"}
	case dto.ColumnFormatTypeImage, dto.ColumnFormatTypeAttachments:
		return ColumnTypeData{LiteralType: "[]Attachment", ValuerFn: "ToAttachments"}
	case dto.ColumnFormatTypeCurrency:
		return ColumnTypeData{LiteralType: "MonetaryAmount", ValuerFn: "ToMonetaryAmount"}
	}

	return ColumnTypeData{LiteralType: "string", ValuerFn: "ToString"}
}
