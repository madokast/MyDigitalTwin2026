# MyDigitalTwin2026

mdk 个人的数字孪生项目

## API

请从 mdk 处获取目前的 base_url 和 token，每个 API 均需要 token 权限。

### 1. 记录 mdk 发言

POST /api/records

```bash
curl -X POST $base_url/api/records \
  -H "Authorization: Bearer $token" \
  -H "Content-Type: application/json" \
  -d '{
    "raw_content": "再稍等一下，我先去洗个澡",
    "objective_context": "mdk 准备去洗澡（收工节奏，洗完后应可休息）。",
    "ai_analysis": "mdk 去洗澡=准备收工的信号。简短温暖支持，洗完可祷告+睡觉。",
    "tags": ["生活", "休息", "闲聊"]
  }'
```

- `raw_content`：mdk 当前对话说的原始内容，请原封不动记录，必填，不能为空
- `objective_context`：客观描述的对话上下文信息，必填，不能为空
- `ai_analysis`：AI 分析，必填，不能为空
- `tags`：标签数组，选填，元素不能为空白串、不能重复，缺省为 `[]`

**成功响应**（`201 Created`）：

```json
{
  "ok": true,
  "status": 201,
  "record": {
    "id": 2504,
    "created_at": "2026-08-11T10:30:00+08:00",
    "raw_content": "再稍等一下，我先去洗个澡",
    "objective_context": "mdk 准备去洗澡（收工节奏，洗完后应可休息）。",
    "ai_analysis": "mdk 去洗澡=准备收工的信号。简短温暖支持，洗完可祷告+睡觉。",
    "tags": ["生活", "休息", "闲聊"]
  }
}
```

说明：`created_at` 由服务端生成，序列化为 `+08:00` 时区。

### 2. 查询记录

GET /api/records

```bash
curl -X GET "$base_url/api/records?page=1&page_size=10&q=洗澡&q=休息&tag=生活&tag=闲聊&from=2026-08-01&to=2026-08-31" \
  -H "Authorization: Bearer $token"
```

- `page`：页码，从 1 开始，默认 1
- `page_size`：每页记录数，1~1000，默认 100
- `q`：模糊搜索关键词，匹配 `raw_content`/`objective_context`/`ai_analysis`/`tags` 任一字段（不区分大小写），可重复传入多个（须全部命中），最多 10 个，不能为空
- `tag`：标签精确过滤，记录必须包含所有传入的标签，可重复传入多个，最多 10 个，不能为空白、不能重复
- `from`：起始时间（含），RFC3339（如 `2026-08-01T00:00:00+08:00`）或日期（如 `2026-08-01`，即当天 00:00 +08:00），选填
- `to`：结束时间（不含），格式同 `from`，选填

**成功响应**（`200 OK`）：

```json
{
  "ok": true,
  "status": 200,
  "page": 1,
  "page_size": 10,
  "total": 2,
  "total_page": 1,
  "records": [
    {
      "id": 2504,
      "created_at": "2026-08-11T10:30:00+08:00",
      "raw_content": "再稍等一下，我先去洗个澡",
      "objective_context": "mdk 准备去洗澡（收工节奏，洗完后应可休息）。",
      "ai_analysis": "mdk 去洗澡=准备收工的信号。简短温暖支持，洗完可祷告+睡觉。",
      "tags": ["生活", "休息", "闲聊"]
    }
  ]
}
```

说明：`records` 按 `id` 升序排列；`total` 为满足条件的总记录数，`total_page` 为总页数。

### 3. 标签 CRUD

Path 主体均为 /api/records/{record_id}/tags

**获取标签**

```bash
curl -X GET $base_url/api/records/{record_id:2504}/tags \
  -H "Authorization: Bearer $token"
```

成功响应（`200 OK`）：

```json
{
  "ok": true,
  "status": 200,
  "tags": ["生活", "休息", "闲聊"]
}
```

**附加标签**

```bash
curl -X PUT $base_url/api/records/{record_id:2504}/tags/健身 \
  -H "Authorization: Bearer $token"
```

- `tag`：要附加的标签，必填，不能为空白；标签含特殊字符时需 URL 编码（如 `/` 编码为 `%2F`）

成功响应（`201 Created`，标签已存在时为 `200 OK`）：

```json
{
  "ok": true,
  "status": 201,
  "attached": true,
  "changed": true,
  "tags": ["生活", "休息", "闲聊", "健身"]
}
```

说明：重复附加同一标签返回 `200` 且 `changed` 为 `false`（幂等）。

**删除标签**

```bash
curl -X DELETE $base_url/api/records/{record_id:2504}/tags/健身 \
  -H "Authorization: Bearer $token"
```

成功响应（`200 OK`）：

```json
{
  "ok": true,
  "status": 200,
  "detached": true,
  "changed": true,
  "tags": ["生活", "休息", "闲聊"]
}
```

说明：删除不存在的标签返回 `200` 且 `detached`/`changed` 均为 `false`（幂等）。

