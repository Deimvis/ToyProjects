import sys
from abc import ABC, abstractmethod
from typing import Any, Union

from .utils import timeit


class WorkerBase(ABC):

    @abstractmethod
    def serialize(self, data: Any) -> Union[str, bytes]:
        pass

    @abstractmethod
    def deserialize(self, data: Union[str, bytes]) -> Any:
        pass

    @staticmethod
    @abstractmethod
    def MODE():
        pass

    def calc_perf_info(self, iterations=1000) -> str:
        serialization_time_results = []
        deserialization_time_results = []
        for _ in range(iterations):
            obj = self.perf_test_case
            serialization_time, serialized_obj = timeit(self.serialize)(obj)
            deserialization_time, deserialized_obj = timeit(self.deserialize)(serialized_obj)  # noqa
            serialization_time_results.append(serialization_time)
            deserialization_time_results.append(deserialization_time)
        avg = lambda arr: sum(arr) / len(arr)
        return ' - '.join([
            self.MODE,
            str(sys.getsizeof(serialized_obj)),
            str(int(sum(serialization_time_results) * 1000)) + 'ms',
            str(int(sum(deserialization_time_results) * 1000)) + 'ms',
        ])

    @property
    def perf_test_case(self) -> Any:
        return {
            'string_key': 'string_value',
            'int_key': 42,
            'float_key': 13.37,
            'array_key': ['val1', 'val2', 'val3', 'val1', 'val2', 'val3', 'val1', 'val2', 'val3', 'val1', 'val2', 'val3', 'val1', 'val2', 'val3', 'val1', 'val2', 'val3', 'val1', 'val2', 'val3', 'val1', 'val2', 'val3', 'val1', 'val2', 'val3', 'val1', 'val2', 'val3'],
            'map_key': {
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'k': 'v',
                'hello': 'world',
            },
        }
