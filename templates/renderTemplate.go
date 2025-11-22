package templates

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

func ConvertComponentToHtml(component templ.Component) ([]byte, error) {
	var b bytes.Buffer
	err := component.Render(context.Background(), &b)
	if err != nil {
		return nil, err
	}
	htmlBytes := b.Bytes()

	return htmlBytes, nil 
}