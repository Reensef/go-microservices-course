package notification

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

var (
	orderPaidTemplate     = template.Must(template.ParseFS(templatesFS, "templates/order_paid.tmpl"))
	shipAssembledTemplate = template.Must(template.ParseFS(templatesFS, "templates/ship_assembled.tmpl"))
)

// TODO разбить по функциям renderOrderPaidTemplate, renderShipAssembledTemplate,
// чтобы не допустить ошибку при передаче данных в шаблон, т.к. разные шаблоны используют разные структуры данных
func renderTemplate(tmpl *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %q: %w", tmpl.Name(), err)
	}

	return buf.String(), nil
}
