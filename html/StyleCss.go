package html

const StyleCss = `
body {
  margin: 4px;
  color: black;
  background-color: #FFFFFF;
  font-family: 'Lato', Helvetica, sans-serif;
  font-size: 1em;
}

a {
  text-decoration: none;
  cursor: pointer;
  color: black;
}

table {
  border-collapse: collapse;
  border: none;
  border-spacing: 0;
}

h2 {
  margin-top: 4px;
  margin-bottom: 4px;
  background-color: grey;
}

h3 {
  margin-top: 4px;
  margin-bottom: 4px;
}

h4 {
  margin-top: 4px;
  margin-bottom: 4px;
}

hr {
  margin: 10px 0 10px 0;
  height: 1px;
  background-color: gray;
  border: none;
  width: 70%;
}

pre {
  font-family: 'Roboto Mono', monospace;
  font-size: 1em;
  color: #eeeeee;
  background-color: #37474f;
  padding: 10px;
  margin: 0;
  border-radius: 10px;
  border: none;
  white-space: pre-wrap;
  max-width: 1200px;
  overflow-x: auto;
}

.endpoint-summary-method {
  font-weight: bold;
}

.endpoint-summary-uri {
  font-weight: bold;
}

.endpoint-summary-text {
  padding-left: 20px;
}

.group-name {
  font-size: 1.05em;
  font-weight: bold;
  color: black;
  padding: 10px 0 10px 0;
}

.role-name {
  padding: 0 10px 0 10px;
}

.group-container td {
  padding: 0 6px 4px 6px; 
}

.endpoint-details-method {
  display: inline-block;
  font-weight: bold;
  border-radius: 6px;
  text-align: center;
  vertical-align: middle;
  width: 90px;
  min-width: 90px;
  height: 28px;
  line-height: 25px;
}

.endpoint-details-uri {
  display: inline-block;
  font-weight: bold;
  color: black;
  padding: 4px 10px 0 10px;
}

.endpoint-details-summary {
  font-weight: bold;
  font-size: 1.5em;
  margin: 10px 0 10px 0;
}

.fields-container-title {
  font-size: 1.3em;
  font-weight: bold;
  margin: 10px 0 8px 0;
}

.http-method-get {
  color: blue;
}

.details-http-method-get {
  color: white;
  background-color: blue;
}

.http-method-post {
  color: green;
}

.details-http-method-post {
  color: white;
  background-color: green;
}

.http-method-put {
  color: orange;
}

.details-http-method-put {
  color: white;
  background-color: orange;
}

.http-method-delete {
  color: red;
}

.details-http-method-delete {
  color: white;
  background-color: red;
}

.http-status {
  font-weight: bold;
  border-radius: 6px;
  text-align: center;
  vertical-align: middle;
  width: 90px;
  min-width: 90px;
  height: 28px;
  line-height: 25px;
}

.http-status-200 {
  color: white;
  background-color: blue;
}

.http-status-400 {
  color: white;
  background-color: red;
}

.http-status-401 {
  color: white;
  background-color: red;
}

.http-status-404 {
  color: white;
  background-color: red;
}

.json-type-object {
  font-weight: bold;
}

.json-type-array {
  font-weight: bold;
}

.json-mandatory-yes {
  text-align: center;
  font-weight: bold;
  color: black;
}

.json-mandatory-no {
  text-align: center;
  font-weight: normal;
  color: gray;
}

.parameters-description table {
  border: solid 1px grey;
}

.parameters-description th {
  border: solid 1px grey;
  padding: 2px 6px 2px 6px;  
}

.parameters-description td {
  border: solid 1px grey;
  padding: 2px 6px 2px 6px;
}

.example {
  display: flex; 
  flex-direction: column; 
  justify-content: flex-start; 
  align-items: flex-start;
  margin-bottom: 20px;
}

.example-container {
  display: flex; 
  flex-direction: column; 
  justify-content: flex-start; 
  align-items: flex-start;
  border: solid 1px gray;
  border-radius: 10px;
  padding: 10px;
}

.example-summary {
  font-size: 1.2em;
  font-weight: bold;
  color: #004d40;
  margin-bottom: 8px;
  max-width: 1000px;
}

.example-description {
  font-size: 0.9em;
  color: black;
  margin-bottom: 8px;
  max-width: 1000px;
}

.example-request {
  margin: 4px 0 8px 0;
}

.example-request-body {
  margin-left: 94px;
}

.example-response {
  margin: 4px 0 8px 0;
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  align-content: flex-start;
}

.example-response-body {
  margin-left: 4px;
}

.access-YES {
  color: white;
  background-color: green;
  font-weight: bold;
  padding: 3px 4px 4px 4px;
  border-radius: 10px;
}

.access-NO {
  color: red;
  background-color: white;
  padding: 3px 4px 4px 4px;
}

`
