import msgpack
from typing import Any, Union

from serialization.worker.base import WorkerBase


class MsgpackWorker(WorkerBase):
    MODE = 'msgpack'

    def serialize(self, data: Any) -> Union[str, bytes]:
        return msgpack.dumps(data)

    def deserialize(self, data: Union[str, bytes]) -> Any:
        return msgpack.loads(data)
