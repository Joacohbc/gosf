export function actionDescargar(e) {
    e.stopPropagation();
    e.preventDefault();

    let archivo = e.target.parentElement.parentElement.getAttribute("info");
    let url = `http://${window.location.host}/downloadfiles/${archivo}`;
    let win = window.open(url, "_blank");
    win.focus();
}

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
    let archivo = e.target.parentElement.parentElement.getAttribute("info");

    //Creo la peticion
    const Http = new XMLHttpRequest();
    const url = `http://${window.location.host}/removefiles/${archivo}`;

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
            let respuesta = JSON.parse(Http.responseText);
            alert(respuesta.error);
        }
    };
    Http.send();
}

export function actionObtener(e) {
    e.stopPropagation();
    e.preventDefault();

    //Leo la info del archivo
    const archivo = e.target.parentElement.parentElement.getAttribute("info");
    const url = `http://${window.location.host}/getfiles/${archivo}`;

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

        window.open(url, "_self");

    }).catch(error => {
        //Si ocurrio un error
        alert('Ocurrio un error:' + error);
        console.error(error);
    });
}

export function actionCargar(e) {
    e.stopPropagation();
    e.preventDefault();

    //Leo la info del archivo
    const archivo = e.target.parentElement.parentElement.getAttribute("info");
    const url = `http://${window.location.host}/getfiles/${archivo}`;

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
            let link = window.URL.createObjectURL(blob);
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
    fetch(`http://${window.location.host}/uploadfiles`, {
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