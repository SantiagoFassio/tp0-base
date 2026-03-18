MESSAGE="ej3"

RESPONSE=$(docker run --rm --network testing_net busybox sh -c "
    echo '$MESSAGE' | nc server 12345
" 2>/dev/null | tr -d '\r')

if [ "$RESPONSE" = "$MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi