import avro.schema
from functools import partial
from google.protobuf.json_format import MessageToDict, ParseDict

from serialization.worker.modes import (
    AvroWorker,
    JsonWorker,
    MsgpackWorker,
    ProtobufWorker,
    PythonNativeWorker,
    XMLWorker,
)
from serialization.schemas.perf_test_case_pb2 import TPerfTestCase


def get_test_dict():
    return {
        'string_key': 'string_value',
        'int_key': 42,
        'float_key': 13.37,
        'array_key': ['val1', 'val2', 'val3'],
        'map_key': {
            'k': 'v',
            'hello': 'world',
        },
    }


def are_same(dict1, dict2):
    dict1['float_key'] = '{:.2f}'.format(dict1['float_key'])
    dict2['float_key'] = '{:.2f}'.format(dict2['float_key'])
    return dict1 == dict2


def run_smoke_test(worker, dict2obj=lambda x: x, obj2dict=lambda x: x):
    obj = dict2obj(get_test_dict())
    serialized_obj = worker.serialize(obj)
    deserialized_obj = worker.deserialize(serialized_obj)
    assert are_same(obj2dict(obj), obj2dict(deserialized_obj))


def test_avro():
    with open('../../../schemas/perf_test_case.avsc', 'rb') as f:
        worker = AvroWorker(avro.schema.parse(f.read()))
    run_smoke_test(worker)


def test_json():
    worker = JsonWorker()
    run_smoke_test(worker)


def test_msgpack():
    worker = MsgpackWorker()
    run_smoke_test(worker)


def test_protobuf():
    worker = ProtobufWorker(message=TPerfTestCase)
    run_smoke_test(worker, dict2obj=lambda x: ParseDict(x, TPerfTestCase()), obj2dict=partial(MessageToDict, preserving_proto_field_name=True))


def test_python_native():
    worker = PythonNativeWorker()
    run_smoke_test(worker)


def test_xml():
    worker = XMLWorker()
    run_smoke_test(worker)
