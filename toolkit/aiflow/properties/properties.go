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
	ParameterName string
	// 节点名称
	Name string
	// 数据库ID
	SourceId string
	// 执行的SQL语句
	Sql string
	// WHERE条件列表
	Wheres []string
	// 参数输入
	Args string
	// 是否打印输入
	PrintInput bool
	// 是否打印输出
	PrintOutput bool
	// 是否保存输入
	InputHistory bool
	// 是否保存输出
	ResultHistory bool
}

// MessageItem 结构体
type MessageItem struct {
	// 变量名
	VarName string
	// 变量值
	VarValue string
}

// EndNodeProperties 结构体
type EndNodeProperties struct {
	// 输出内容列表
	Messages []MessageItem
}

type FunctionNodeProperties struct {
	ParameterName string
	Name          string
	ChannelId     string
	Model         string
	FuncCall      string
	Messages      []Message
	PrintInput    bool
	PrintOutput   bool
	InputHistory  bool
	ResultHistory bool
}

type HttpNodeProperties struct {
	ParameterName string
	Name          string
	Url           string
	ContentType   string
	Method        string
	Headers       string
	Cookies       string
	Body          string
	Expression    string
	DataKey       string
	PrintInput    bool
	PrintOutput   bool
	InputHistory  bool
	ResultHistory bool
}

type LineProperties struct {
	Name           string
	StartCondition string
	Expression     string
	EndCondition   string
	ForType        string
	ForData        string
	ForSize        int
}

type McpNodeProperties struct {
	ParameterName string
	Name          string
	ChannelId     string
	Model         string
	McpType       int
	Command       string
	Messages      []Message
	PrintInput    bool
	PrintOutput   bool
	InputHistory  bool
	ResultHistory bool
}

type VectorNodeProperties struct {
	ParameterName   string
	Name            string
	ChannelId       string
	Model           string
	CollectionId    string
	Filename        string
	Limit           int
	SimilarityScore float32
	Analyse         bool
	System          string
	Message         string
	PrintInput      bool
	PrintOutput     bool
	InputHistory    bool
	ResultHistory   bool
	PrintComplete   bool
}
