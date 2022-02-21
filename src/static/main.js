document.addEventListener('DOMContentLoaded', function() {
    const btns = document.querySelectorAll("#del");
    btns.forEach(b => {
        b.addEventListener("click", (e) =>{

            if(confirm("Quiere borrar el archivo seleccionado?")){

                /*
                Obtengi el elemento padre que es un <td>, y de ese <td> obtengo 
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
                        //Digo que es correcto y luego recargo para
                        //que ya no se muestre el archivo
                        alert("El archivo se borro con exito");
                        window.location.reload();
                    }else { 
                        let respuesta = JSON.parse(Http.responseText);
                        alert(respuesta.error);
                    }
                };
                Http.send();
            }

            e.stopPropagation();
        });
    });
});