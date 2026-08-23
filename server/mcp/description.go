package mcp

const ServerDescription = `mdk 的个人数字孪生：持续记录真实发言、客观上下文、AI 分析和标签，并提供历史查询与当前时间。不规定人格或聊天风格。`

const ServerInstructions = `这是 mdk 的个人记忆与上下文服务。用列出的工具读写，不要猜测未公开的工具名。

当前可用：
- probe_health：确认服务是否可用，并读取服务器当前时间（毫秒，+08:00）。不要用它查询历史或写入记录。
- time：读取服务器当前本地时间（Asia/Shanghai），含中文日期、时刻和星期。需要当前时刻做上下文时用它，不要用它探活或读写记录。`
