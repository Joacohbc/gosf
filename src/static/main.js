document.addEventListener('DOMContentLoaded', function() {
    //Selecciono todos los botones de borrar
    document.querySelectorAll("#del").forEach(b => {
        //Y agrego a todos los botones esos eventos
        b.addEventListener("click", (e) =>{
            e.stopPropagation();
            e.preventDefault();
            if(!confirm("¿Quiere borrar el archivo seleccionado?")){
                return;
            }

            /*
            Obtengo el elemento padre que es un <td>, y de ese <td> obtengo 
            el padre que es el que es el <tr>.
            
            De ese <tr> obtengo el atributo "info" que es donde esta el nombre
            del archivo
            */
            let archivo = b.parentElement.parentElement.getAttribute("info");

            //Creo la peticion
            const Http = new XMLHttpRequest();
            const url=`http://${window.location.host}/removefiles/${archivo}`;
            
            Http.open("POST", url);
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
                if(Http.readyState != XMLHttpRequest.DONE){
                    return;
                }

                //Si el status es 202, es el StatusAccepted que envia el servidor
                if (Http.status == 202 ){
                    //Si es correcto  recargo para que ya no se muestre el archivo
                    window.location.reload();
                }else { 
                    let respuesta = JSON.parse(Http.responseText);
                    alert(respuesta.error);
                }
            };
            Http.send();
        });
    });

    //Le agrego el evento "click" al boton de subir archivos 
    document.getElementById("btnSubir").addEventListener("click", function(e){
        e.stopPropagation();
        e.preventDefault();

        //Obtengo los archivos del selector de archivos
        const files = document.getElementById("selectorDeArchivos").files;

        if(files.length <= 0){
            alert("Debe seleccinar algun archivo");
            return;
        }

        if(!confirm("¿Quiere enviar el archivo/s seleccionado?")){
            return;
        }

        const data = new FormData();

        //Agrego todos los archivos en seleccionados
        //al FormData
        for(let i=0; i<files.length; i++) {
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
            //Si todo ocurrio exitosamente
            window.location.reload();
        })
        .catch(error => {
            //Si ocurrio un error
            alert('Ocurrio un error:' + error.error);
            console.error(error);
        }).finally(() => {
            
        }); 
    });

});
