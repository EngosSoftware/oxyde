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
