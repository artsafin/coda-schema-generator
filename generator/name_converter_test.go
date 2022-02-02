package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_nameConverter_ConvertNameToGoSymbol(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{"Invoices", "Invoices"},
		{"Working employees names", "WorkingEmployeesNames"},
		{"1) employees", "No1Employees"},
		{"[999] employees", "No999Employees"},
		{"~999 all employees", "Tilde999AllEmployees"},
		{"test ;Non;letter characters", "TestNonLetterCharacters"},
		{"non-standard 🔥 name 0", "NonStandardName0"},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			d := newNameConverter()
			assert.Equal(t, tt.expected, d.ConvertNameToGoSymbol(tt.value))
		})
	}
}

func Test_nameConverter_ConvertNameToGoType(t *testing.T) {
	tests := []struct {
		name     string
		suffix   string
		expected string
	}{
		{"Working employees names", "Table", "_workingEmployeesNamesTable"},
		{"1 non-standard 🔥 name 2", "SUFFIX", "_no1NonStandardName2SUFFIX"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newNameConverter()
			assert.Equal(t, tt.expected, d.ConvertNameToGoType(tt.name, tt.suffix))
		})
	}
}
