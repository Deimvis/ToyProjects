#!/usr/bin/env python3

import argparse
import logging
import os

from serialization.worker.manager import WorkersManager
from proxy import run_proxy
from worker import run_worker


logging.basicConfig(
    level=logging.INFO if not os.getenv('DEBUG') else logging.DEBUG,
    format='[%(asctime)s] %(levelname)s [%(name)s] %(message)s',
    datefmt='%d/%b/%Y %H:%M:%S',
)


def parse_args():
    common_parser = argparse.ArgumentParser(add_help=False)
    common_parser.add_argument('--host', type=str, default='0.0.0.0')
    common_parser.add_argument('--port', type=int, default=80)
    parser = argparse.ArgumentParser(description='Run server')
    subparsers = parser.add_subparsers()

    proxy = subparsers.add_parser('proxy', help='Run proxy', parents=[common_parser])
    proxy.set_defaults(run=run_proxy)
    proxy.add_argument('--avro-worker-address', default='127.0.0.1:1001')
    proxy.add_argument('--json-worker-address', default='127.0.0.1:1002')
    proxy.add_argument('--msgpack-worker-address', default='127.0.0.1:1003')
    proxy.add_argument('--protobuf-worker-address', default='127.0.0.1:1004')
    proxy.add_argument('--python_native-worker-address', default='127.0.0.1:1005')
    proxy.add_argument('--xml-worker-address', default='127.0.0.1:1006')

    worker = subparsers.add_parser('worker', help='Run worker', parents=[common_parser])
    worker.set_defaults(run=run_worker)
    worker.add_argument('mode', choices=WorkersManager().mode2worker.keys(), help='Serialization mode')

    return parser.parse_args()


def main():
    args = parse_args()
    args.run(args)


if __name__ == '__main__':
    main()
