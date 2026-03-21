OUTPUT_FILE=$1
NUM_CLIENTS=$2

cat <<EOF > $OUTPUT_FILE
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
EOF

for ((i=1; i<=NUM_CLIENTS; i++))
do
cat <<EOF >> $OUTPUT_FILE

  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
      - CLI_NOMBRE=Nombre
      - CLI_APELLIDO=Base
      - CLI_DNI=11111111
      - CLI_NACIMIENTO=1111-11-11
      - CLI_NUMERO=1
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
EOF
done

cat <<EOF >> $OUTPUT_FILE

networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
EOF

echo "Archivo $OUTPUT_FILE generado con $NUM_CLIENTS clientes."
echo "Nota del Ejercicio2: Se agregaron los volumenes para evitar la reconstruccion de las imagenes de docker"
