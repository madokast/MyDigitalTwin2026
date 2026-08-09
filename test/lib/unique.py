import re
import uuid
from collections import defaultdict

def resolve_unique(data: dict | list):
    uv = defaultdict(lambda:str(uuid.uuid4()))
    _resolve(data, uv)

def _resolve(data: dict | list, uv:defaultdict):
    if isinstance(data, list):
        for i, e in enumerate(data):
            if isinstance(e, str):
                data[i] = _resolve_str(e, uv)
            elif isinstance(e, (dict, list)):
                _resolve(e, uv)
    elif isinstance(data, dict):
        for k, v in data.items():
            if isinstance(v, str):
                data[k] = _resolve_str(v, uv)
            elif isinstance(v, (dict, list)):
                _resolve(v, uv)

def _resolve_str(s: str, uv:defaultdict) -> str:
    return re.sub(
        r"\{unique:(\d+)\}",
        lambda m: uv[int(m.group(1))],
        s,
    )

