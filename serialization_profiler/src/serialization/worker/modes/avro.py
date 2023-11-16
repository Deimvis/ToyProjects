import avro.schema
from avro.io import DatumReader, DatumWriter
from io import BytesIO
from os.path import dirname, join as pj
from pathlib import Path
from typing import Any, Union

from serialization.worker.base import WorkerBase


AVRO_SCHEMA_FILE_PATH = Path(pj(dirname(__file__), '../../schemas/perf_test_case.avsc')).resolve()


class AvroWorker(WorkerBase):
    MODE = 'avro'

    def __init__(self, schema=None):
        self.schema = None
        self.writer = None
        self.reader = None
        self.reset(schema)

    def serialize(self, data: Any) -> Union[str, bytes]:
        bytes_io = BytesIO()
        self.writer.write(data, avro.io.BinaryEncoder(bytes_io))
        return bytes_io.getvalue()

    def deserialize(self, data: Union[str, bytes]) -> Any:
        bytes_io = BytesIO(data)
        return self.reader.read(avro.io.BinaryDecoder(bytes_io))

    @property
    def perf_test_case(self) -> Any:
        with open(AVRO_SCHEMA_FILE_PATH, 'rb') as f:
            self.reset(avro.schema.parse(f.read()))
        return super().perf_test_case

    def reset(self, schema):
        self.schema = schema
        self.writer = DatumWriter(self.schema)
        self.reader = DatumReader(schema)
