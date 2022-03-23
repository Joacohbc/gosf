//Eventos de botones

//Accion de descargar archivos
export function actionDescargar(e) {
    e.stopPropagation();
    e.preventDefault();

    const archivo = e.target.parentElement.parentElement.getAttribute("info");
    const url = `http://${window.location.host}/api/api/ownloadfiles/${archivo}`;
    const win = window.open(url, "_self");
}

//Accion de borrar un archivo
export function actionBorrar(e) {
    e.stopPropagation();
    e.preventDefault();

    if (!confirm("¿Quiere borrar el archivo seleccionado?")) {
        return;
    }

    /*
    Obtengo el elemento padre que es un <td>, y de ese <td> obtengo 
    el padre que es el que es el <tr>.
    
    De ese <tr> obtengo el atributo "info" que es donde esta el nombre
    del archivo
    */
    const archivo = e.target.parentElement.parentElement.getAttribute("info");

    //Creo la peticion
    const Http = new XMLHttpRequest();
    const url = `http://${window.location.host}/api/removefiles/${archivo}`;

    Http.open("DELETE", url);
    Http.onreadystatechange = () => {

        //Si no esta completada la transaccion(Estado Nro 4)
        /* 
            Si no remite la funcion 4 veces
            0: no inicializado. Indica que no se ha abierto la conexión con el servidor (no se ha llamado a open)

            1: conexión con servidor establecida. (no se ha llamado a open)

            2: recibida petición en servidor. (se ha llamado a send)

            3: enviando información. (se ha llamado a send)

            4: completado. Se ha recibido la información del servidor y esta listo para operar
        */
        if (Http.readyState != XMLHttpRequest.DONE) {
            return;
        }

        //Si el status es 202, es el StatusAccepted que envia el servidor
        if (Http.status == 202) {
            //Si es correcto  recargo para que ya no se muestre el archivo
            window.location.reload();
        } else {
            const respuesta = JSON.parse(Http.responseText);
            alert(respuesta.error);
        }
    };
    Http.send();
}

//Accion para obtener un archivo
export function actionObtener(e) {
    e.stopPropagation();
    e.preventDefault();

    //Leo la info del archivo
    const archivo = e.target.parentElement.parentElement.getAttribute("info");
    const url = `http://${window.location.host}/api/getfiles/${archivo}`;

    //Hago la peticion
    fetch(url, {
        method: 'GET',
    }).then(responseData => {
        
        //Si el status no es Ok, sigmifica que algo fallo
        //y se notificara con un JSON
        if (!responseData.ok) {

            //Sabiendo que sera un JSON, proceso el json
            responseData.json().then(json => {
                //Y notifico el error
                alert('Ocurrio un error al acceder al archivo: ' + json.error);
            }).catch((err) => {
                alert('Ocurrio un error al leer la respuesta del servidor');
                console.log(err);
            });
            return;
        }

        window.open(url, "_self");

    }).catch(error => {
        //Si ocurrio un error
        alert('Ocurrio un error:' + error);
        console.error(error);
    });
}

//Accion de cargar el HTML apartir
export function actionCargarDir(e){
    e.stopPropagation();
    e.preventDefault();

    //Vacio la tabla
    document.querySelector("tbody").innerHTML = "";

    //Leo la info del archivo
    const archivo = e.target.parentElement.parentElement.getAttribute("info");
    cargarHTML(archivo);
}

//Accion de subir un archivo
export function actionSubir(e) {

    e.stopPropagation();
    e.preventDefault();

    //Obtengo los archivos del selector de archivos
    const files = document.getElementById("selectorDeArchivos").files;

    //Comprobar si el selector de archivos tiene algun archivo
    //sino no aceptar
    if (files.length <= 0) {
        alert("Debe seleccinar algun archivo");
        return;
    }

    if (!confirm("¿Quiere enviar el archivo/s seleccionado?")) {
        return;
    }

    const data = new FormData();

    //Agrego todos los archivos en seleccionados
    //al FormData
    for (let i = 0; i < files.length; i++) {
        data.append('fileToUpload', files[i]);
    }

    //Envio los archivos a /uploadfiles con metodo POST
    fetch(`http://${window.location.host}/api/uploadfiles`, {
            method: 'POST',
            body: data
        })
        ///Convierto la respuesta JSON
        .then(response => response.json())
        .then(data => {
            //Si todo ocurrio exitosamente recargo
            //para ver los resultados
            window.location.reload();
        })
        .catch(error => {
            //Si ocurrio un error
            alert('Ocurrio un error:' + error.error);
            console.error(error);
        });
}

