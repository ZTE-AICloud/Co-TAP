import json
import re
import time
from typing import Any, Type, Union

import json5
from pydantic import BaseModel

CHINESE_CHAR_RE = re.compile(r'[\u4e00-\u9fff]')


def json_loads(text: str) -> dict:
    text = text.strip('\n')
    if text.startswith('```') and text.endswith('\n```'):
        text = '\n'.join(text.split('\n')[1:-1])
    try:
        return json.loads(text)
    except json.decoder.JSONDecodeError as json_err:
        try:
            return json5.loads(text)
        except ValueError:
            raise json_err
    except Exception as e:
        raise


def json_loads_str(text: str) -> dict | str:
    try:
        return json_loads(text)
    except Exception as e:
        print(f'Json load error: {e}\n Text: {text}')
        return text


def has_chinese_chars(data: Any) -> bool:
    text = f'{data}'
    return bool(CHINESE_CHAR_RE.search(text))


def generate_template_from_model(model: Type[BaseModel]) -> Union[dict, list, str]:
    """
    根据 Pydantic 模型生成 JSON 模板，使用字段描述作为示例值
    """
    if hasattr(model, "__fields__"):  # 是 Pydantic 模型
        result = {}
        for field_name, field_info in model.model_fields.items():
            # 处理 List[SubModel]
            if hasattr(field_info.annotation, "__origin__") and field_info.annotation.__origin__ is list:
                sub_model_type = field_info.annotation.__args__[0]
                # 如果列表元素是 Pydantic 模型，则递归
                if hasattr(sub_model_type, "__fields__"):
                    result[field_name] = [generate_template_from_model(sub_model_type)]
                # 如果列表元素是原始类型，使用字段本身的描述作为列表的示例值
                else:
                    # 对于原始类型列表，通常我们希望显示原始类型本身的示例。
                    # 描述是针对 *列表* 本身的，而不是它的元素。
                    # 所以，我们为元素类型生成一个示例值。
                    # 如果需要将列表的描述作为列表中的唯一元素，可以改为 [field_info.description]
                    example_value = f"example_{sub_model_type.__name__.lower()}"
                    result[field_name] = [field_info.description or example_value]  # 将描述作为列表的示例值，如果不存在则用默认
            else:
                result[field_name] = field_info.description or f"example_{field_name}"
        return result
    else:
        return f"example_{model.__name__.lower()}"


def retry_on_exception(max_retries=3, delay=5):
    def decorator(func):
        def wrapper(*args, **kwargs):
            retries = 0
            while retries < max_retries:
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    print(f"Request failed with exception: {e}. Retrying...")
                    retries += 1
                    time.sleep(delay)
            raise Exception(f"Failed after {max_retries} attempts.")

        return wrapper

    return decorator
