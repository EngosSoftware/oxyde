package oxyde

import (
	"bytes"
	"fmt"
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

func StartPreview(dc *DocContext) {
	model := CreatePreviewModel(dc)
	runPreviewServer(model, 16100)
}

func runPreviewServer(model *PreviewModel, port int) {

	index := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		wrapInPage(w, indexTemplate, model)
	}

	styleCss := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		_, err := io.WriteString(w, StyleCss)
		panicOnError(err)
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/style.css", styleCss)
	mux.HandleFunc("/endpoint-details", endpointDetails)
	fmt.Printf("Documentation server started and listening on port: %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func preparePageTemplate() *template.Template {
	t, err := template.New("pageTemplate").Parse(PageTemplate)
	panicOnError(err)
	return t
}

func prepareIndexTemplate() *template.Template {
	t, err := template.New("indexTemplate").Parse(IndexTemplate)
	panicOnError(err)
	return t
}

func prepareEndpointTemplate() *template.Template {
	t, err := template.New("endpointTemplate").Parse(EndpointTemplate)
	panicOnError(err)
	return t
}

func prepareErrorTemplate() *template.Template {
	t, err := template.New("errorTemplate").Parse(ErrorTemplate)
	panicOnError(err)
	return t
}

func wrapInPage(w http.ResponseWriter, t *template.Template, data interface{}) {
	var out bytes.Buffer
	outWriter := io.Writer(&out)
	err := t.Execute(outWriter, data)
	panicOnError(err)
	err = pageTemplate.Execute(w, out.String())
	panicOnError(err)
}
