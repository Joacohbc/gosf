# Gosf

Gosf es un servidor HTTP que permite servir archivos. Este servidor utiliza tempaltes de HTML, CSS y scripts de JS para poder funcionar, asi que se necesita indicar esas rutas (predeterminadamente usar $HOME/.config/gosf/static/). Actualmente permite:

- Acceder a la carpeta servida
- Descargar archivos
- Borrar archivos
- Subir archivos
- Pide una autentificacion basica para acceder, de manera predeterminada el usuario es "admin" y contraseña la "admin"

```bash
#Con -D se indica el directorio de templates
#Con -d el directorio que se quiere servir
#Con -p en el puerto que se quiere servir
gosf -d ~/Documentos/Videos -p 80
```

## Instalacion (only Linux)

La instalacion consta de darle permiso de ejecucion al binario y moverlo al el directorio de binarios, en Linux es /usr/bin, y ademas mover la carpeta que contiene el HTML, CSS y JS (static) a la carpeta de configuracion ($HOME/.config/gosf/). Esto ya lo hace automaticamente el install.sh:

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
