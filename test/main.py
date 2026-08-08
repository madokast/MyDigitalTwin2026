import re
import sys
import json
import time
from pathlib import Path
import urllib.request
import urllib.parse
from datetime import datetime

BASE_URL = "http://localhost:29301"
TOKEN = "dt-20260807"

def do_request(api: str, headers: dict | None, data: dict | None) -> dict:
    url = f"{BASE_URL}{api}"
    req = urllib.request.Request(url)
    
    if headers is not None:
        for key, value in headers.items():
            req.add_header(key, value)
    
    if data is not None:
        if data == "{empty}":
            req.data = b""
        else:
            req.data = json.dumps(data).encode()
    
    try:
        with urllib.request.urlopen(req) as response:
            response_data = response.read()
            return json.loads(response_data)
    except urllib.error.HTTPError as e:
        # 读取错误响应体
        error_body = e.read().decode('utf-8')
        try:
            return json.loads(error_body)
        except json.JSONDecodeError:
            return {"error": f"HTTP {e.code}", "message": error_body}
    except urllib.error.URLError as e:
        return {"error": "Connection error", "message": str(e)}

def request(r: dict):
    headers = r.get("headers", {})
    if headers.get("Authorization") == "{token}":
        headers["Authorization"] = f"Bearer {TOKEN}"
    return do_request(r["api"], headers, r.get("data", None))

def check(case: dict):
    for r in case.get("before", []):
        res = request(r)
        assert res["status"]//100 == 2, f"before {r} failed: {res}"

    expected = case["expected"]
    result = request(case)
    check_result(expected, result)

def check_result(expected: str | dict, result: str | dict):
    if isinstance(expected, str):
        if expected != result:
            raise AssertionError(f"Expected: {expected}, but got: {result}")
        return

    for key, value in expected.items():
        if key not in result:
            raise AssertionError(f"Missing key '{key}' in response: {result}")
        
        if value == "{time}": # 确定时间合法
            try: datetime.fromisoformat(result[key])
            except ValueError: 
                raise AssertionError(f"Expected '{key}' to be a valid time, but got: {result[key]}")
        elif value == "{integer}": # 确定是数字
            if not isinstance(result[key], int):
                raise AssertionError(f"Expected '{key}' to be a integer, but got: {result[key]}")
        elif value == "{existed}": # 确定存在
            if result[key] is None:
                raise AssertionError(f"Expected '{key}' to exist, but got: {result[key]}")
        elif value == "{array}": # 确定是数组
            if not isinstance(result[key], list):
                raise AssertionError(f"Expected '{key}' to be a array, but got: {result[key]}")
        elif (m := re.fullmatch(r"\{array:(\d+)\}", str(value))): # 固定长度的数组
            if not isinstance(result[key], list):
                raise AssertionError(f"Expected '{key}' to be a array, but got: {result[key]}")
            if len(result[key]) != int(m.group(1)):
                raise AssertionError(f"Expected '{key}' to be a array with length {m.group(1)}, but got: {result[key]}")
        elif isinstance(value, list):
            if not isinstance(result[key], list):
                raise AssertionError(f"Expected '{key}' to be a array, but got: {result[key]}")
            if len(result[key]) != len(value):
                raise AssertionError(f"Expected '{key}' to be a array with length {len(value)}, but got: {result[key]}")
            for e, r in zip(value, result[key]):
                check_result(e, r) # 递归检查数组元素
        elif isinstance(value, dict): # 递归检查字典
            if not isinstance(result[key], dict):
                raise AssertionError(f"Expected '{key}' to be a dict, but got: {result[key]}")
            check_result(value, result[key])
        elif result[key] != value:
            raise AssertionError(f"Expected '{key}': {value}, but got: {result[key]}")

def main():
    parent = Path(__file__).parent
    json_files = [parent / f for f in sys.argv[1:]]
    if not json_files:
        json_files = list(parent.glob("*.json"))

    all_pass = True
    for json_file in json_files:
        with open(json_file, "r", encoding="utf-8") as f:
            cases = json.load(f)
        for name, case in cases.items():
            try:
                start = time.time()
                check(case)
                print(f"{json_file} - {name} passed - {time.time() - start:.2f}s")
            except Exception as e:
                all_pass = False
                print(f"{json_file} - {name} failed: {e}")
    
    if all_pass:
        print("All tests passed")
    else:
        print("Some tests failed")

if __name__ == "__main__":
    main()
