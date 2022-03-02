import cargarDOM from './js/cargarDOM.js';

document.addEventListener('DOMContentLoaded', function() {
    cargarDOM().catch(error => {
        //Si ocurrio un error
        alert('Ocurrio un error:' + error);
        console.error(error);
    });
});

