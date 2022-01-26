package api

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

type ColumnFormat struct {
	Type                  string      `json:"type"`
	IsArray               bool        `json:"isArray"`
	Table                 TableFormat `json:"table"`                 // lookup, person
	Precision             int         `json:"precision"`             // number, percent, currency, duration ([0 .. 10])
	UseThousandsSeparator bool        `json:"useThousandsSeparator"` // number, percent
	Format                string      `json:"format"`                // currency (enum "currency" "accounting" "financial"), date, time
	DateFormat            string      `json:"dateFormat"`            // dateTime
	TimeFormat            string      `json:"timeFormat"`            // dateTime
	CurrencyCode          string      `json:"currencyCode"`          // currency
	MaxUnit               string      `json:"maxUnit"`               // duration (enum "days" "hours" "minutes" "seconds")
	Minimum               string      `json:"minimum"`               // slider
	Maximum               string      `json:"maximum"`               // slider, scale
	Step                  string      `json:"step"`                  // slider
	Icon                  string      `json:"icon"`                  // scale (enum "star" "circle" "fire" "bug" "diamond" "bell" "thumbsup" "heart" "chili" "smiley" "lightning" "currency" "coffee" "person" "battery" "cocktail" "cloud" "sun" "checkmark" "lightbulb")
	Label                 string      `json:"label"`                 // button (formula)
	DisableIf             string      `json:"disableIf"`             // button (formula)
	Action                string      `json:"action"`                // button (formula)
}

type Column struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Formula    string       `json:"formula"`
	Calculated bool         `json:"calculated"`
	Format     ColumnFormat `json:"format"`
}

type TableColumns struct {
	TableID   string
	TableType string
	Items     []Column `json:"items"`
}

func (tc TableColumns) HasMutableRefColumns() bool {
	return len(tc.GetMutableRefColumns()) > 0
}

func (tc TableColumns) GetMutableRefColumns() (cs []Column) {
	if tc.TableType == "view" {
		return
	}

	for _, c := range tc.Items {
		if c.Format.Type == "lookup" && !c.Calculated {
			cs = append(cs, c)
		}
	}

	return cs
}
