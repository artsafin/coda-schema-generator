package dto

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
	Minimum               int         `json:"minimum"`               // slider
	Maximum               int         `json:"maximum"`               // slider, scale
	Step                  float64     `json:"step"`                  // slider
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

func (c Column) IsMutableLookup() bool {
	return c.IsLookup() && !c.Calculated
}

func (c Column) IsLookup() bool {
	return c.Format.Type == ColumnFormatTypeLookup && c.Format.Table.ID != SpecialTableGlobalExternalConnections
}
