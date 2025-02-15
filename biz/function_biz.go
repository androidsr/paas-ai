package biz

import (
	"encoding/json"
	"errors"
	"log"
	"paas-ai/entity"

	"github.com/androidsr/sc-go/mapper"
	"github.com/androidsr/sc-go/model"
	"github.com/androidsr/sc-go/sbuilder"
	"github.com/androidsr/sc-go/shttp"
	"github.com/androidsr/sc-go/sno"
	"github.com/tmc/langchaingo/llms"
)

type FunctionBiz struct {
	db *mapper.Mapper[entity.Function]
}

func NewFunctionBiz() *FunctionBiz {
	db := mapper.NewHelper[entity.Function]()
	return &FunctionBiz{db: db}
}

func (m *FunctionBiz) Get(id string) model.HttpResult {
	data := new(entity.Function)
	data.Id = id
	err := m.db.SelectOne(data)
	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return model.NewFailDefault("查询失败")
	}
	return model.NewOK(data)
}

func (m *FunctionBiz) Add(task *entity.Function) model.HttpResult {
	task.Id = sno.GetString()
	err := m.db.Insert(task)
	if err != nil {
		log.Printf("添加失败: %v\n", err)
		return model.NewFailDefault("添加失败")
	}
	return model.NewOK(task)
}

func (m *FunctionBiz) Edit(task *entity.Function) model.HttpResult {
	err := m.db.Update(task, "id", task.Id).Error
	if err != nil {
		log.Printf("修改失败: %v\n", err)
		return model.NewFailDefault("修改失败")
	}
	return model.NewOK(task)
}

func (t *FunctionBiz) Delete(id string) model.HttpResult {
	err := t.db.Delete("id", id).Error
	if err != nil {
		log.Printf("删除失败: %v\n", err)
		return model.NewFailDefault("删除失败")
	}
	return model.NewOK(true)
}

type FunctionQuery struct {
	Page     *model.PageInfo `json:"page" column:"-"`
	FuncName string          `json:"funcName" column:"a.func_name" keyword:"like"`
	RoleId   string          `json:"roleId" column:"a.role_id" keyword:"eq"`
}

func (m *FunctionBiz) Page(query *FunctionQuery) model.HttpResult {
	sql := `select * from function a where 1=1 `
	data := make([]entity.Function, 0)
	b := sbuilder.StructToBuilder(query, sql)
	sql, values := b.Build()
	sql += " order by func_name,id desc "
	result := m.db.SelectPage(&data, query.Page, sql, values...)
	return model.NewOK(result)
}

func (m *FunctionBiz) GetList() model.HttpResult {
	sql := `select id as value, func_name as label from function a where 1=1 `
	data := make([]model.SelectVO, 0)
	b := sbuilder.Builder(sql)
	sql, values := b.Build()
	sql += " order by func_name,id desc "
	err := m.db.SelectSQL(&data, sql, values...)
	if err != nil {
		log.Printf("查询列表失败: %v\n", err)
		return model.NewFailDefault("查询列表失败")
	}
	return model.NewOK(data)
}

func (m *FunctionBiz) ExecFuncCall(funcCall *llms.FunctionCall) (bool, string, error) {
	switch funcCall.Name {
	case "searchNetSource":
		//获取互联网内容
		//return searchNetSource(funcCall.Arguments)
	case "getTables":
		//查询表结构
		return getTables(funcCall.Arguments)
	case "getTableDetail":
		//查询表详情
		return getTableDetail(funcCall.Arguments)
	case "queryExecSql":
		return queryExecSql(funcCall.Arguments)
	default:
		impl := NewFunctionImplBiz().FindByName(funcCall.Name)
		if impl != nil {
			return true, "非找到对应的函数内容", nil
		}
		headers := make(map[string]string, 0)
		if impl.Headers == "" {
			json.Unmarshal([]byte(impl.Headers), &headers)
		}
		body := make(map[string]interface{}, 0)
		if funcCall.Arguments != "" {
			err := json.Unmarshal([]byte(funcCall.Arguments), &body)
			if err != nil {
				return true, "", errors.New("解析AI参数非JSON字符串")
			}
		}
		if impl.ContentType == "application/json" {

			switch impl.Method {
			case "GET":
				result, err := shttp.Get(impl.Url, impl.ContentType, headers)
				if err != nil {
					return true, "", errors.New("执行函数实现出错：" + err.Error())
				}
				return true, string(result), nil
			case "POST":
				result, err := shttp.Post(impl.Url, impl.ContentType, headers, &body)
				if err != nil {
					return true, "", errors.New("执行函数实现出错：" + err.Error())
				}
				return true, string(result), nil

			case "PUT":
				result, err := shttp.Put(impl.Url, impl.ContentType, headers, &body)
				if err != nil {
					return true, "", errors.New("执行函数实现出错：" + err.Error())
				}
				return true, string(result), nil

			case "DELETE":
				result, err := shttp.Delete(impl.Url, impl.ContentType, headers, &body)
				if err != nil {
					return true, "", errors.New("执行函数实现出错：" + err.Error())
				}
				return true, string(result), nil

			}
		} else {
			result, err := shttp.PostForm(impl.Url, impl.ContentType, headers, body)
			if err != nil {
				return true, "", errors.New("执行函数实现出错：" + err.Error())
			}
			return true, string(result), nil
		}
	}
	return false, "", errors.New("未找到对应的函数内容")
}

