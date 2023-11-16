import json
import logging
import socket
import traceback
from serialization.worker.base import WorkerBase
from serialization.worker import WorkersManager


class Worker:
    def __init__(self, socket_: socket.socket, serialization_worker: WorkerBase):
        self.socket = socket_
        self.serialization_worker = serialization_worker

    def poll(self):
        logging.info('Worker: begin polling')
        while True:
            data, address = self.socket.recvfrom(512)
            logging.info(f'Worker: received {len(data)} bytes from {address}')
            try:
                self.handle(data, address)
            except Exception as error:
                msg = f'Got exception:\n{error}\nTraceback:\n{traceback.format_exc()}'
                logging.error(msg)
                self.socket.sendto(msg.encode(), address)

    def handle(self, data, address):
        request = data.decode().strip()
        match request:
            case 'get_result':
                self.handle_get_result(request, address)
            case request if request.startswith('get_result_json'):
                self.handle_get_result_json(request, address)
            case _:
                raise RuntimeError(f'Bad request: {request}')

    def handle_get_result(self, request, address):
        perf_info = self.serialization_worker.calc_perf_info()
        self.socket.sendto((perf_info + '\n').encode(), address)

    def handle_get_result_json(self, request, address):
        request = json.loads(request.split('\n')[1])
        perf_info = self.serialization_worker.calc_perf_info()
        result = json.dumps({'answer': perf_info + '\n'} | request)
        self.socket.sendto(f'result_json\n{result}'.encode(), ('proxy', 2000))
        logging.debug(f'Respond on `{request}` to `("proxy", 2000)`')


def create_udp_socket(host: str, port: int):
    udp_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    udp_socket.bind((host, port))
    return udp_socket


def run_worker(args):
    udp_socket = create_udp_socket(args.host, args.port)
    serialization_worker = WorkersManager().mode2worker[args.mode]
    worker = Worker(udp_socket, serialization_worker)
    try:
        worker.poll()
    except Exception as error:
        logging.error(f'Got error during polling:\n{error}\nTraceback:\n{traceback.format_exc()}')
