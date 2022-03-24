# Gosf

Gosf es un servidor HTTP que permite servir archivos. Este servidor utiliza tempaltes de HTML, CSS y scripts de JS para poder funcionar, asi que se necesita indicar esas rutas (predeterminadamente usar $HOME/.config/gosf/static/). Actualmente permite:

- Acceder a la carpeta servida (y a sus subidrectorio)
- Descargar archivos
- Borrar archivos
- Subir archivos
- Pide una autentificacion basica para acceder a las funciones de enviar y borrar archivos, de manera predeterminada el usuario es "admin" y contraseña la "admin"

```bash
#Con -d el directorio que se quiere servir
#Con -p en el puerto que se quiere servir
#Con -U se indica el usuario y contraña
gosf -d ~/Documentos/ -p 8081 -U user,pass

Output:
2022/03/23 22:59:01 Flags:
2022/03/23 22:59:01 - Servidor abierto en: 8081
2022/03/23 22:59:01 - Ruta servida: /home/user/Documentos
2022/03/23 22:59:01 - Tiempo de cierre automatico: 0s
2022/03/23 22:59:01 - Se esta usando el directorio "/home/user/.config/gosf/static" para templates
2022/03/23 22:59:01 - Usuario de adminstrador: map[user:pass]
2022/03/23 22:59:01 Iniciando servidor...
```

Para apagar el servidor simplemente hay que hacer una señal de interrucion, un simple Ctrl+C(^C)

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
