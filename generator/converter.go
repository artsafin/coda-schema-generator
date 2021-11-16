package generator

import (
	"regexp"
	"strings"
)

type nameConverter struct {
	symbolizer         *strings.Replacer
	nonLetterOrDigitRe *regexp.Regexp
	pascalCaseRe       *regexp.Regexp
}

// Taken from https://gist.github.com/elliotchance/d419395aa776d632d897
func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			if v[i] == -1 || v[i+1] == -1 {
				groups = append(groups, "")
			} else {
				groups = append(groups, str[v[i]:v[i+1]])
			}
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

func NewNameConverter() nameConverter {
	return nameConverter{
		symbolizer: strings.NewReplacer(
			"#", "Hash",
			"$", "Dollar",
			"%", "Percent",
			"&", "And",
			"*", "Star",
			"+", "Plus",
			"<", "Less",
			"=", "Equal",
			">", "Greater",
			"@", "At",
			"^", "Accent",
			"~", "Tilde",
		),
		nonLetterOrDigitRe: regexp.MustCompile("[^0-9a-zA-Z]"),
		pascalCaseRe:       regexp.MustCompile("(?m)(^|[^0-9a-zA-Z])([0-9a-zA-Z])"),
	}
}

func (d *nameConverter) ConvertNameToGoSymbol(name string) string {
	name = ReplaceAllStringSubmatchFunc(d.pascalCaseRe, name, func(g []string) string {
		return g[1] + strings.ToUpper(g[2])
	})

	name = d.symbolizer.Replace(name)
	name = d.nonLetterOrDigitRe.ReplaceAllLiteralString(name, "")

	if name[0] >= 0x30 && name[0] <= 0x39 { // First byte is number
		name = "No" + name
	}

	return name
}

func (d *nameConverter) ConvertNameToGoType(name string, suffix string) string {
	sym := []rune(d.ConvertNameToGoSymbol(name))

	return "_" + strings.ToLower(string(sym[0])) + string(sym[1:]) + suffix
}
