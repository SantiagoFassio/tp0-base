import socket

def recv_line(sock: socket.socket) -> str:
    buffer = b""
    while b"\n" not in buffer:
        chunk = sock.recv(1024)
        if not chunk:
            raise ConnectionError("Client closed the connection")
        buffer += chunk
    
    msg_bytes, _, rest = buffer.partition(b"\n")
    return msg_bytes.decode('utf-8'), rest

def recv_batch(sock: socket.socket, n, initial_buffer=b""):
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
