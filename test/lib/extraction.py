import re

def extract(res: dict, extractions: list | None) -> dict:
    extracted = {}
    for ks in (extractions or []):
        data = res
        for key in ks.split('.'):
            data = data[key]
        extracted[f"{{extracted:{ks}}}"] = data
    return extracted

def resolve_extracted(data: dict | list, extracted: dict):
    if isinstance(data, list):
        for i, e in enumerate(data):
            if isinstance(e, str):
                data[i] = _resolve_str(e, extracted)
            elif isinstance(e, (dict, list)):
                resolve_extracted(e, extracted)
    elif isinstance(data, dict):
        for k, v in data.items():
            if isinstance(v, str):
                data[k] = _resolve_str(v, extracted)
            elif isinstance(v, (dict, list)):
                resolve_extracted(v, extracted)

def _resolve_str(s: str, extracted: dict) -> any:
    if s in extracted:
        return extracted[s]

    return re.sub(
        r"(\{extracted:[\w.]+\})",
        lambda m: str(extracted[m.group(1)]),
        s,
    )
