package dto

type TableColumns struct {
	TableID   string
	TableType string
	Items     []Column `json:"items"`
}

func (tc TableColumns) HasReferencedTables() bool {
	if tc.TableType == TableTypeView {
		return false
	}

	for _, c := range tc.Items {
		if c.Format.Type == ColumnFormatTypeLookup && !c.Calculated {
			return true
		}
	}
	return false
}

func (tc TableColumns) GetReferencedTables() []TableFormat {
	if tc.TableType == TableTypeView {
		return nil
	}

	fm := map[string]TableFormat{}

	for _, c := range tc.Items {
		if c.Format.Type == ColumnFormatTypeLookup && !c.Calculated {
			fm[c.Format.Table.Name] = c.Format.Table
		}
	}

	formats := make([]TableFormat, 0, len(fm))
	for _, v := range fm {
		formats = append(formats, v)
	}

	return formats
}

func (tc TableColumns) GetColumnsReferencedTo(refTable string) (cs []Column) {
	for _, c := range tc.Items {
		if c.Format.Type == ColumnFormatTypeLookup && c.Format.Table.Name == refTable {
			cs = append(cs, c)
		}
	}
	return
}
