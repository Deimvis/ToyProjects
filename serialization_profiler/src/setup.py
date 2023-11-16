#!/usr/bin/env python3

from setuptools import setup, find_packages, Command
import subprocess as sp

class BuildProtoCommand(Command):
    """
    A custom command to run the protobuf compiler and generate _pb2.py files
    """
    description = 'build protobuf files'
    user_options = []

    def initialize_options(self):
        pass

    def finalize_options(self):
        pass

    def run(self):
        protoc_command = [
            'protoc',
            '--python_out=.',
            'serialization/schemas/perf_test_case.proto'
        ]
        self.announce('Running command: %s' % ' '.join(protoc_command), level=3)
        sp.check_call(protoc_command)

setup(
    name='Serialization Performance Testing System',
    version='1.0',
    packages=find_packages(),
    install_requires=[
        'avro',
        'msgpack',
        'protobuf',
        'xml-marshaller',
    ],
    cmdclass={
        'build_proto': BuildProtoCommand
    },
    package_data={
        '': ['**/*_pb2.py', '**/*.avsc']
    },
)
