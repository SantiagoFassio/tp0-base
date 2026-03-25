import socket
import logging
from common.utils import Bet, has_won, store_bets, load_bets
from common.parser import parse_bet, process_header, parse_winners
from common.reader import recv_line, recv_batch
from common.sender import send_response
from common.result import Result
from multiprocessing import Manager, Process, Semaphore
import signal

DEFAULT_MAX_WORKERS = 10

class Server:
    def __init__(self, port, listen_backlog, agencies):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        self._server_socket.settimeout(1)

        self._total_agencies = agencies

        manager = Manager()
        self._agencies_done = manager.dict()
        self._lock = manager.Lock()

        self._max_workers = min(DEFAULT_MAX_WORKERS, agencies)
        self._sem = manager.Semaphore(self._max_workers)

        #Graceful shutdown event
        self._is_running = True
        self._children = []

        signal.signal(signal.SIGTERM, self.__handle_shutdown)
        signal.signal(signal.SIGINT, self.__handle_shutdown)
    
    def __handle_shutdown(self, signum, frame):
        logging.info(f"action: shutdown | result: in_progress")
        self._is_running = False
        self._server_socket.close()

    def run(self):
        """
        Server Loop
        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        logging.info('action: accept_connections | result: in_progress')

        while self._is_running:
            idle_cycles = 0
            try:
                client_sock = self.__accept_new_connection()

                idle_cycles = 0
                self._sem.acquire()

                p = Process(
                    target=self.__handle_client_connection_wrapper,
                    args=(client_sock,)
                )
                p.start()
                self._children.append(p)
                client_sock.close()
            except OSError:
                break
            except socket.timeout:
                idle_cycles += 1
                if idle_cycles >= 5:
                    break
        
        for p in self._children:
            p.join(timeout = 1)
            if p.is_alive():
                p.terminate()
        logging.info(f"action: shutdown | result: success")

    def process_signal(self, client_sock: socket.socket, signal: str, n: int, initial_buffer=b""):
        if signal == "B":
            lines = recv_batch(client_sock, n, initial_buffer)

            valid_bets = []
            invalid = 0

            for line in lines:
                try:
                    bet = parse_bet(line)
                    valid_bets.append(bet)
                except ValueError:
                    invalid += 1

            if len(valid_bets) > 0:
                with self._lock:
                    store_bets(valid_bets)
            
            if invalid > 0:
                return Result(False)
            
            return Result(True)
            
        elif signal == "E":
            with self._lock:
                if n < 1 or n > self._total_agencies:
                    return Result(False, "Invalid agency")
                elif n in self._agencies_done:
                    return Result(False, "Duplicate agency")
                self._agencies_done[n] = True
            
                logging.info(f'action: agencia_finalizada | result: success | agencia: {str(n)}')
                if len(self._agencies_done) == self._total_agencies:
                    logging.info(f'action: sorteo | result: success')

                return Result(True)
        
        elif signal == "G":
            if n < 1 or n > self._total_agencies:
                return Result(False, "Invalid agency")
            with self._lock:
                if len(self._agencies_done) != self._total_agencies:
                    return Result(False, "Not all agencies have completed their batches")
            
            bets = load_bets()
            winners = [bet for bet in bets if bet.agency == n and has_won(bet)]

            return Result(True, payload=winners)
        
        else:
            return Result(False, "Invalid signal")
        
    def __handle_client_connection_wrapper(self, client_sock):
        try:
            self.__handle_client_connection(client_sock)
        finally:
            self._sem.release()

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        response = b"NOK\n"
        try:
            # Leer la primera linea para obtener el número de apuestas en el batch
            header, buffer = recv_line(client_sock)
            addr = client_sock.getpeername()

            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')

            signal, n = process_header(header)

            result = self.process_signal(client_sock, signal, n, buffer)

            if not result.ok:
                if signal == "B":
                    logging.error(f'action: apuesta_recibida | result: fail | cantidad: {str(n)}')
                raise ValueError(result.message)
            
            if signal == "B":
                logging.info(f'action: apuesta_recibida | result: success | cantidad: {str(n)}')

            elif signal == "G":
                winners = result.payload or []
                response = parse_winners(winners)

            if signal != "G":
                response = b"OK\n"

        except Exception as e:
            response = b"NOK\n"
        finally:
            send_response(client_sock, response)
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections
        """
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
