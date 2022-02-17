package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

const tpl string = `
<!DOCTYPE html>
<html>
<head>
    <meta charset='utf-8'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge'>
    <title>Archivos</title>
    <meta name='viewport' content='width=device-width, initial-scale=1'>

    <style>
body{
    background: rgba(0,0,0,0.9);
    color: #fff
}

a {
    text-align: center;
    color: #fff;
    font-size: 1.5em;
}

.files{
    margin: auto;
    display: block;
}
    </style>
</head>
    <body>
        <div class="files">
            {{range .Files}}
            <a id="{{.Index}}" href="{{.Link}}">{{.Name}}<a> <br>            
            {{end}}
        </div>
    </body>
</html>
`

const PathHtml string = "./template.html"

func CreateTemplate() error {

	if _, err := os.Stat(PathHtml); err == nil {
		log.Println("Usando template.html ya existente")
		return nil
	}
	return ioutil.WriteFile("template.html", []byte(tpl), fs.ModeAppend)
}
