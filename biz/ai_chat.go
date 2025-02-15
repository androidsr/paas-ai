package biz

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"paas-ai/access"
	"paas-ai/dv"
	"paas-ai/entity"
	"paas-ai/toolkit"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sc"
	"github.com/androidsr/sc-go/sgin"
	"github.com/androidsr/sc-go/sno"
	"github.com/androidsr/sc-go/syaml"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
)

var (
	ctx      = context.Background()
	lock     sync.Mutex
	chatLock sync.Mutex
	messages = make([]llms.MessageContent, 0)
	cfg      *dv.Config
	channel  *entity.AiChannel
)

func init() {
	_, err := os.Stat("./sc-go-router")
	if err != nil {
		os.WriteFile("./sc-go-router", []byte(`{"AiChannelBiz":{},"AiChannelQuery":{},"AiChat":{"CleanLogs":"@Router [get] [/chat/cleanLogs]\n","Download":"@Router [get] [/chat/download/:filename]\n","File":"@Router [post] [/chat/file]\n","Flow":"@Router [post] [/chat/flow]\n","Image":"@Router [post] [/chat/image]\n","Text":"@Router [post] [/chat/text]\n","Tools":"@Router [post] [/chat/tools]\n","Vector":"@Router [post] [/chat/vector]\n"},"BaseDTO":{},"ChatBiz":{},"CollectionBiz":{},"DbConfigBiz":{},"DbDataSourceQuery":{},"DocumentBiz":{},"DocumentQuery":{},"FileDTO":{},"FlowDTO":{},"FunctionBiz":{},"FunctionQuery":{},"FwConfigBiz":{},"FwConfigQuery":{},"ImageDTO":{},"PromptBiz":{},"PromptQuery":{},"TextDTO":{},"ToolsDTO":{},"VectorDTO":{}}`), 0644)
	}
	router := sgin.New(&syaml.GinInfo{Port: 11083, Scan: &syaml.GinScanInfo{Pkg: "biz", Filter: "@Router"}})
	router.Cors()
	go func() {
		router.RunServer()
	}()
	sgin.AddRouter(AiChat{prompt: NewPromptBiz(), function: NewFunctionBiz()})
	fmt.Println("初始化web层成功！")
}

func AddMessage(llmType llms.ChatMessageType, message string, cacheLimit int) []llms.MessageContent {
	lock.Lock()
	defer lock.Unlock()
	if messages == nil {
		messages = make([]llms.MessageContent, 0)
	}
	if message != "" {
		text := llms.TextParts(llmType, message)
		messages = append(messages, text)
	}
	if len(messages) > cacheLimit*2 {
		messages = messages[1:]
	}
	return messages
}

type AiChat struct {
	prompt   *PromptBiz
	function *FunctionBiz
}

// @Router [get] [/chat/download/:filename]
func (m AiChat) Download(c *gin.Context) {
	filename := c.Param("filename")
	dir := "temp/"
	if !sc.IsFile(dir + "/" + filename) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "文件不存在"})
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(dir + "/" + filename)
}

// 生成新的文件名
func generateFileName(originalName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(originalName)
	// 生成新的文件名
	newFileName := fmt.Sprintf("%s%s", sno.GetString(), ext)
	return newFileName
}

