# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |

### Servidor

Se trata de un "echo server", en donde los mensajes recibidos por el cliente se responden inmediatamente y sin alterar. 

Se ejecutan en bucle las siguientes etapas:

1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor retorna al paso 1.


### Cliente
 se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma:
 
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
4. Servidor responde al mensaje.
5. Servidor desconecta al cliente.
6. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

### Ejemplo

Al ejecutar el comando `make docker-compose-up`  y luego  `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```


## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).


### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `


### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).


### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

## Condiciones de Entrega
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios y que aprovechen el esqueleto provisto tanto (o tan poco) como consideren necesario.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deberán existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.
 (hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Se espera que se redacte una sección del README en donde se indique cómo ejecutar cada ejercicio y se detallen los aspectos más importantes de la solución provista, como ser el protocolo de comunicación implementado (Parte 2) y los mecanismos de sincronización utilizados (Parte 3).

Se proveen [pruebas automáticas](https://github.com/7574-sistemas-distribuidos/tp0-tests) de caja negra. Se exige que la resolución de los ejercicios pase tales pruebas, o en su defecto que las discrepancias sean justificadas y discutidas con los docentes antes del día de la entrega. 

El incumplimiento de las pruebas es condición de desaprobación, pero su cumplimiento no es suficiente para la aprobación.  Se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección informados  [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).
Respetar el formato y contenido las entradas de logs descritas en los ejercicios, pues son las que se chequean en cada uno de los tests.

# Seccion personal

# Parte 1

## Ejercicio 1

Se creo un archivo generar-compose.sh que genera un archivo .yaml. Este archivo yaml que crea los siguientes containers:
- server: crea el servidor, que acepta conexiones de clientes y recibe mensajes
- n clientes: crea n clientes (n siendo parametro al llamar al generador).
Estos clientes envian un mensaje al cliente y se desconectan.
### Generar al archivo
En la carpeta principal:
```bash
./generar-compose.sh {nombre-del-archivo} {n}
```
Como parametros, se envian el nombre del archivo (.yaml) y n siendo la cantidad de clientes.

## Ejercicio 2

Para evitar reconstruir las imagenes al modificar la configuracion, se sacaron las configuraciones de las
imagenes docker y se inyectan en el tiempo de ejecucion mediante un volumen.

Se agregaron volumenes en el generador de compose para que los archivos de configuracion se monten dentro
del contenedor. Cualquier cambio local se refleja. Ademas, se elimino que se copie la configuracion del cliente en el dockerfile, evitando que la configuracion quede ligada a la imagen.

La forma de ejecucion no cambia respecto al Ejercicio 1.

## Ejercicio 3

Se creo el script "validar-echo-server.sh". Con la imagen del servidor levantada, ejecutar en la carpeta principal del proyecto:
```bash
./validar-echo-server.sh
```
Este script levanta un contenedor temporal que se conecta a la red del servidor. Mediante netcat interno (no es requerido instalar netcat) Se comunica con el servidor para testear la capacidad de echo server, limpiando y comparando la respuesta.

Escribe un mensaje particular segun si lo que recibe es igual o no a lo que se envio (correspondiente a un echo server)

## Ejercicio 4

Para detener el loop inicial, se creo un evento usando threading (en el caso del servidor) que espera una senal de
SIGTERM, al momento de detectar la senal, se avisa al loop del servidor y se termina de aceptar clientes. Tambien se agrego un timeout para futuros ejercicios (parte 2 y 3).

En el caso del cliente, se uso Notify para detectar el SIGTERM en go y frenar el cliente de continuar el loop.

# Parte 2

## Ejercicio 5

Se implemento en el yaml base (al generar el docker compose yaml (ej1)) con los datos de un cliente.

El cliente ahora construye un objeto bet, lo serializa en un mensaje y lo envia al servidor.

El servidor ahora recibe el mensaje de cliente, verifica la integridad, crea un objeto Bet y lo guarda usando los comandos de utils.py

Tanto cliente como servidor tienen en cuenta el short write y short read y 
tienen en cuenta enviar los datos hasta que no haya mas por enviar y leen 
hasta que encuentran el delimitador de mensaje (\n).

Se agrego que el servidor, al no recibir nuevas conexiones por 
un tiempo determinado (3 segundos), se corte definitivamente.

### Protocolo

Request:
- Formato: String Delimitado por el caracter "|"
- Terminacion: Newline (\n)

Campos: Todos tipo string en orden
- nombre
- apellido
- dni (XXXXXXX)
- nacimiento (YYYY-MM-DD)
- numero

Respuesta:
- String
- Campo unico: status (OK o NOK segun resultado)

Se envia NOK si hay un error al parsear la Bet, y se descarta.

## Ejercicio 6

El cliente ya no loopea un "loop amount" de veces, sino que StartClientLoop fue reemplazado por sendBatch().
Esta funcion tiene el objetivo de encargar a leer los archivos (ver batch_reader.go), 
serializar cada bet correctamente (ver serializer.go) y enviar el batch de bets al servidor.

batch_reader.go resuelve de forma interesante la limitacion de maxAmount y 8kb de datos. El lector lee el
csv correspondiente de su cliente, hasta llegar a leer maxAmount bets, o al llegar a 8kb de datos,
lo que ocurra antes. Para evitar perder datos, guarda un buffer con la ultima bet que no pudo ser agregada
al mensaje, y espera a que el cliente vuelva a loopear, esperando un nuevo batch de datos.

El servidor, por su parte, se encarga de recibir primero una linea. 
Esta linea forma parte del nuevo protocolo y contiene la cantidad de bets enviadas en el batch.
Luego, el servidor recibe datos hasta recibir la cantidad de bets esperadas.

### Protoolo

La diferencia en el cambio del protocolo con el ejercicio 5 es la presencia de "x\n" al inicio del mensaje
del cliente. x representa la cantidad de bets que van a ser enviadas en ese batch de bets.

un ejemplo se ve:

```
2\n
agencia|nombre|apellido|dni|nacimiento|numero\n
agencia|nombre|apellido|dni|nacimiento|numero\n
```

El servidor verifica todas las bets enviadas parseandolas en Bets. Si hay un error en una de ellas,
la ignora y devuelve un mensaje de error al cliente.
Notas:
- El cliente ignora si recibe un mensaje de OK o NOK. Solo envia hasta quedarse sin bets para enviar
(o un SIGTERM)
- El servidor guarda las bets que se encuentren en el mismo batch que una bet incorrecta. Simplemente
no se agrega la bet erronea a la lista.

## Ejercicio 7

Se cambio el protocolo para acomodar que el cliente pueda notificar al servidor el tipo de mensaje
que se envia de parte del cliente. Se agregaron cambios de protocolo

### Protocolo

El cambio con respecto a ejercicios anteriores es el hecho de que se agrega una letra al inicio del mensaje,
creando un "header" para el servidor. Este header contiene la letra correspondiente y un numero.
Las senales son:

#### B (Batch)

La senal corresponde al ya existente sistema de enviar batches del cliente respecto al ejercicio anterior.

Ejemplo:
```
B|2\n
agencia|nombre|apellido|dni|nacimiento|numero\n
agencia|nombre|apellido|dni|nacimiento|numero\n
```

El servidor responde de la misma forma que el Ejercicio anterior.

#### E (End)

La senal indica al servidor que el cliente ya finalizo de enviar todos los batches de bets, 
y que no hay mas bets que enviar. El cliente envia este mensaje hasta recibir el OK del servidor,
pero recibir un NOK es un caso borde.

Ejemplo:
```
E|1\n
```
El cliente envia la senal y un numero, este numero siendo el numero de agencia del cliente.

El servidor contesta "OK\n" y agrega la agencia a una lista de agencias que ya terminaron. En caso
de que el cliente ya se encontrase en dicha lista, se envia "NOK\n", como caso borde.

#### G (Get)

La senal indica al servidor que el cliente de agencia x solicita los ganadores del concurso
correspondientes a su agencia. En caso de que no todos los clientes hayan enviado el mensaje 
con senal E (que indica que ya finalizaron de enviar bets), el servidor contesta "NOK\n".
El cliente, al recibir "NOK\n", hace polling hasta recibir los usuarios.

Ejemplo (cliente):
```
G|1\n
```
Siendo el numero (1 en el ejemplo) el numero de la agencia.

Ejemplo (servidor):
```
NOK\n
```
En caso de que todos los otros clientes tambien hayan terminado, se leen los ganadores y se envia una respuesta con sus documentos separados por "\n"
```
3\n
11111111\n
22222222\n
33333333\n
```

# Parte 3

## Ejercicio 8

Se implemento una solucion multiproceso. Es decir, cada nueva conexion aceptada por el servidor
genera un nuevo proceso (Process). Cada proceso maneja el ciclo completo de comunicacion con un cliente.
De esta forma, se evita el impacto de GIL de CPython, ya que todos los procesos tienen
su propia memoria.

Realizar multiprocessing hay problemas tanto de secciones criticas como de estados compartidos.
Principalmente, los problemas a resolver fueron del estado compartido de agencies_done (la cantidad de
agencias que terminaron de enviar bets) y store_bets, la funcion que guarda las bets enviadas por un cliente.

Mediante un monitor, se logro hacer que se pueda compartir el estado evitando race conditions. Tambien
logra limitar el acceso a la escritura de datos a solo uno a la vez, evitando la corrupcion de datos y dirty writes.

Nota: Este lock protege ambos recursos a la vez. Es de esta forma para simplificar la implementacion
sin penalizar de forma excesiva la performance.

A su vez, se implemento un semaforo para limitar la cantidad maxima de procesos concurrentes.
Por defecto, la cantidad de procesos simultaneos esta limitada al minimo entre el numero total de clientes
y un valor maximo (10 por defecto). Esto tiene el objetivo de que se cree una cantidad ilimitada
de procesos y que no se agoten los recursos (CPU, memoria).

# Ejecucion

Habiendo provisto de los archivos necesarios para la cantidad correspondiente de clientes (n clientes, n archivos),
crear el docker-compose.yaml y luego hacer ejecutar las imagenes mediante el makefile:
```bash
./generar-compose.sh {nombre-del-archivo} {n}
make docker-compose-up
```
Para observar los logs:
```bash
make docker-compose-logs
```
Para apagar todos los contenedores:
```bash
make docker-compose-down
```
