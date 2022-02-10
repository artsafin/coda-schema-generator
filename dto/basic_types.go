package dto

type Schema struct {
	Tables   TableList
	Columns  map[string]TableColumns
	Formulas EntityList
	Controls EntityList
}

type ItemsContainer interface {
	Count() int
}

type Entity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Table struct {
	Entity
	TableType string `json:"tableType"`
}

type EntityList struct {
	Items []Entity `json:"items"`
}

func (e EntityList) Count() int {
	return len(e.Items)
}

type TableList struct {
	Items []Table `json:"items"`
}

func (t TableList) Count() int {
	return len(t.Items)
}

type TableFormat struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	TableType   string `json:"tableType"`
	Name        string `json:"name"`
	Href        string `json:"href"`
	BrowserLink string `json:"browserLink"`
}
