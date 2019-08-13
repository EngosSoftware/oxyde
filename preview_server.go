/*
 * MIT License
 *
 * Copyright (c) 2017-2019 Dariusz Depta Engos Software
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE. 
 */
 
package oxyde

import (
    "bytes"
    "fmt"
    h "github.com/EngosSoftware/oxyde/html"
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
        _, err := io.WriteString(w, h.StyleCss)
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

    http.HandleFunc("/", index)
    http.HandleFunc("/style.css", styleCss)
    http.HandleFunc("/endpoint-details", endpointDetails)
    fmt.Printf(">> API preview server started and listening on port: %d\n", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func preparePageTemplate() *template.Template {
    t, err := template.New("pageTemplate").Parse(h.PageTemplate)
    panicOnError(err)
    return t
}

func prepareIndexTemplate() *template.Template {
    t, err := template.New("indexTemplate").Parse(h.IndexTemplate)
    panicOnError(err)
    return t
}

func prepareEndpointTemplate() *template.Template {
    t, err := template.New("endpointTemplate").Parse(h.EndpointTemplate)
    panicOnError(err)
    return t
}

func prepareErrorTemplate() *template.Template {
    t, err := template.New("errorTemplate").Parse(h.ErrorTemplate)
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
