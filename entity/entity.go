package entity

type BaseEntity struct {
	Id string `json:"id" gorm:"primaryKey;column:id"`
}

type Config struct {
	BaseEntity
	Content        string `json:"content" gorm:"column:content"`
	LegalStatement string `json:"legalStatement" gorm:"column:legal_statement"`
}

type DbConfig struct {
	BaseEntity
	DbType   string `json:"dbType" gorm:"column:db_type"`
	Dbname   string `json:"dbname" gorm:"column:dbname"`
	Url      string `json:"url" gorm:"column:url"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

type FwConfig struct {
	BaseEntity
	Name      string `json:"name" gorm:"column:name"`
	Remark    string `json:"remark" gorm:"column:remark"`
	Content   string `json:"content" gorm:"column:content"`
	FlowType  string `json:"flowType" gorm:"column:flow_type"`
	FlowState string `json:"flowState" gorm:"column:flow_state"`
}

type Function struct {
	BaseEntity
	FuncName    string `json:"funcName" gorm:"column:func_name"`
	FuncContent string `json:"funcContent" gorm:"column:func_content"`
	RoleId      string `json:"roleId" gorm:"column:role_id"`
}

type FunctionImpl struct {
	BaseEntity
	Title       string `json:"title" gorm:"column:title"`
	Name        string `json:"name" gorm:"column:name"`
	Url         string `json:"url" gorm:"column:url"`
	Method      string `json:"method" gorm:"column:method"`
	ContentType string `json:"contentType" gorm:"column:content_type"`
	Headers     string `json:"headers" gorm:"column:headers"`
}

type Prompt struct {
	BaseEntity
	RoleName string `json:"roleName" gorm:"column:role_name"`
	Prompt   string `json:"prompt" gorm:"column:prompt"`
}

type Document struct {
	BaseEntity
	CollectionName string `json:"collectionName" gorm:"column:collection_name"`
	Filename       string `json:"filename" gorm:"column:filename"`
	Content        string `json:"content" gorm:"column:content"`
	Metadata       string `json:"metadata" gorm:"column:metadata"`
}

type AiChannel struct {
	BaseEntity
	Name        string `json:"name" gorm:"column:name"`
	Url         string `json:"url" gorm:"column:url"`
	Token       string `json:"token" gorm:"column:token"`
	MaxToken    int    `json:"maxToken" gorm:"column:max_token"`
	Models      string `json:"models" gorm:"column:models"`
	Priority    int    `json:"priority" gorm:"column:priority"`
	Remark      string `json:"remark" gorm:"column:remark"`
	OriginalUrl string `json:"originalUrl" gorm:"column:original_url"`
}
