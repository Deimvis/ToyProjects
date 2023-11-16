import pickle
from typing import Any, Union

from serialization.worker.base import WorkerBase


class PythonNativeWorker(WorkerBase):
    MODE = 'python_native'

    def serialize(self, data: Any) -> Union[str, bytes]:
        return pickle.dumps(data)

    def deserialize(self, data: Union[str, bytes]) -> Any:
        return pickle.loads(data)
