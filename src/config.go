package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const TextTemplate string = `

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

.error{
    color: rgb(223, 80, 80);
}
    </style>

</head>
    <body>
        <div class="files">
            {{range .Files}}
            <a id="{{.Index}}" href="{{.Link}}">{{.Name}}<a> <br>            
            {{end}}
        </div>

        <div class="error">
            <p>{{.Error}}</p>
        </div>

    </body>
</html>`

var (
	PathTempalteHtml string = "./template.html"
	NameTemplateHtml string = filepath.Base(PathTempalteHtml)
)

func CreateTemplate() error {
	if _, err := os.Stat(PathTempalteHtml); err == nil {
		log.Println("Usando " + NameTemplateHtml + " ya existente")
		return nil
	}

	return ioutil.WriteFile(NameTemplateHtml, []byte(TextTemplate), fs.ModeAppend)
}