// @Router [get] [/chat/cleanLogs]
func (m AiChat) CleanLogs(c *gin.Context) any {
	f, _ := os.OpenFile("system.log", os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	f.Seek(0, 0)
	f.WriteAt([]byte{}, 0)
	f.Seek(0, 0)
	return model.NewOK("清理成功")
}

type BaseDTO struct {
	ChannelId   string `json:"channelId"`   //渠道ID
	Stream      bool   `json:"stream"`      //是否流式输出
	SystemId    string `json:"systemId"`    //系统提示词id
	System      string `json:"system"`      //输入系统提示词
	Message     string `json:"message"`     //消息体
	CacheLimit  int    `json:"cacheLimit"`  //会话历史大小
	Model       string `json:"model"`       //模型名称
	Temperature string `json:"temperature"` //温度为采样温度，取值范围在0到1之间。
	TopK        string `json:"TopK"`        // TopK是top-k采样时要考虑的令牌数量。
	TopP        string `json:"TopP"`        // TopP是top-p采样的累积概率。
}

type TextDTO struct {
	BaseDTO
	JsonFormat bool `json:"jsonFormat"` // json格式化输出
}

type VectorDTO struct {
	BaseDTO
	CollectionName  string         `json:"collectionName"`  //向量库集合名称
	SimilarityScore string         `json:"similarityScore"` //向量库相似度
	Filter          map[string]any `json:"filter"`          //向量库过虑条件
	PageSize        int            `json:"pageSize"`        //向最库最大条数
}

type ToolsDTO struct {
	BaseDTO
	FuncCall   string `json:"funcCall"`   // 指定函数
	DatabaseId string `json:"databaseId"` //数据库id
}

type FlowDTO struct {
	BaseDTO
	FlowId string `json:"flowId"` //流程图对话
}

type FileDTO struct {
	BaseDTO
	UploadFileId string `json:"uploadFileId"` //文件ID
}

type ImageDTO struct {
	BaseDTO
	UploadFileId string `json:"uploadFileId"` //文件ID
}

func settings(c *gin.Context, input BaseDTO) []llms.CallOption {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("panic recovered:\n%s\n%s", err, debug.Stack())
			sendMessage(c.Writer, []byte(err.(error).Error()))
		}
	}()
	options := make([]llms.CallOption, 0)
	topKstr := input.TopK
	topK, _ := strconv.ParseInt(topKstr, 10, 64)
	if topK != 0 {
		options = append(options, llms.WithTopK(int(topK)))
	}
	topPstr := input.TopP
	topP, _ := strconv.ParseFloat(topPstr, 64)
	if topP != 0 {
		options = append(options, llms.WithTopP(topP))
	}

	temperatureStr := input.Temperature
	temperature, _ := strconv.ParseFloat(temperatureStr, 64)

	if temperature != 0 {
		options = append(options, llms.WithTemperature(temperature))
	}
	channelDb := mapper.NewHelper[entity.AiChannel]()
	channel = new(entity.AiChannel)
	channel.Id = input.ChannelId
	channelDb.SelectOne(channel)
	if channel.MaxToken != -1 {
		options = append(options, llms.WithMaxTokens(int(channel.MaxToken)))
	}
	return options
}

func (m AiChat) loadSystem(input BaseDTO, chat *toolkit.OpenAI) {
	var system string
	if input.SystemId != "" {
		system = m.GetSystemInfo(input.SystemId)
		if system != "" {
			chat.SetContent(llms.TextParts(llms.ChatMessageTypeSystem, system))
		}
	}

	if input.System != "" {
		chat.SetContent(llms.TextParts(llms.ChatMessageTypeSystem, input.System))
	}
}

// @Router [post] [/chat/text]
func (m AiChat) Text(c *gin.Context) {
	chatLock.Lock()
	defer chatLock.Unlock()
	input := new(TextDTO)
	c.BindJSON(input)
	options := settings(c, input.BaseDTO)
	if input.JsonFormat {
		options = append(options, llms.WithJSONMode())
	}
	message := input.Message
	urlRegex := regexp.MustCompile(`\bhttps?://\S+\b`)
	urls := urlRegex.FindAllString(message, -1)
	if len(urls) != 0 {
		browser, err := toolkit.NewBrowser(cfg)
		if err == nil {
			for _, url := range urls {
				browser.NewPage(url)
				data, err := browser.GetLinkContent("body", "text")
				context := bytes.Buffer{}
				for _, v := range data {
					context.WriteString(v.(string))
				}
				if err != nil {
					message = strings.Replace(message, url, "", 1)
				} else {
					message = strings.Replace(message, url, context.String(), 1)
				}
			}
			browser.Close()
		}
	}
	chat, err := toolkit.NewOpenAI(channel.Url, channel.Token, input.Model)
	if err != nil {
		sendMessage(c.Writer, []byte("模型连接失败"))
		return
	}
	m.loadSystem(input.BaseDTO, chat)
	chat.SetContent(AddMessage(llms.ChatMessageTypeHuman, message, input.CacheLimit)...)
	result := m.output(c, chat, input.Stream, options)
	AddMessage(llms.ChatMessageTypeAI, result.String(), input.CacheLimit)
}

// @Router [post] [/chat/vector]
func (m AiChat) Vector(c *gin.Context) {
	chatLock.Lock()
	defer chatLock.Unlock()
	input := new(VectorDTO)
	c.BindJSON(input)
	options := settings(c, input.BaseDTO)
	message := input.Message

	chat, err := toolkit.NewOpenAI(channel.Url, channel.Token, input.Model)
	if err != nil {
		sendMessage(c.Writer, []byte("模型连接失败"))
		return
	}
	m.loadSystem(input.BaseDTO, chat)

	score, err := strconv.ParseFloat(input.SimilarityScore, 32)
	if err != nil {
		score = 0
	}
	filter := make(map[string]string, 0)
	for k, v := range input.Filter {
		if v != "" && v != nil {
			filter[k] = fmt.Sprint(v)
		}
	}
	vc := toolkit.NewCollection().Get(input.CollectionName)
	vdoc := toolkit.NewDocument(vc)
	docs, err := vdoc.Query(message, input.PageSize, filter)
	if err != nil || len(docs) == 0 {
		log.Println("获取到向量库数据失败：", err, filter)
		sendMessage(c.Writer, []byte("向量库未匹配到内容"))
		return
	}

	if len(docs) == 0 {
		sendMessage(c.Writer, []byte("##### 本次问答未找到任何本地知识库内容 \n  "))
	} /* else {
		log.Println("获取到本地数据：", len(docs))
		sendMessage(c.Writer, []byte(fmt.Sprintf("#### 检索到本地知识库数量：%d  \n  ", len(docs))))
	} */

	for _, v := range docs {
		if float64(v.Similarity) < score {
			continue
		}
		content := v.Content
		txt := []rune(content)
		if len(txt) > 20 {
			txt = txt[:20]
		}
		docConfig := mapper.NewHelper[entity.Document]()
		data := &entity.Document{}
		data.Id = v.ID
		err := docConfig.SelectOne(data)
		if err != nil {
			continue
		}
		sendMessage(c.Writer, []byte(fmt.Sprintf("##### 来源：[%s]、相似程度：%.2f、内容概要：%s  \n  ", data.Filename, v.Similarity, strings.TrimSpace(string(txt)))))
		chat.SetContent(llms.TextParts(llms.ChatMessageTypeHuman, content))
	}

	chat.SetContent(AddMessage(llms.ChatMessageTypeHuman, message, input.CacheLimit)...)
	result := m.output(c, chat, input.Stream, options)
	AddMessage(llms.ChatMessageTypeAI, result.String(), input.CacheLimit)

}

// @Router [post] [/chat/tools]
func (m AiChat) Tools(c *gin.Context) {
	chatLock.Lock()
	defer chatLock.Unlock()
	input := new(ToolsDTO)
	c.BindJSON(input)
	options := settings(c, input.BaseDTO)
	var toolsContent string
	if input.FuncCall != "" {
		funcInfo := &entity.Function{}
		funcInfo.Id = input.FuncCall
		mapper.NewHelper[entity.Function]().SelectOne(funcInfo)
		if funcInfo != nil {
			toolsContent = funcInfo.FuncContent
		}
	}

	chat, err := toolkit.NewOpenAI(channel.Url, channel.Token, input.Model)
	if err != nil {
		sendMessage(c.Writer, []byte("模型连接失败"))
		return
	}
	m.loadSystem(input.BaseDTO, chat)
	chat.SetContent(AddMessage(llms.ChatMessageTypeHuman, input.Message, input.CacheLimit)...)

	if !strings.HasPrefix(strings.TrimSpace(toolsContent), "[") {
		toolsContent = "[" + toolsContent + "]"
	}
	result := bytes.Buffer{}
	funcResult, err := chat.GenerateFunction(ctx, toolsContent, func(fun *llms.FunctionCall) (bool, string, error) {
		return NewFunctionBiz().ExecFuncCall(fun)
	}, func(ctx context.Context, chunk []byte) error {
		result.Write(chunk)
		sendMessage(c.Writer, chunk)
		return nil
	}, options...)

	if err != nil {
		sendMessage(c.Writer, []byte(err.Error()))
		return
	}

	AddMessage(llms.ChatMessageTypeAI, result.String(), input.CacheLimit)
	fmt.Println(funcResult)
	sendMessage(c.Writer, []byte(funcResult))
}

// @Router [post] [/chat/flow]
func (m AiChat) Flow(c *gin.Context) {
	chatLock.Lock()
	defer chatLock.Unlock()
	input := new(FlowDTO)
	c.BindJSON(input)

	startData := make(map[string]any, 0)
	startData["message"] = input.Message
	/* service.NewFlowActuator(input.FlowId, func(data string, err error) {
		if err != nil {
			sendMessage(c.Writer, []byte(err.Error()))
			return
		}
		sendMessage(c.Writer, []byte(data))
	}).StartFlow(startData) */
}

