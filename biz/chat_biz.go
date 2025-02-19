package biz

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"paas-ai/dv"
	"paas-ai/entity"
	"paas-ai/toolkit"
	"strings"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sc"
	"github.com/androidsr/sc-go/sno"
	"github.com/androidsr/sc-go/syaml"
	"github.com/tmc/langchaingo/llms"
)

type ChatBiz struct {
}

func NewChatBiz() *ChatBiz {
	return &ChatBiz{}
}

func LoadConfig() *dv.Config {
	db := mapper.NewHelper[entity.Config]()
	data := new(entity.Config)
	data.Id = "1"
	db.SelectOne(data)
	cfg = &dv.Config{}
	json.Unmarshal([]byte(data.Content), cfg)

	toolkit.Init(cfg)
	return cfg
}

func (m *ChatBiz) Clean() {
	messages = make([]llms.MessageContent, 0)
}

func (m *ChatBiz) Upload(files map[string]string) any {
	err := os.MkdirAll("temp/", os.ModePerm)
	if err != nil {
		return model.NewFailDefault(err.Error())
	}

	var fileNames []string
	for filename, dataStr := range files {
		newFileName := generateFileName(filename)
		dst, err := os.Create("temp/" + newFileName)
		if err != nil {
			return model.NewFailDefault(err.Error())
		}
		defer dst.Close()
		idx := strings.Index(dataStr, ",")
		data, err := base64.StdEncoding.DecodeString(dataStr[idx+1:])
		if err != nil {
			return model.NewFailDefault(err.Error())
		}
		dst.Write(data)
		dst.Close()
		fileNames = append(fileNames, newFileName)
	}
	fmt.Println(fileNames)
	return model.NewOK(strings.Join(fileNames, ","))
}

func (m *ChatBiz) GetImage(filename string) any {
	dir := "temp/"
	filePath := dir + "/" + filename

	// 检查文件是否存在
	if !sc.IsFile(filePath) {
		return model.NewFailDefault("文件不存在")
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return model.NewFailDefault("读取文件失败")
	}
	// 转换为 Base64 编码
	base64Data := base64.StdEncoding.EncodeToString(data)
	return model.NewOK(base64Data)
}

