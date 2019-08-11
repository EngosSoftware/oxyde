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
 
package html

const IndexTemplate = `
{{range .Groups}}
  <div class="group-name">{{.Name}}</div>
  <div class="group-container">
    <table>
      <tbody>
        {{range .Endpoints}}
          <tr>
            <td>
              <div class="endpoint-summary-method http-method-{{.MethodLo}}">{{.MethodUp}}</div>
            </td>
            <td>
              <div class="endpoint-summary-uri"><a href="/endpoint-details?id={{.Id}}">{{.UrlPath}}</a></div>
            </td>
            <td>
              <div class="endpoint-summary-text">{{.Summary}}</div>
            </td>
          </tr>
        {{end}}
      </tbody>
    </table>
  </div>
{{end}}

<h1>MACIERZ UPRAWNIEŃ</h1>
<div>
  <table>
    <thead>
      <tr>
        <th></th>
        <th></th>
        {{range .RoleNames}}
          <th class="role-name">{{.}}</th>
        {{end}}
      </tr>
    </thead>
    <tbody>
      {{range .Groups}}
        <tr>
         <td colspan="2" class="group-name">{{.Name}}</td>
        </tr>  
        {{range .Endpoints}}
          <tr>
            <td>
              <div class="endpoint-summary-method http-method-{{.MethodLo}}">{{.MethodUp}}</div>
            </td>
            <td>
              <div class="endpoint-summary-uri"><a href="/endpoint-details?id={{.Id}}">{{.UrlPath}}</a></div>
            </td>
            {{range .Access}}
              <td style="text-align: center">
                <div class="access-{{.}}">{{.}}</div>
              </td>
            {{end}}
          </tr>
        {{end}}
      {{end}}
    </tbody>
  </table>
</div>
`