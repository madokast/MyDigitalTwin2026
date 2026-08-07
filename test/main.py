import json
from pathlib import Path
import urllib.request
import urllib.parse
from datetime import datetime

BASE_URL = "http://localhost:29301"
TOKEN = "dt-20260807"

def request(api: str, headers: dict | None, data: dict | None) -> dict:
    url = f"{BASE_URL}{api}"
    req = urllib.request.Request(url)
    
    if headers is not None:
        for key, value in headers.items():
            req.add_header(key, value)
    
    if data is not None:
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

def check(file: str, name: str, case: dict):
    headers = case.get("headers", {})
    if headers.get("Authorization") == "{token}":
        headers["Authorization"] = f"Bearer {TOKEN}"

    expected = case["expected"]
    result = request(case["api"], headers, case.get("data", None))

    for key, value in expected.items():
        if key not in result:
            raise AssertionError(f"Missing key '{key}' in response: {result}")
        
        if value == "{time}":
            # 确定时间合法
            try:
                datetime.fromisoformat(result[key])
            except ValueError:
                raise AssertionError(f"Expected '{key}' to be a valid time, but got: {result[key]}")
        elif value == "{integer}":
            # 确定是数字
            if not isinstance(result[key], int):
                raise AssertionError(f"Expected '{key}' to be a integer, but got: {result[key]}")
        elif value == "{existed}":
            # 确定存在
            if result[key] is None:
                raise AssertionError(f"Expected '{key}' to exist, but got: {result[key]}")
        elif result[key] != value:
            raise AssertionError(f"Expected '{key}': {value}, but got: {result[key]}")

def main():
    # 获取同目录下所有 json
    json_files = list(Path(__file__).parent.glob("*.json"))
    for json_file in json_files:
        with open(json_file, "r", encoding="utf-8") as f:
            cases = json.load(f)
        for name, case in cases.items():
            try:
                check(json_file.name, name, case)
                print(f"{json_file.name} - {name} passed")
            except Exception as e:
                print(f"{json_file.name} - {name} failed: {e}")

if __name__ == "__main__":
    main()
