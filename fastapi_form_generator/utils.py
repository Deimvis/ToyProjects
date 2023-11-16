import inspect
from fastapi import Form
from pydantic import BaseModel
from pydantic.fields import FieldInfo
from typing import Annotated, Dict, List, Type


def _unflatten_dict(d: Dict, sep: str) -> Dict:
    result = {}
    for k, v in d.items():
        node = result
        path = k.split(sep)
        for p in path[:-1]:
            if p not in node:
                node[p] = {}
            node = node[p]
        node[path[-1]] = v
    return result


class WithFormGen:
    def __init__(self, method_name: str = 'form', by_alias: bool = False, submodel_sep='__'):
        self.method_name = method_name
        self.by_alias = by_alias
        self.submodel_sep = submodel_sep

    def __call__(self, cls: Type[BaseModel]):
        async def generate_form(**flattened_data):
            return cls.model_validate(_unflatten_dict(flattened_data, sep=self.submodel_sep))

        signature = self.generate_form_signature(cls, prefix='')
        generate_form.__signature__ = inspect.signature(generate_form).replace(parameters=signature)  # noqa
        setattr(cls, self.method_name, generate_form)
        return cls

    def generate_form_signature(self, cls: Type[BaseModel], prefix: str = '') -> List[inspect.Parameter]:
        def build_annotation(field_info: FieldInfo) -> bool:
            return Annotated[field_info.annotation, *(field_info.metadata or []), Form()]

        def is_model(t: Type) -> bool:
            return inspect.isclass(t) and issubclass(t, BaseModel)

        signature = []
        for field_name, field_info in cls.model_fields.items():
            field_info: FieldInfo
            param_name = prefix + (field_info.alias if (self.by_alias and field_info.alias is not None) else field_name)
            if is_model(field_info.annotation):
                signature.extend(self.generate_form_signature(field_info.annotation, prefix=param_name+self.submodel_sep))
            else:
                # NOTE! If there is no default value for a field, then None will be set as a default value
                #       So expect validation error in case when field wasn't marked as `required` in html and wasn't set
                signature.append(
                    inspect.Parameter(
                        param_name,
                        inspect.Parameter.KEYWORD_ONLY,
                        default=field_info.get_default(call_default_factory=True),
                        annotation=build_annotation(field_info),
                    )
                )
        return signature
