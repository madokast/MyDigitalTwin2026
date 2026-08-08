import re
import uuid
from collections import defaultdict

def solve_unique(data: dict | list):
    uv = defaultdict(lambda:str(uuid.uuid4()))
    _solve(data, uv)

def _solve(data: dict | list, uv:defaultdict):
    if isinstance(data, list):
        for i, e in enumerate(data):
            if isinstance(e, str):
                data[i] = _solve_str(e, uv)
            elif isinstance(e, (dict, list)):
                _solve(e, uv)
    elif isinstance(data, dict):
        for k, v in data.items():
            if isinstance(v, str):
                data[k] = _solve_str(v, uv)
            elif isinstance(v, (dict, list)):
                _solve(v, uv)

def _solve_str(s: str, uv:defaultdict) -> str:
    return re.sub(
        r"\{unique:(\d+)\}",
        lambda m: uv[int(m.group(1))],
        s,
    )

