#!/bin/sh

MESSAGE="ej3"

RESPONSE=$(docker run --rm --network tp0_testing_net busybox sh -c "
    printf '%s\n' '$MESSAGE' | nc -w 2 server 12345
" 2>/dev/null | tr -d '\r\n')

if [ "$RESPONSE" = "$MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi