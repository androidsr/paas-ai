package dv

type Config struct {
	ChannelId      string  `json:"channelId"`
	EmbeddingUrl   string  `json:"embeddingUrl"`
	EmbeddingToken string  `json:"embeddingToken"`
	EmbeddingModel string  `json:"embeddingModel"`
	Browser        string  `json:"browser"`
	ShowBrowser    bool    `json:"showBrowser"`
	Timeout        float64 `json:"timeout"`
	PageLimit      float64 `json:"pageLimit"`
	PureType       bool    `json:"pureType"`
	StartInterval  float64 `json:"startInterval"`
	EndInterval    float64 `json:"endInterval"`
	LinkText       float64 `json:"linkText"`
}

type TableInfo struct {
	TableName string `json:"tableName"`
	TableDesc string `json:"tableDesc"`
}
type TableDetail struct {
	TableName string              `json:"tableName"`
	Columns   []map[string]string `json:"columns"`
}