// @Router [post] [/chat/file]
func (m AiChat) File(c *gin.Context) {
	chatLock.Lock()
	defer chatLock.Unlock()
	input := new(FileDTO)
	c.BindJSON(input)
	options := settings(c, input.BaseDTO)

	sendData := bytes.Buffer{}
	data, err := access.AutoLoad("temp/" + input.UploadFileId)
	if err != nil {
		sendMessage(c.Writer, []byte(err.Error()))
		return
	}
	for _, v := range data {
		sendData.WriteString(v.PageContent)
	}
	chat, err := toolkit.NewOpenAI(channel.Url, channel.Token, input.Model)
	if err != nil {
		sendMessage(c.Writer, []byte("模型连接失败"))
		return
	}
	m.loadSystem(input.BaseDTO, chat)
	chat.SetContent(llms.TextParts(llms.ChatMessageTypeHuman, sendData.String()+";"+input.Message))
	result := m.output(c, chat, input.Stream, options)
	AddMessage(llms.ChatMessageTypeAI, result.String(), input.CacheLimit)
}

// @Router [post] [/chat/image]
func (m AiChat) Image(c *gin.Context) {
	chatLock.Lock()
	defer chatLock.Unlock()
	input := new(ImageDTO)
	c.BindJSON(input)
	options := settings(c, input.BaseDTO)

	message := input.Message
	fmt.Println(input.Model)
	chat, err := toolkit.NewOpenAI(channel.Url, channel.Token, input.Model)
	if err != nil {
		sendMessage(c.Writer, []byte("模型连接失败"))
		return
	}
	m.loadSystem(input.BaseDTO, chat)
	ids := strings.Split(input.UploadFileId, ",")
	parts := make([]llms.ContentPart, 0)
	for _, id := range ids {
		data, err := compressImageFromFile("temp/"+id, 300)
		if err != nil {
			return
		}
		base64Data := base64.StdEncoding.EncodeToString(data)
		parts = append(parts, llms.ImageURLPart("data:image/png;base64,"+base64Data))
	}

	chat.SetContent(llms.MessageContent{Role: llms.ChatMessageTypeHuman, Parts: parts})

	chat.SetContent(AddMessage(llms.ChatMessageTypeHuman, message, input.CacheLimit)...)
	result := m.output(c, chat, input.Stream, options)
	AddMessage(llms.ChatMessageTypeAI, result.String(), input.CacheLimit)
}

// compressImageFromFile 压缩图片文件并返回字节流
func compressImageFromFile(inputPath string, maxSizeKB int) ([]byte, error) {
	// 打开文件
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 解码图片
	img, err := imaging.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 初始化宽高和质量参数
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	quality := 90 // 初始质量
	var buf bytes.Buffer

	// 压缩处理循环
	for {
		// 等比例缩小图片（每次缩小 10%）
		img = imaging.Resize(img, width*9/10, height*9/10, imaging.Lanczos)

		// 重置缓冲区
		buf.Reset()

		// 编码为 JPEG
		err := imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(quality))
		if err != nil {
			return nil, fmt.Errorf("failed to encode image: %w", err)
		}

		// 检查是否符合目标大小
		if buf.Len() <= maxSizeKB*1024 {
			break
		}

		// 动态调整参数
		width = width * 9 / 10
		height = height * 9 / 10
		quality -= 5 // 减小质量
		if quality < 30 {
			return nil, fmt.Errorf("failed to compress image within size limits")
		}
	}
	return buf.Bytes(), nil
}

func (m AiChat) output(c *gin.Context, chat *toolkit.OpenAI, stream bool, options []llms.CallOption) *bytes.Buffer {
	result := new(bytes.Buffer)
	if stream {
		_, err := chat.GenerateCallback(ctx, func(ctx context.Context, chunk []byte) error {
			result.Write(chunk)
			sendMessage(c.Writer, chunk)
			return nil
		}, options...)
		if err != nil {
			sendMessage(c.Writer, []byte(err.Error()))
			return nil
		}
	} else {
		resp, err := chat.GenerateCallback(ctx, nil, options...)
		if err != nil {
			sendMessage(c.Writer, []byte(err.Error()))
			return nil
		}
		choice1 := resp.Choices[0]
		sendMessage(c.Writer, []byte(choice1.Content))
		result.WriteString(choice1.Content)
	}
	return result
}

func (m AiChat) GetSystemInfo(roleId string) string {
	data := new(entity.Prompt)
	data.Id = roleId
	mapper.NewHelper[entity.Prompt]().SelectOne(data)
	return data.Prompt
}

func sendMessage(writer gin.ResponseWriter, message []byte) {
	defer func() {
		recover()
	}()
	writer.WriteString(fmt.Sprintf("data: %s\n\n", base64.StdEncoding.EncodeToString(message)))
	writer.Flush()
}
