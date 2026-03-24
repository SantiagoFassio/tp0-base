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
      - AGENCIES=$NUM_CLIENTS
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
      - CLI_AGENCY=$i
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data:/data
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
