// Package appcharts builds ECharts-based chart components from query read models.
package charts

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
)

// Renderable is any chart that can render its HTML representation to a
// writer. It is satisfied by the chart types of the go-echarts library.
type Renderable interface {
	Render(w io.Writer) error
}

// ConvertChartToTemplComponent cuts unnecessary html code from echarts
// (only get what's inside <body></body>) and makes it templ component.
func ConvertChartToTemplComponent(chart Renderable) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		buf := bytes.NewBuffer(nil)
		err := chart.Render(buf)
		if err != nil {
			return err
		}
		html := buf.String()
		bodyStart := strings.Index(html, "<body>")
		bodyEnd := strings.Index(html, "</body>")
		fragment := html[bodyStart+len("<body>") : bodyEnd]

		idStart := strings.Index(fragment, "id=\"")
		id := fragment[idStart+len("id=\"") : idStart+len("id=\"")+12]

		if _, err := io.WriteString(w, fragment); err != nil {
			return err
		}

		_, err = fmt.Fprintf(
			w, `<script>
		// block resize when chart is animating
		let initialized_%s = false;
		const observer_%s = new ResizeObserver(() => {
			if (!initialized_%s) {
				initialized_%s = true;
				return;
			}
			goecharts_%s.resize()
		})
		observer_%s.observe(document.getElementById("%s"))
		</script>`, id, id, id, id, id, id, id,
		)
		return err
	})
}
