import avro.schema
from os.path import dirname, join as pj
from pathlib import Path

from serialization.worker.modes import (
    AvroWorker,
    JsonWorker,
    MsgpackWorker,
    ProtobufWorker,
    PythonNativeWorker,
    XMLWorker,
)


AVRO_SCHEMA_FILE_PATH = Path(pj(dirname(__file__), '../schemas/perf_test_case.avsc')).resolve()


class WorkersManager:
    def __init__(self):
        self.workers = [
            self.init_AvroWorker(),
            JsonWorker(),
            MsgpackWorker(),
            ProtobufWorker(),
            PythonNativeWorker(),
            XMLWorker(),
        ]
        self._mode2worker = {
            worker.MODE: worker for worker in self.workers
        }

    def init_AvroWorker(self):
        with open(AVRO_SCHEMA_FILE_PATH, 'rb') as f:
            return AvroWorker(avro.schema.parse(f.read()))

    @property
    def mode2worker(self):
        return self._mode2worker
 