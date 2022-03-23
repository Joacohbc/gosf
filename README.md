# Gosf

Es un servidor HTTP que permite servir carpetas. Este servidor utiliza tempaltes de HTML, CSS y scripts de JS para poder funcionar, asi que se necesita indicar esas rutas (predeterminadamente usar $HOME/.config/gosf/static/). Actualmente permite:

- Acceder a la carpeta servida
- Descargar archivos
- Borrar archivos
- Subir archivos

```bash
#Con -D se indica el directorio de templates
#Con -d el directorio que se quiere servir
#Con -p en el puerto que se quiere servir
gosf -D ./src/static/ -d ~/Archivos/JoacoFiles/Videos -p 80
```

## Instalar

La instalacion consta de mover el binario a el directorio de binarios, generalmente /usr/bin, y la carpeta de templates
moverla a la carpeta de configuracion ($HOME/.config/gosf/). Esto ya lo hace el install.sh:

- Clonar el repositorio

```bash
    git clone https://github.com/Joacohbc/gosf.git
```

- Ingresar a la carpeta

```bash
    cd gosf
```

- Ejecutar el intalador

```bash
    sh ./install.sh
```
