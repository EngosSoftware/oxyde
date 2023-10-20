package server

import (
	"bytes"
	"fmt"
	"github.com/wisbery/oxyde/common"
	d "github.com/wisbery/oxyde/doc"
	h "github.com/wisbery/oxyde/html"
	m "github.com/wisbery/oxyde/model"
	"io"
	"log"
	"net/http"
	"text/template"
)

var (
	pageTemplate     = preparePageTemplate()
	indexTemplate    = prepareIndexTemplate()
	endpointTemplate = prepareEndpointTemplate()
	errorTemplate    = prepareErrorTemplate()
)

func StartPreview(dc *d.Context) {
	model := m.CreateModel(dc)
	runPreviewServer(model, 16100)
}

func runPreviewServer(model *m.Model, port int) {

	index := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrapInPage(w, indexTemplate, model)
	}

	styleCss := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		_, err := io.WriteString(w, h.StyleCss)
		common.PanicOnError(err)
	}

	endpointDetails := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		ids, ok := req.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			wrapInPage(w, errorTemplate, "no endpoint identifier")
			return
		}
		if endpoint := model.FindEndpointById(ids[0]); endpoint != nil {
			wrapInPage(w, endpointTemplate, endpoint)
			return
		}
		wrapInPage(w, indexTemplate, model)
		return
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/style.css", styleCss)
	http.HandleFunc("/endpoint-details", endpointDetails)
	fmt.Printf(">> API preview server started and listening on port: %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func preparePageTemplate() *template.Template {
	t, err := template.New("pageTemplate").Parse(h.PageTemplate)
	common.PanicOnError(err)
	return t
}

func prepareIndexTemplate() *template.Template {
	t, err := template.New("indexTemplate").Parse(h.IndexTemplate)
	common.PanicOnError(err)
	return t
}

func prepareEndpointTemplate() *template.Template {
	t, err := template.New("endpointTemplate").Parse(h.EndpointTemplate)
	common.PanicOnError(err)
	return t
}

func prepareErrorTemplate() *template.Template {
	t, err := template.New("errorTemplate").Parse(h.ErrorTemplate)
	common.PanicOnError(err)
	return t
}

func wrapInPage(w http.ResponseWriter, t *template.Template, data interface{}) {
	var out bytes.Buffer
	outWriter := io.Writer(&out)
	err := t.Execute(outWriter, data)
	common.PanicOnError(err)
	err = pageTemplate.Execute(w, out.String())
	common.PanicOnError(err)
}