func init() {
	mapper.Initdb(&syaml.GormInfo{
		Driver:  "sqlite",
		Url:     "sqlite.db",
		MaxOpen: 10,
		MaxIdle: 2,
		ShowSql: true,
	})
	sno.New(syaml.SnowflakeInfo{WorkerId: 1})

	db := mapper.NewHelper[entity.Config]()
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "config" (
			"id" integer NOT NULL,
			"content" TEXT NOT NULL,
			"legal_statement" TEXT NOT NULL,
			PRIMARY KEY ("id")
		);
		
		CREATE TABLE IF NOT EXISTS "db_config" (
			"id" integer NOT NULL,
			"db_type" TEXT DEFAULT '',
			"dbname" TEXT DEFAULT '',
			"url" TEXT DEFAULT '',
			"username" TEXT DEFAULT '',
			"password" TEXT DEFAULT '',
			PRIMARY KEY ("id")
		);

		CREATE TABLE IF NOT EXISTS "fw_config" (
			"id" integer NOT NULL,
			"name" TEXT DEFAULT '',
			"remark" TEXT DEFAULT '',
			"content" TEXT DEFAULT '',
			"flow_type" TEXT DEFAULT '',
			"flow_state" TEXT DEFAULT '',
			PRIMARY KEY ("id")
		);

		CREATE TABLE IF NOT EXISTS "function" (
			"id" integer NOT NULL,
			"func_name" TEXT DEFAULT '',
			"func_content" TEXT DEFAULT '',
			"role_id" TEXT DEFAULT '',
			PRIMARY KEY ("id")
		);

		CREATE TABLE IF NOT EXISTS "function_impl" (
			"id" integer NOT NULL,
			"title" TEXT DEFAULT '',
			"name" TEXT DEFAULT '',
			"url" TEXT DEFAULT '',
			"content_type" TEXT DEFAULT '',
			"method" TEXT DEFAULT '',
			"headers" TEXT DEFAULT '',
			PRIMARY KEY ("id")
		);

		CREATE TABLE IF NOT EXISTS "prompt" (
			"id" integer NOT NULL,
			"role_name" TEXT DEFAULT '',
			"prompt" TEXT DEFAULT '',
			PRIMARY KEY ("id")
		);
		
		CREATE TABLE IF NOT EXISTS "document" (
			"id" integer NOT NULL,
			"collection_name" TEXT DEFAULT '',
			"filename" TEXT DEFAULT '',
			"content" TEXT DEFAULT '',
			"metadata" TEXT DEFAULT '',
			PRIMARY KEY ("id")
		);

		CREATE TABLE IF NOT EXISTS "ai_channel" (
			"id" integer NOT NULL,
			"name" TEXT DEFAULT '',
			"url" TEXT DEFAULT '',
			"token" TEXT DEFAULT '',
			"max_token" integer NOT NULL,
			"models" TEXT DEFAULT '',
			"remark" TEXT DEFAULT '',
			"priority" integer NOT NULL,
			"original_url" TEXT NOT NULL,
			PRIMARY KEY ("id")
		);
	`).Error
	if err != nil {
		panic(err)
	}
	db.Exec(`
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889136494312951808, '本机ollama', 'http://localhost:11434/v1', 'token', 11, 'deepseek-r1:8b,qwen2.5-coder:7b,ZimaBlueAI/MiniCPM-o-2_6:latest', -1, NULL, NULL);
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889158531488157696, '文心一言', 'https://qianfan.baidubce.com/v2', '', 99, 'ernie-4.0-8k-latest,ernie-4.0-turbo-128k', -1, '大模型token申请地址：
		https://console.bce.baidu.com/iam/#/iam/apikey/list', NULL);
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889183926614757376, '通义千问', 'https://dashscope.aliyuncs.com/compatible-mode/v1', '', 3, 'qwen-max-latest', -1, '大模型token申请地址：
		https://bailian.console.aliyun.com/?apiKey=1#/api-key', 'https://tongyi.aliyun.com/qianwen');
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889184913454796800, '智谱清言', 'https://open.bigmodel.cn/api/paas/v4', '', 99, 'glm-4-plus,glm-4-air,glm-4-airx,glm-4-long,glm-4-flashx,glm-4-flash', -1, '大模型token申请地址：
		https://www.bigmodel.cn/usercenter/proj-mgmt/apikeys', 'https://chatglm.cn/main/alltoolsdetail?lang=zh');
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889185292728930304, '科大讯飞', 'https://spark-api-open.xf-yun.com/v1', '', 6, 'generalv3.5', -1, '大模型token申请地址：
		https://console.xfyun.cn/app/myapp', 'https://xinghuo.xfyun.cn/desk');
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889185484219879424, '百川大模型', 'https://api.baichuan-ai.com/v1', '', 2, 'Baichuan4-Turbo,Baichuan4-Air,Baichuan4,Baichuan3-Turbo,Baichuan3-Turbo-128k,Baichuan2-Turbo', -1, '模型token申请地址：
		https://platform.baichuan-ai.com/console/apikey
		', 'https://ying.baichuan-ai.com/chat?from=home');
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889186198371438592, '腾讯混元', 'https://api.hunyuan.cloud.tencent.com/v1', '', 5, 'hunyuan-turbo-latest', -1, '大模型token申请地址：
		https://console.cloud.tencent.com/hunyuan/api-key', 'https://yuanbao.tencent.com/chat/naQivTmsDa');
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889221866447441920, 'DeepSeek', 'https://api.deepseek.com', '', 99, 'deepseek-chat,deepseek-reasoner', -1, '大模型token 申请地址：
		https://platform.deepseek.com/api_keys', 'https://chat.deepseek.com/');
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1889297593876353024, '硅基流动', 'https://api.siliconflow.cn/v1', '', 6, 'deepseek-ai/DeepSeek-R1-Distill-Llama-70B,deepseek-ai/DeepSeek-R1,deepseek-ai/DeepSeek-R1-Distill-Qwen-32B,deepseek-ai/DeepSeek-R1-Distill-Qwen-14B,deepseek-ai/DeepSeek-R1-Distill-Llama-8B,deepseek-ai/DeepSeek-R1-Distill-Qwen-7B,deepseek-ai/DeepSeek-R1-Distill-Qwen-1.5B,Pro/deepseek-ai/DeepSeek-R1-Distill-Llama-8B,Pro/deepseek-ai/DeepSeek-R1-Distill-Qwen-7B,Pro/deepseek-ai/DeepSeek-R1-Distill-Qwen-1.5B', -1, '大模型token申请地址：
		https://cloud.siliconflow.cn/account/ak', 'https://cloud.siliconflow.cn/models');
		INSERT INTO "ai_channel" ("id", "name", "url", "token", "priority", "models", "max_token", "remark", "original_url") VALUES (1891785744574320640, 'ChatGPT', 'https://chat.openai.com/chat', '1', 1, '1', -1, '', 'https://chat.openai.com/chat');
	`)
	LoadConfig()
	fmt.Println("初始化数据库成功！")

}
