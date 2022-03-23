##Colores##
ESC=$(printf '\033')
RESET="${ESC}[0m"
RED="${ESC}[31m"
GREEN="${ESC}[32m"

##Funciones con colores##
greenprint() { printf "${GREEN}%s${RESET}\n" "$1"; }
redprint() { printf "${RED}%s${RESET}\n" "$1"; }

error() {
    echo "$(redprint "$1")"
}

exito() {
    echo "$(greenprint "$1")"
}

mv ./bin/gosf.bin ./bin/gosf 
if [ $? -eq 0 ];then 
    exito "Nombre cambiado gosf.bin -> ./bin/gosf..." 
else
    error "Error al renombar el binario"
    exit 1
fi 

chmod +x ./bin/gosf 
if [ $? -eq 0 ];then 
    exito "Permisos de ejecucion dados..." 
else
    error "Error al dar los permisos de ejecucion"
    exit 1
fi 

sudo mv ./bin/gosf /usr/bin/ 
if [ $? -eq 0 ];then 
    echo "Archivo movido con exito..." 
else
    error "Error al mover el binario"
    exit 1
fi 

#Si $XDG_CONFIG_HOME esta vacia
if [ -z "$XDG_CONFIG_HOME" ]; then
    XDG_CONFIG_HOME="$HOME/.config"
    echo "La varaible \$XDG_CONFIG_HOME fue definida con en exito \"\$HOME/.config\""
fi

#Defino donde esta el tempalte dir
CONFIG_DIR="$XDG_CONFIG_HOME"/gosf/static/

#CCreo el directorio
mkdir -p $CONFIG_DIR
if [ $? -eq 0 ]; then
    exito "Directorio de configuracion creado exitosamente"
else
    error "El directorio de configuracion no se pudo crear"
    exit 1
fi

#Copio la carpetas de templates
cp -r ./src/static/* $CONFIG_DIR
if [ $? -eq 0 ]; then
    exito "Archivos de templates copiados exitosamente"
else
    error "Los archivos de templates no se pudieron copiar"
    exit 1
fi

echo "Instalado con exito"

