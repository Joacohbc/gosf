# Gosf

Gosf es un servidor HTTP que permite servir archivos. Este servidor utiliza templates de HTML, CSS y scripts de JS para poder funcionar, así que se necesita indicar esas rutas (predeterminadamente usar $HOME/.config/gosf/static/). Actualmente permite:

- Acceder a la carpeta servida (y a sus subdirectorio)
- Descargar archivos
- Borrar archivos
- Subir archivos
- Pide una autentificación básica para acceder a las funciones de enviar y borrar archivos, de manera predeterminada el usuario es "admin" y contraseña la "admin"

```bash
#Con -p en el puerto que se quiere servir
#Con -U se indica el usuario y contraseña
gosf server ~/Documentos/ -p 8081 -U user,pass

Output:
2022/03/23 22:59:01 Flags:
2022/03/23 22:59:01 - Servidor abierto en: 8081
2022/03/23 22:59:01 - Ruta servida: /home/user/Documentos
2022/03/23 22:59:01 - Tiempo de cierre automático: 0s
2022/03/23 22:59:01 - Se esta usando el directorio "/home/user/.config/gosf/static" para templates
2022/03/23 22:59:01 - Usuario de administrador: map[user:pass]
2022/03/23 22:59:01 Iniciando servidor...
```

Para apagar el servidor simplemente hay que hacer una señal de interrupción, un simple Ctrl+C(^C)

## Instalación (only Linux)

La instalación consta de darle permiso de ejecución al binario y moverlo al el directorio de binarios, en Linux es /usr/bin, y además mover la carpeta que contiene el HTML, CSS y JS (static) a la carpeta de configuración ($HOME/.config/gosf/). Esto ya lo hace automaticamente el install.sh:

- Clonar el repositorio

```bash
    git clone https://github.com/Joacohbc/gosf.git
```

- Ingresar a la carpeta

```bash
    cd gosf
```

- Ejecutar el instalador

```bash
    sh ./install.sh
```
