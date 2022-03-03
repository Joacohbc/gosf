import * as events from "./btnEvents.js";

export default function cargarDOM() {
    
    //Hago la consulta por los archivo
    const respuesta = fetch(`http://${window.location.host}/getfiles`, {
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
            link.addEventListener("click", events.actionObtener);

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
            size.innerHTML = events.formatBytes(data[i].size);
            //Y agrego la columna de peso a la fila
            tr.appendChild(size);
            
            //Creo la columna de botones
            const botones = document.createElement("td");

            //Creo el "boton de Borrar"
            const btnBorrar = document.createElement("button");
            btnBorrar.setAttribute("id", "del");
            btnBorrar.innerHTML = "Borrar";
            btnBorrar.addEventListener("click", events.actionBorrar);
            //Agrego el boton a la columna
            botones.appendChild(btnBorrar);

            //Creo el "boton de Descargar"
            const btnDescargar = document.createElement("button");
            btnDescargar.setAttribute("id", "download");
            btnDescargar.innerHTML = "Descargar";
            btnDescargar.addEventListener("click", events.actionDescargar);
            //Agrego los botones
            botones.appendChild(btnDescargar);
            
            //Creo el "boton de Cargar"
            const btnCargar = document.createElement("button");
            btnCargar.setAttribute("id", "open");
            btnCargar.innerHTML = "Cargar";
            btnCargar.addEventListener("click", events.actionCargar);
            
            //Agrego los botones
            botones.appendChild(btnCargar);

            //Y agrego la columna de botones a la fila
            tr.appendChild(botones);

            //Y agrego al final la fila a la parte del documetno
            part.appendChild(tr);
        }

        //Y agego el documento a la tabla(con la clase fila)
        document.querySelector("tbody").appendChild(part);
        
        //Le agrego el evento "click" al boton de subir archivos 
        document.getElementById("btnSubir").addEventListener("click", events.actionSubir);

    }).catch(error => {
        //Si ocurrio un error
        alert('Ocurrio un error:' + error);
        console.error(error);
    });
}

