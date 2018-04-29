package main

const htmlTemplate = `<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <title>{{ .Title }}</title>
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
    h1 {
      font-weight: normal;
    }
    main {
      padding-top: 100px;
      max-width: 940px;
      margin-left: auto;
      margin-right: auto;
    }
  </style>
</head>
<body>
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
	{{ if (len .Siblings) gt 1 }}
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
