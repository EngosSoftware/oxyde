package html

const PageTemplate = `
<!DOCTYPE html>
<html lang="pl">
<head>
    <base href="/">
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="keywords" content="oxyde.io">
    <title>API Preview</title>
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Lato:400,700&display=swap&subset=latin-ext">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto+Mono:400,500,700&display=swap">
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
    <link rel="stylesheet" href="style.css" type="text/css">
</head>
<body>
  {{.}}
</body>
</html>
`
