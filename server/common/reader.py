import socket

def recv_line(sock: socket.socket) -> str:
    """
    Receives data from the socket until a newline character is encountered.
    Returns the line as a string (without the newline).
    Also returns any remaining data extracted.
    If the client closes the connection before sending a newline, a ConnectionError is raised.
    """
    buffer = b""
    while b"\n" not in buffer:
        chunk = sock.recv(1024)
        if not chunk:
            raise ConnectionError("Client closed the connection")
        buffer += chunk
    
    msg_bytes, _, rest = buffer.partition(b"\n")
    return msg_bytes.decode('utf-8'), rest

def recv_batch(sock: socket.socket, n, initial_buffer=b""):
    """
    Receives a batch of n lines from the socket. Ocurrs when signal "B" is received.
    Uses an initial buffer to handle cases where part of the batch has already been received.
    Returns a list of lines as strings (without the newline).
    If the client closes the connection before sending n lines, a ConnectionError is raised.
    """
    buffer = initial_buffer
    lines = []

    while len(lines) < n:
        if b"\n" in buffer:
            line_bytes, _, buffer = buffer.partition(b"\n")
            lines.append(line_bytes.decode('utf-8'))
        else:
            chunk = sock.recv(1024)
            if not chunk:
                raise ConnectionError("Client closed the connection")
            buffer += chunk

    return lines