/** 获取数据库中所有表 **/
func getTables(arguments string) (bool, string, error) {
	input := make(map[string]string, 0)
	err := json.Unmarshal([]byte(arguments), &input)
	if err != nil {
		return true, "", errors.New("解析AI参数出错：" + err.Error())
	}
	dbname := input["dbname"]
	if dbname == "" {
		return true, "", errors.New("无效函数，未找到dbname参数")
	}

	info, err := NewDbConfigBiz().GetTables(dbname)
	if err != nil {
		log.Println("获取所有表信息出错")
		return true, "", errors.New("获取所有表信息出错")
	}
	bs, err := json.Marshal(info)
	if err != nil {
		return true, "", errors.New("表信息转换JSON出错：" + err.Error())
	}
	return false, string(bs), nil
}

/** 获取数据库表结构信息 **/
func getTableDetail(arguments string) (bool, string, error) {
	input := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(arguments), &input)
	if err != nil {
		return true, "", errors.New("解析AI参数出错：" + err.Error())
	}
	dbname, ok := input["dbname"].(string)
	if !ok || dbname == "" {
		return true, "", errors.New("无效函数，未找到dbname参数")
	}
	tableNames, ok := input["tableNames"].([]interface{})
	if !ok {
		log.Println("转换数据不对", ok)
		return true, "", errors.New("无效函数，未找到tableNames参数")
	}
	log.Println("执行函数：", dbname, tableNames)
	data, err := NewDbConfigBiz().GetTableDetail(dbname, tableNames)
	return false, data, err
}

/** 获取数据库表结构信息 **/
func queryExecSql(arguments string) (bool, string, error) {
	input := make(map[string]string, 0)
	err := json.Unmarshal([]byte(arguments), &input)
	if err != nil {
		return true, "", errors.New("解析AI参数出错：" + err.Error())
	}
	dbname := input["dbname"]
	if dbname == "" {
		return true, "", errors.New("无效函数，未找到dbname参数")
	}
	sql := input["sql"]
	if sql == "" {
		return true, "", errors.New("无效函数，未找到tableName参数")
	}
	data, err := NewDbConfigBiz().QueryExecSql(dbname, sql)
	return true, data, err
}

/** 在线搜索信息 **/
/* func searchNetSource(arguments string) (bool, string, error) {
	input := make(map[string]string, 0)
	err := json.Unmarshal([]byte(arguments), &input)
	if err != nil {
		return true, "", errors.New("解析AI参数出错：" + err.Error())
	}
	url := input["url"]
	if url == "" {
		return true, "", errors.New("无效函数，未找到url参数")
	}
	content := input["content"]
	if url == "" {
		return true, "", errors.New("无效函数，未找到url参数")
	}

	if content != "" {
		url += content
	}

	buf := bytes.Buffer{}
	web := NewChromedp()
	defer web.Close()

	cfg := NewCrawlerSelecterService().GetDomainInfo(url, "link")
	var selecter string
	if cfg == nil {
		selecter = "a"
	} else {
		selecter = cfg.LinkSelecter
	}
	links := web.GetLink(url, selecter)
	if links == nil {
		return false, "", nil
	}
	for _, v := range links {
		cfg := NewCrawlerSelecterService().GetDomainInfo(url, "content")
		if cfg == nil {
			selecter = "body"
		} else {
			selecter = cfg.TextSelecter
		}
		text, err := web.GetContent(v, selecter)
		if err != nil {
			log.Printf("获取文档内容处理失败：%s", err.Error())
			continue
		}
		tmp := bytes.NewBufferString(text)
		for {
			line, err := tmp.ReadString('\n')
			if err != nil {
				break
			}
			if line == "" || utf8.RuneCountInString(line) <= 8 {
				continue
			}
			buf.WriteString(line)
		}
	}
	return true, buf.String(), nil
}
*/
