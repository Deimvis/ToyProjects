from xml_marshaller import xml_marshaller
from typing import Any, Union

from serialization.worker.base import WorkerBase


class XMLWorker(WorkerBase):
    MODE = 'xml'

    def serialize(self, data: Any) -> Union[str, bytes]:
        return xml_marshaller.dumps(data)

    def deserialize(self, data: Union[str, bytes]) -> Any:
        return xml_marshaller.loads(data)
