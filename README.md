# MyDigitalTwin2026

mdk 个人的数字孪生项目

## API

请从 mdk 处获取目前的 base_url 和 token，每个 API 均需要 token 权限。

### 0. 获取当前时间

GET /api/time

```bash
curl -X GET "$base_url/api/time" \
  -H "Authorization: Bearer $token"
```

**成功响应**（`200 OK`）：

```json
{
  "ok": true,
  "status": 200,
  "datetime": "2026-08-11T10:30:00+08:00",
  "timezone": "Asia/Shanghai",
  "local": {
    "date": "2026年8月11日",
    "time": "10点30分00秒",
    "weekday": "星期二"
  }
}
```

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

- `raw_content`：mdk 当前对话说的原始内容，请原封不动记录，必填，不能为空白串
- `objective_context`：客观描述的对话上下文信息，必填，不能为空白串
- `ai_analysis`：AI 分析，必填，不能为空白串
- `tags`：标签数组，选填，元素会 trim space，不能为空白串、不能重复，缺省为 `[]`

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

## 4. 使用方法

### 4.1 新建聊天：初始化mdk  mdk 上下文

每次创建新的聊天会话时，AI 应首先向 mdk 获取当前 MyDigitalTwin2026 的：

* `base_url`
* `token`

mdk 提供后，按照以下顺序初始化上下文。

#### 第一步：获取 mdk 最早的 200 条记录

调用：

```http
GET /api/records?page=1&page_size=200
```

读取最早的 200 条记录，综合其中的 `raw_content`、`objective_context`、`ai_analysis` 和 `tags`，建立对 mdk 长期情况的初步认识。

重点了解：

* mdk 的基本背景与长期特征
* 兴趣、习惯与偏好
* 工作、学习、生活等长期状态
* mdk 过去反复出现的问题
* mdk 与 AI 的长期交流特点
* 其他对后续聊天具有持续参考价值的信息

基于这些记录，在**当前会话内部**形成一份 mdk 画像。

不要求将 mdk 画像再次写入 MyDigitalTwin2026，除非 mdk 明确要求记录。

#### 第二步：获取最近的记录

继续查询记录，获取 mdk 最新的记录。

由于 `/api/records` 默认按照 `id` 升序返回，因此应根据上一步响应中的 `total` 和 `total_page` 获取最新的记录，例如：

```http
GET /api/records?page={total_page-1}&page_size=200
GET /api/records?page={total_page}&page_size=200
```

综合分析这些记录，重点了解 mdk **近期生活状态**：

* 最近正在做什么
* 最近关注什么问题
* 最近的工作、学习和生活状态
* 最近发生的重要事件
* 最近形成或改变的计划
* 最近持续讨论的话题
* 最近的情绪、行为和生活节奏变化

形成一份“mdk 最新生活情况”的临时上下文。

> 如果记录总数不足 500 条，请直接获取全部记录。

#### 第三步：获取当前时间

调用：

```http
GET /api/time
```

获取：

* 当前日期
* 当前时间
* 星期
* 时区

结合 mdk 长期画像、近期生活情况以及当前时间，推断 mdk 此刻可能处于什么生活场景。

例如：

* 工作日晚上 → 可能已经下班，进入休息时间
* 深夜 → 可能准备睡觉
* 工作日上午 → 可能正在工作或学习
* 周末下午 → 可能处于休闲、外出或个人事务时间

这种推断仅作为生成回复的上下文，**不得将推断当成确定事实**。

#### 第四步：开始聊天

综合：

```text
mdk 长期画像
        ↓
mdk 近期生活情况
        ↓
当前日期 / 时间 / 星期
        ↓
当前聊天内容
        ↓
生成个性化回复
```

生成聊天开场白。

开场白应根据上述上下文自然生成，而不是机械地向 mdk 复述画像。

例如，如果当前时间、近期记录和历史记录共同表明mdk 刚结束一天的工作，可以自然地从其当前状态切入，而不必直接说：

> “根据我的数据库，你今天应该已经下班了。”

---

### 4.2 聊天过程中：实时记录 mdk 发言

在聊天过程中，mdk **每发送一条消息**，都应先对该消息进行分析，并调用：

```http
POST /api/records
```

将 mdk 发言记录到 MyDigitalTwin2026。

其中：

* `raw_content`：mdk 当前消息，**必须原封不动保存**
* `objective_context`：当前消息对应的客观上下文
* `ai_analysis`：AI 对该消息的分析
* `tags`：根据实际内容添加适当标签

例如 mdk 说：

```text
今天终于把那个项目收尾了，累死我了
```

可以记录为：

```json
{
  "raw_content": "今天终于把那个项目收尾了，累死我了",
  "objective_context": "mdk 表示今天完成了一个项目的收尾工作，并表达了明显的疲劳感。",
  "ai_analysis": "mdk 当前可能处于项目完成后的放松阶段；后续回复可以关注其休息状态，但不应过度解读其情绪。",
  "tags": ["工作", "项目", "休息"]
}
```

记录成功后，再生成对 mdk 的正常回复。

**记录失败不应阻塞正常聊天。** 如果 API 暂时不可用，AI 应继续回答mdk ，而不是因为记录失败而中断对话。

---

### 4.3 聊天风格

MyDigitalTwin2026 **不规定具体聊天风格、语气、人格或回答方式**。

具体聊天方式完全由 mdk 在当前 AI 系统中的设定决定。

MyDigitalTwin2026 的职责主要是：

```text
提供mdk 历史信息
        +
持续记录mdk 当前信息
        +
提供当前时间
        ↓
帮助 AI 更准确地理解mdk 
```

而不是规定 AI 应该：

* 温柔
* 严肃
* 活泼
* 专业
* 像朋友
* 像秘书
* 像心理咨询师

具体行为由上层 AI 的系统提示词、mdk 设置以及当前对话决定。

---

### 4.4 联网搜索

推荐在聊天过程中，根据当前 mdk 消息和上下文，**适当搜索与 mdk 当前话题相关的最新信息**。

例如 mdk 正在讨论：

* 最新科技新闻
* 软件版本
* 当前产品价格
* 历史事件的最新研究
* 当地商家
* 体育比赛
* 政策或时事
* 当前网络服务状态

则可以主动搜索相关信息，再结合 mdk 的历史兴趣和近期情况进行回答。

搜索的目的不是单纯提供搜索结果，而是：

```text
mdk 当前问题
    +
mdk 长期兴趣
    +
mdk 近期状态
    +
最新外部信息
    ↓
更加个性化的回答
```

但是否搜索、搜索什么以及如何使用搜索结果，均由 AI 根据具体聊天内容自主决定。

---

### 4.5 整体工作流程

```text
新建聊天
   │
   ├── 向 mdk 获取 base_url + token
   │
   ├── GET /api/records?page=1&page_size=100
   │       ↓
   │   建立长期mdk 画像
   │
   ├── GET /api/records?page=<最后一页>&page_size=100
   │       ↓
   │   分析mdk 近期生活情况
   │
   ├── GET /api/time
   │       ↓
   │   获取当前时间
   │
   └── 综合以上信息
           ↓
       个性化开场白
           ↓
       开始聊天
           │
           ▼
    mdk 发送一条消息
           │
           ├── 分析消息
           │
           ├── POST /api/records
           │
           ├── 根据需要搜索外部信息
           │
           └── 生成回复
           │
           ▼
       mdk 继续聊天
           │
           └────── 循环
```

