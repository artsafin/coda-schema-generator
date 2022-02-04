package generator

import (
	"github.com/artsafin/coda-schema-generator/dto"
)

const (
	lookupTypeSuffix    = "Lookup"
	lookupRefTypeSuffix = "RowRef"
)

type field struct {
	Name, Type, ConvertFn string
}

type lookupField struct {
	LookupValuesType string
	LookupRowRefType string
	DTOType          string
	TableName        string
}

type fieldMapper struct {
	namer        nameConverter
	Fields       map[string][]field
	lookupFields map[string]lookupField
}

func newFieldMapper(namer nameConverter) *fieldMapper {
	return &fieldMapper{
		namer:        namer,
		Fields:       make(map[string][]field),
		lookupFields: make(map[string]lookupField),
	}
}

func (m *fieldMapper) GetLookupFields() (r []lookupField) {
	for _, v := range m.lookupFields {
		r = append(r, v)
	}
	return
}

func (m *fieldMapper) registerField(tableID string, c dto.Column) {
	typeData := mapColumnFormatToGoType(m.namer, c)

	m.Fields[tableID] = append(m.Fields[tableID], field{
		Name:      m.namer.ConvertNameToGoSymbol(c.Name),
		Type:      typeData.LiteralType,
		ConvertFn: typeData.ValuerFn,
	})

	if c.Format.Type == dto.ColumnFormatTypeLookup {
		m.registerLookup(c)
	}
}

func (m *fieldMapper) registerLookup(c dto.Column) {
	if _, ok := m.lookupFields[c.Format.Table.ID]; ok {
		return
	}

	referencedTableSym := m.namer.ConvertNameToGoSymbol(c.Format.Table.Name)

	m.lookupFields[c.Format.Table.ID] = lookupField{
		TableName:        c.Format.Table.Name,
		LookupValuesType: referencedTableSym + lookupTypeSuffix,
		LookupRowRefType: referencedTableSym + lookupRefTypeSuffix,
		DTOType:          referencedTableSym,
	}
}

type columnTypeData struct {
	LiteralType string
	ValuerFn    string
}

func mapColumnFormatToGoType(namer nameConverter, c dto.Column) columnTypeData {
	switch c.Format.Type {
	case dto.ColumnFormatTypeDateTime:
		return columnTypeData{LiteralType: "time.Time", ValuerFn: "ToDateTime"}
	case dto.ColumnFormatTypeTime:
		return columnTypeData{LiteralType: "time.Time", ValuerFn: "ToTime"}
	case dto.ColumnFormatTypeDate:
		return columnTypeData{LiteralType: "time.Time", ValuerFn: "ToDate"}
	case dto.ColumnFormatTypeScale:
		return columnTypeData{LiteralType: "uint8", ValuerFn: "ToUint8"}
	case dto.ColumnFormatTypeNumber, dto.ColumnFormatTypeSlider:
		return columnTypeData{LiteralType: "float64", ValuerFn: "ToFloat64"}
	case dto.ColumnFormatTypeCheckbox:
		return columnTypeData{LiteralType: "bool", ValuerFn: "ToBool"}
	case dto.ColumnFormatTypeLookup:
		t := namer.ConvertNameToGoSymbol(c.Format.Table.Name) + lookupTypeSuffix
		return columnTypeData{
			LiteralType: t,
			ValuerFn:    "To" + t,
		}
	case dto.ColumnFormatTypePerson, dto.ColumnFormatTypeReaction:
		return columnTypeData{LiteralType: "[]Person", ValuerFn: "ToPersons"}
	case dto.ColumnFormatTypeImage, dto.ColumnFormatTypeAttachments:
		return columnTypeData{LiteralType: "[]Attachment", ValuerFn: "ToAttachments"}
	case dto.ColumnFormatTypeCurrency:
		return columnTypeData{LiteralType: "MonetaryAmount", ValuerFn: "ToMonetaryAmount"}
	}

	return columnTypeData{LiteralType: "string", ValuerFn: "ToString"}
}
