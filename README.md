# Gosf

Es un servidor HTTP que permite servir carpetas. Este servidor utiliza tempaltes de HTML, CSS y scripts de JS para poder funcionar, asi que se necesita indicar esas rutas. Actualmente permite:

- Acceder a la carpeta servida
- Descargar archivos
- Borrar archivos
- Subir archivos

```bash
#Con -D se indica el directorio de templates
#Con -d el directorio que se quiere servir
#Con -p en el puerto que se quiere servir
#Y aunque de manera predeterminada viene desactivado con -r se activa el modo recursivo
gosf -D ./src/static/ -d ~/Archivos/JoacoFiles/Videos -p 80
```
