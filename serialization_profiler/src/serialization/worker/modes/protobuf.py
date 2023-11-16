from typing import Any, Union
from google.protobuf.json_format import ParseDict

from serialization.schemas.perf_test_case_pb2 import TPerfTestCase
from serialization.worker.base import WorkerBase


class ProtobufWorker(WorkerBase):
    MODE = 'protobuf'

    def __init__(self, message=None):
        self.message = message

    def serialize(self, data: Any) -> Union[str, bytes]:
        return data.SerializeToString()

    def deserialize(self, data: Union[str, bytes]) -> Any:
        msg = self.message()
        msg.ParseFromString(data)
        return msg

    @property
    def perf_test_case(self) -> Any:
        self.message = TPerfTestCase
        test_case_dict = super().perf_test_case
        return ParseDict(test_case_dict, TPerfTestCase())
