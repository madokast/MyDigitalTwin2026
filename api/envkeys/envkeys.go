package envkeys

const ServerPort = "DT_SERVER_PORT"
const DatabaseUrl = "DT_DATABASE_URL"
const Token = "DT_TOKEN"

const QQBotAppID = "DT_QQBOT_APP_ID"
const QQBotAppSecret = "DT_QQBOT_APP_SECRET"
const QQBotUserOpenID = "DT_QQBOT_USER_OPENID"

const TestMode = "DT_TEST_MODE" // 测试模式

const MCP = "DT_MCP" // mcp 服务

var MustHave = []string{
	ServerPort,
	Token,
	DatabaseUrl,
	QQBotAppID,
	QQBotAppSecret,
	QQBotUserOpenID,
}
