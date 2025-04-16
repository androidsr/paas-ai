package properties

type Message struct {
	// 消息类型（human：用户消息，system：系统消息，ai：AI响应）
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
}

// NewMessage 构造一个新的 Message 实例
func NewMessage(messageType, message string) *Message {
	return &Message{
		MessageType: messageType,
		Message:     message,
	}
}

type ScriptNodeProperties struct {
	// 参数名称
	ParameterName string `json:"parameterName"`
	// 节点名称
	Name string `json:"name"`
	// 代码块内容
	Content string `json:"content"`
	// 是否打印输入
	PrintInput bool `json:"printInput"`
	// 是否打印输出
	PrintOutput bool `json:"printOutput"`
	// 是否保存输入
	InputHistory bool `json:"inputHistory"`
	// 是否保存输出
	ResultHistory bool `json:"resultHistory"`
}

type LlmNodeProperties struct {
	// 参数名称
	ParameterName string `json:"parameterName"`
	// 节点名称
	Name string `json:"name"`
	// AI平台的ID
	ChannelId string `json:"channelId"`
	// 选择的模型
	Model string `json:"model"`
	// 消息列表
	Messages []Message `json:"messages"`
	// 是否打印输入
	PrintInput bool `json:"printInput"`
	// 是否打印输出
	PrintOutput bool `json:"printOutput"`
	// 是否保存输入
	InputHistory bool `json:"inputHistory"`
	// 是否保存输出
	ResultHistory bool `json:"resultHistory"`
	// 结束打印
	PrintComplete bool `json:"printComplete"`
}

type DatabaseNodeProperties struct {
	// 参数名称
	ParameterName string `json:"parameterName"`
	// 节点名称
	Name string `json:"name"`
	// 数据库ID
	SourceId string `json:"sourceId"`
	// 执行的SQL语句
	Sql string `json:"sql"`
	// WHERE条件列表
	Wheres []string `json:"wheres"`
	// 参数输入
	Args string `json:"args"`
	// 是否打印输入
	PrintInput bool `json:"printInput"`
	// 是否打印输出
	PrintOutput bool `json:"printOutput"`
	// 是否保存输入
	InputHistory bool `json:"inputHistory"`
	// 是否保存输出
	ResultHistory bool `json:"resultHistory"`
}

// MessageItem 结构体
type MessageItem struct {
	// 变量名
	VarName string `json:"varName"`
	// 变量值
	VarValue string `json:"varValue"`
}

// EndNodeProperties 结构体
type EndNodeProperties struct {
	// 输出内容列表
	Messages []MessageItem `json:"messages"`
}

type FunctionNodeProperties struct {
	ParameterName string    `json:"parameterName"`
	Name          string    `json:"name"`
	ChannelId     string    `json:"channelId"`
	Model         string    `json:"model"`
	FuncCall      string    `json:"funcCall"`
	Messages      []Message `json:"messages"`
	PrintInput    bool      `json:"printInput"`
	PrintOutput   bool      `json:"printOutput"`
	InputHistory  bool      `json:"inputHistory"`
	ResultHistory bool      `json:"resultHistory"`
}

type HttpNodeProperties struct {
	ParameterName string `json:"parameterName"`
	Name          string `json:"name"`
	Url           string `json:"url"`
	ContentType   string `json:"contentType"`
	Method        string `json:"method"`
	Headers       string `json:"headers"`
	Cookies       string `json:"cookies"`
	Body          string `json:"body"`
	Expression    string `json:"expression"`
	DataKey       string `json:"dataKey"`
	PrintInput    bool   `json:"printInput"`
	PrintOutput   bool   `json:"printOutput"`
	InputHistory  bool   `json:"inputHistory"`
	ResultHistory bool   `json:"resultHistory"`
}

type LineProperties struct {
	Name           string `json:"name"`
	StartCondition string `json:"startCondition"`
	Expression     string `json:"expression"`
	EndCondition   string `json:"endCondition"`
	ForType        string `json:"forType"`
	ForData        string `json:"forData"`
	ForSize        int    `json:"forSize"`
}

type McpNodeProperties struct {
	ParameterName string    `json:"parameterName"`
	Name          string    `json:"name"`
	ChannelId     string    `json:"channelId"`
	Model         string    `json:"model"`
	McpType       int       `json:"mcpType"`
	Command       string    `json:"command"`
	Messages      []Message `json:"messages"`
	PrintInput    bool      `json:"printInput"`
	PrintOutput   bool      `json:"printOutput"`
	InputHistory  bool      `json:"inputHistory"`
	ResultHistory bool      `json:"resultHistory"`
}

type VectorNodeProperties struct {
	ParameterName   string  `json:"parameterName"`
	Name            string  `json:"name"`
	ChannelId       string  `json:"channelId"`
	Model           string  `json:"model"`
	CollectionId    string  `json:"collectionId"`
	Filename        string  `json:"filename"`
	Limit           int     `json:"limit"`
	SimilarityScore float32 `json:"similarityScore"`
	Analyse         bool    `json:"analyse"`
	System          string  `json:"system"`
	Message         string  `json:"message"`
	PrintInput      bool    `json:"printInput"`
	PrintOutput     bool    `json:"printOutput"`
	InputHistory    bool    `json:"inputHistory"`
	ResultHistory   bool    `json:"resultHistory"`
	PrintComplete   bool    `json:"printComplete"`
}