//Funcion de cargar HTML
export function cargarHTML(archivo) {
    const url = `http://${window.location.host}/api/getfiles/${archivo}`;

    //Hago la consulta por los archivo
    const respuesta = fetch(url, {
        method: 'GET',
    }).then(response => response.json()).then(data => {

        /*
            Si 'data' es null sigmifica que el json esta vacio, 
            si esta vacio sigmifica que no se subio ningun archivo
        */
        if(data == null){
            //Creo la fila (table row)
            const tr = document.createElement("tr");

            //Creo la columna de la no archivos
            const noFiles = document.createElement("td");
            noFiles.setAttribute("colspan","4");
            noFiles.innerHTML = "No hay archivos en el directorio servido";
            //Y agrego la columna de tiempo a la fila
            tr.appendChild(noFiles);

            document.querySelector(".files").appendChild(tr);
            return;
        }

        /*
            Creo la parte del documento donde agregare
            las filas
        */
        const part = document.createDocumentFragment();
        for (let i = 0; i < data.length; i++) {

            //Creo la fila (table row)
            const tr = document.createElement("tr");
            tr.setAttribute("info", data[i].link);

            //Creo la columna donde va el "Nombre"
            const name = document.createElement("th");

            //Agrego el <a> en el nombre
            const link = document.createElement("a");
            link.href = "/getfiles/"+data[i].link;
            link.innerHTML = data[i].name;
            if (data[i].isDir == true) {
                link.addEventListener("click", actionCargarDir);
            }else{
                link.addEventListener("click", actionObtener);
            }

            name.appendChild(link);
            //Y agrego la columna de nombres a la fila
            tr.appendChild(name);

            //Creo la columna de la "Modificación del tiempo"
            const modTime = document.createElement("td");
            modTime.innerHTML = data[i].sModTime;
            //Y agrego la columna de tiempo a la fila
            tr.appendChild(modTime);

            //Creo la columna del peso del archivo
            const size = document.createElement("td");
            size.innerHTML = formatBytes(data[i].size);
            //Y agrego la columna de peso a la fila
            tr.appendChild(size);
            
            //Creo la columna de botones
            const botones = document.createElement("td");

            //Creo el "boton de Borrar"
            const btnBorrar = document.createElement("button");
            btnBorrar.setAttribute("id", "del");
            btnBorrar.innerHTML = "Borrar";
            btnBorrar.addEventListener("click", actionBorrar);
            //Agrego el boton a la columna
            botones.appendChild(btnBorrar);

            //Creo el "boton de Descargar"
            const btnDescargar = document.createElement("button");
            btnDescargar.setAttribute("id", "download");
            btnDescargar.innerHTML = "Descargar";
            btnDescargar.addEventListener("click", actionDescargar);
            //Agrego los botones
            botones.appendChild(btnDescargar);

            //Y agrego la columna de botones a la fila
            tr.appendChild(botones);

            //Y agrego al final la fila a la parte del documetno
            part.appendChild(tr);
        }

        //Y agego el documento a la tabla(con la clase fila)
        document.querySelector("tbody").appendChild(part);
        
        //Le agrego el evento "click" al boton de subir archivos 
        document.getElementById("btnSubir").addEventListener("click", actionSubir);

    }).catch(error => {
        //Si ocurrio un error
        alert('Ocurrio un error:' + error);
        console.error(error);
    });
}

/* 
    Sacado de:
    https://stackoverflow.com/questions/15900485/correct-way-to-convert-size-in-bytes-to-kb-mb-gb-in-javascript
*/
export function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

/*
export function actionCargar(e) {
    e.stopPropagation();
    e.preventDefault();

    //Leo la info del archivo
    const archivo = e.target.parentElement.parentElement.getAttribute("info");
    const url = `http://${window.location.host}/api/getfiles/${archivo}`;

    //Hago la peticion
    fetch(url, {
        method: 'GET',
    }).then((response) => response).then((responseData) => {

        //Si el status no es Ok, sigmifica que algo fallo
        //y se notificara con un JSON
        if (!responseData.ok) {

            //Sabiendo que sera un JSON, proceso el json
            responseData.json().then(json => {
                //Y notifico el error
                alert('Ocurrio un error al acceder al archivo: ' + json.error);
            }).catch((err) => {
                alert('Ocurrio un error al leer la respuesta del servidor');
                console.log(err);
            });
            return;
        }
        
        //Si el status fue Ok, se que llegara un blob
        //entoces proceso un blob
        responseData.blob().then((blob) => {
			
            //Creo un url a ese blbo
            const link = window.URL.createObjectURL(blob);
            window.open(link, "_self");
            
        }).catch(error => {
            alert('Ocurrio un error al procesar el archivo:' + error);
            console.error(error);
        });

    }).catch(error => {
        //Si ocurrio un error
        alert('Ocurrio un error:' + error);
        console.error(error);
    });
}
*/