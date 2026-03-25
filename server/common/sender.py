import socket

def send_response(sock: socket.socket, response: str):
    """
    Uses the client socket to send a message back to the client.
    """
    total_sent = 0
    while total_sent < len(response):
        sent = sock.send(response[total_sent:])
        if sent == 0:
            raise ConnectionError("Client closed the connection")
        total_sent += sent
        