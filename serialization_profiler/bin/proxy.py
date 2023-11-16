import json
import logging
import socket
import traceback

from serialization.worker import WorkersManager
from serialization.worker.modes import (
    AvroWorker,
    JsonWorker,
    MsgpackWorker,
    ProtobufWorker,
    PythonNativeWorker,
    XMLWorker,
)


class Proxy:
    def __init__(self, socket_, serialization_mode2address):
        self.socket = socket_
        self.serialization_mode2address = serialization_mode2address

    def poll(self):
        logging.info('Proxy: begin polling')
        while True:
            data, address = self.socket.recvfrom(512)
            logging.info(f'Proxy: received {len(data)} bytes from {address}')

            try:
                self.handle(data, address)
            except Exception as error:
                msg = f'Got exception:\n{error}\nTraceback:\n{traceback.format_exc()}'
                logging.error(msg)
                if address not in self.serialization_mode2address.values():
                    self.socket.sendto(msg.encode(), address)

    def handle(self, data, address):
        request = data.decode().strip()
        match request:
            case request if request.startswith('get_result'):
                self.handle_get_result(request, address)
            case request if request.startswith('result_json'):
                self.handle_result_json(request, address)
            case 'help' | 'info' | '?':
                self.handle_help(request, address)
            case _:
                raise RuntimeError(f'Bad request: {request}')
    
    def handle_get_result(self, request, address):
        mode = request.split(' ')[1]
        query = json.dumps({'client_address__0': address[0], 'client_address__1': address[1]})
        msg = f'get_result_json\n{query}'
        self.socket.sendto(msg.encode(), self.serialization_mode2address[mode])
        logging.debug(f'Redirected {request} request to worker: {self.serialization_mode2address[mode]}')
        logging.debug(f'Sent request {msg} to worker: {self.serialization_mode2address[mode]}')

    def handle_result_json(self, request, address):
        result = json.loads(request.split('\n')[1])
        answer = result['answer']
        client_address = (result['client_address__0'], result['client_address__1'])
        self.socket.sendto(answer.encode(), client_address)
        logging.debug(f'Received result from worker: {request}')

    def handle_help(self, request, address):
        msg = '\n'.join(['get_result {mode}', '', 'Available modes:'] + list(WorkersManager().mode2worker.keys()) + [''])
        self.socket.sendto(msg.encode(), address)

def create_udp_socket(host: str, port: int):
    udp_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    udp_socket.bind((host, port))
    return udp_socket


def address_arg2socket_address(address):
    host, port = address.split(':')
    return (host, int(port))

def run_proxy(args):
    udp_socket = create_udp_socket(args.host, args.port)
    serialization_mode2address = {
        AvroWorker.MODE: address_arg2socket_address(args.avro_worker_address),
        JsonWorker.MODE: address_arg2socket_address(args.json_worker_address),
        MsgpackWorker.MODE: address_arg2socket_address(args.msgpack_worker_address),
        ProtobufWorker.MODE: address_arg2socket_address(args.protobuf_worker_address),
        PythonNativeWorker.MODE: address_arg2socket_address(args.python_native_worker_address),
        XMLWorker.MODE: address_arg2socket_address(args.xml_worker_address),
    }
    logging.debug(serialization_mode2address)
    proxy = Proxy(udp_socket, serialization_mode2address)
    try:
        proxy.poll()
    except Exception as error:
        logging.error(f'Got error during polling:\n{error}\nTraceback:\n{traceback.format_exc()}')
