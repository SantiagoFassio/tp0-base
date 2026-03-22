import socket
import logging
import threading
from common.utils import Bet, store_bets
from common.parser import parse_bet, parse_batch
from server.common.reader import recv_line, recv_batch

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        # To shutdown gracefully
        self._shutdown_event = threading.Event()
        # Timeout to unblock accept()
        self._server_socket.settimeout(1)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        logging.info('action: accept_connections | result: in_progress')

        idle_cycles = 0

        while not self._shutdown_event.is_set():
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)
                idle_cycles = 0
            except socket.timeout:
                idle_cycles += 1
                if idle_cycles >= 3:
                    break
            except OSError:
                break

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # Leer la primera linea para obtener el número de apuestas en el batch
            header, buffer = recv_line(client_sock)
            n = int(header.strip())

            # Leer el batch completo
            lines = recv_batch(client_sock, n, buffer)

            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')

            try:
                bets = parse_batch(n , lines)
                store_bets(bets)
            except ValueError:
                logging.error(f'action: apuesta_recibida | result: fail | cantidad: {str(n)}')
                raise ValueError('Invalid batch format.')

            logging.info(f'action: apuesta_recibida | result: success | cantidad: {str(n)}')

            response = "OK\n".encode('utf-8')
        except ValueError as e:
            response = "NOK\n".encode('utf-8')
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            total_sent = 0
            while total_sent < len(response):
                sent = client_sock.send(response[total_sent:])
                if sent == 0:
                    raise ConnectionError("Client closed the connection")
                total_sent += sent
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections
        """

        # Connection arrived
        
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    
    def shutdown(self):
        """
        Shutdown server when signal is received
        """
        self._shutdown_event.set()
        self._server_socket.close()
