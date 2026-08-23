package mcp

const ServerDescription = `mdk's personal digital twin: continuously records real utterances, objective context, AI analysis, and tags, and provides history queries and the current time. It does not prescribe a persona or chat style.`

const ServerInstructions = `This is mdk's personal memory and context service. Use the listed tools to read and write; do not guess unpublished tool names.

Currently available:
- probe_health: confirm the service is up and read the current server time (milliseconds, +08:00). Do not use this to query history or write records.
- time: read the server's current local time (Asia/Shanghai), including the date, clock time, and weekday in Chinese. Use this when you need the current time as context; do not use it as a health check or to read or write records.`
