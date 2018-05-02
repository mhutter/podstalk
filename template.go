package podstalk

const htmlTemplate = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
{{- if gt .Refresh 0 }}
  <meta http-equiv="refresh" content="{{ .Refresh }}">
{{- end }}
  <title>{{ .Title }}</title>
  <link rel="shortcut icon" type="image/png" href="https://kubernetes.io/images/favicon.png">
  <link href="https://fonts.googleapis.com/css?family=Ubuntu:400,700" rel="stylesheet">
  <style type="text/css">
    html, body {
      margin: 0;
      padding: 0;
      font-family: Ubuntu,"Helvetica Neue",Arial,Helvetica,sans-serif;
      font-size: 20px;
      color: #ffffff;
      background: #002F47;
    }
    h1,h2,h3 {
      font-weight: normal;
    }
    header,main {
      padding: 8px;
      max-width: 940px;
      margin-left: auto;
      margin-right: auto;
    }
    header {
      text-align: center;
    }
    header img {
      width: 100%;
      max-width: 432px;
    }
  </style>
</head>
<body>
  <header>
    <img src="{{ .BasePath }}/appuioli.png"/>
  </header>
  <main>
    <h1>Hi, I'm <strong>{{ .Name }}</strong>!</h1>
    <p>
      I live in the
      <strong>{{ .Namespace }}</strong>
      namespace. If you want to call me, my IP is
      <strong>{{ .IP }}</strong>
    </p>
    <p>
      I run on node
      <strong>{{ .NodeName }}</strong>
      (which listens to IP
      <strong>{{ .NodeIP }}</strong>
      )
    </p>
    <p>
      The current time is <strong>{{ .Now }}</strong>, plus/minus a bit.
    </p>
{{ if .Siblings }}
    <h3>Siblings</h3>
    {{ if gt (len .Siblings) 1 }}
    <p>
      I have a few siblings, together we're:
      <ul>
        {{ range .Siblings }}<li>{{ . }}</li>{{ end }}
      </ul>
    </p>
    {{ else }}
    <p>I don't have any siblings, I'm an only child :-(<p>
    {{ end }}
{{ end }}
  </main>
</body>
</html>`
