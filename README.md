# MyDigitalTwin2026

mdk 个人的数字孪生项目

## API

### 1. 记录 mdk 发言

POST /api/records

```bash
curl -X POST http://$base_url/api/records \
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
    "id": 123,
    "created_at": "2026-08-11T10:30:00+08:00",
    "raw_content": "再稍等一下，我先去洗个澡",
    "objective_context": "mdk 准备去洗澡（收工节奏，洗完后应可休息）。",
    "ai_analysis": "mdk 去洗澡=准备收工的信号。简短温暖支持，洗完可祷告+睡觉。",
    "tags": ["生活", "休息", "闲聊"]
  }
}
```

说明：`created_at` 由服务端生成，序列化为 `+08:00` 时区。

