import json
from typing import Any, Union

from serialization.worker.base import WorkerBase


class JsonWorker(WorkerBase):
    MODE = 'json'

    def serialize(self, data: Any) -> Union[str, bytes]:
        return json.dumps(data)

    def deserialize(self, data: Union[str, bytes]) -> Any:
        return json.loads(data)
