package helpers

import (
	"strings"
	"text/template"
)

const listTemplate = `
{{- range . }}
- {{ . }}
{{- end }}
`

var parsedListTemplate = template.Must(template.New("list").Parse(listTemplate))

func StringSliceToBulletedList(items []string) string {
	var sb strings.Builder

	if err := parsedListTemplate.Execute(&sb, items); err != nil {
		panic(err)
	}

	return sb.String()
}
