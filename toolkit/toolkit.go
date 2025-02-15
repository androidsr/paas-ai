package toolkit

import (
	"bytes"
	"fmt"
	"log"
	"paas-ai/dv"
	"strings"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
	"github.com/playwright-community/playwright-go"
)

// PageToolkit 封装的网页操作工具
type PageToolkit struct {
	Browser         playwright.Browser
	Context         playwright.BrowserContext
	Page            playwright.Page
	isRecording     bool
	PageLimit       float64
	timeout         *float64
	PureType        bool
	recordedActions []string
	links           map[string]string
	linkText        float64
}

// NewBrowser 创建新的浏览器实例
func NewBrowser(conf *dv.Config) (*PageToolkit, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("无法启动 Playwright: %v", err)
	}
	browserPath := conf.Browser
	showBrowser := conf.ShowBrowser
	pageLimit := conf.PageLimit
	if pageLimit == 0 {
		pageLimit = 20
	}
	timeout := conf.Timeout
	if timeout == 0 {
		timeout = 5000
	}
	pureType := conf.PureType

	linkText := conf.LinkText
	if linkText == 0 {
		linkText = 6
	}
	var interval float64
	startInterval := conf.StartInterval
	if startInterval == 0 {
		interval = 650
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		ExecutablePath: playwright.String(browserPath),
		Headless:       playwright.Bool(!showBrowser),
		Timeout:        playwright.Float(timeout),  // 设置浏览器启动超时为 30 秒
		SlowMo:         playwright.Float(interval), // 每个操作延迟 50 毫秒
	})
	if err != nil {
		return nil, fmt.Errorf("无法启动 Chromium 浏览器: %v", err)
	}

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		AcceptDownloads: playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("无法创建浏览器上下文: %v", err)
	}
	return &PageToolkit{
		Browser:   browser,
		Context:   context,
		PageLimit: pageLimit,
		timeout:   playwright.Float(timeout),
		PureType:  pureType,
		linkText:  linkText,
	}, nil
}

// NewPage 打开新的页面
func (p *PageToolkit) NewPage(url string) (*PageToolkit, error) {
	var err error

	if p.PureType {
		p.Page, err = p.Context.NewPage()
	} else {
		p.Page, err = p.Browser.NewPage()
	}

	if err != nil {
		return nil, fmt.Errorf("无法创建新页面: %v", err)
	}

	if _, err := p.Page.Goto(url); err != nil {
		return nil, fmt.Errorf("无法跳转到指定 URL: %v", err)
	}
	return p, nil
}

// Close 关闭浏览器
func (p *PageToolkit) Close() error {
	if err := p.Browser.Close(); err != nil {
		return fmt.Errorf("无法关闭浏览器: %v", err)
	}
	return nil
}

// GetCookie 获取Cookie信息
func (p *PageToolkit) GetCookie(url string) (map[string]string, error) {
	err := p.Page.WaitForLoadState()
	if err != nil {
		return nil, fmt.Errorf("无法获取 Cookie: %v", err)
	}
	cs, err := p.Page.Context().Cookies(url)
	if err != nil {
		return nil, fmt.Errorf("无法获取 Cookie: %v", err)
	}
	result := make(map[string]string, 0)
	for _, c := range cs {
		result[c.Name] = c.Value
	}
	return result, nil
}

// GetAttribute 获取元素的指定属性值
func (p *PageToolkit) GetAttribute(selector, attribute string) (string, error) {
	locator := p.Page.Locator(selector).First()
	if count, _ := locator.Count(); count == 0 {
		return "", fmt.Errorf("未找到选择器 %v 对应的元素", selector)
	}
	attr, err := locator.GetAttribute(attribute)
	if err != nil {
		return "", fmt.Errorf("无法获取选择器 %v 的属性 %v: %v", selector, attribute, err)
	}
	return attr, nil
}

// GetAllAttributes 获取所有匹配元素的指定属性值
func (p *PageToolkit) GetAllAttributes(selector, attribute string) ([]string, error) {
	locators := p.Page.Locator(selector).First()
	count, err := locators.Count()
	if err != nil || count == 0 {
		return nil, fmt.Errorf("未找到选择器 %v 对应的任何元素", selector)
	}
	var attributes []string
	for i := 0; i < count; i++ {
		attr, _ := locators.Nth(i).GetAttribute(attribute)
		if attr != "" {
			attributes = append(attributes, attr)
		}
	}
	return attributes, nil
}

// SetInputValue 设置输入框的值
func (p *PageToolkit) SetInputValue(selector, value string) error {
	locator := p.Page.Locator(selector).First()
	if count, _ := locator.Count(); count == 0 {
		return fmt.Errorf("未找到选择器 %v 对应的输入框", selector)
	}
	if err := locator.Fill(value); err != nil {
		return fmt.Errorf("无法填充输入框: %v", err)
	}
	return nil
}

// Click 点击指定元素
func (p *PageToolkit) Click(selector string) error {
	locator := p.Page.Locator(selector).First()
	if count, _ := locator.Count(); count == 0 {
		return fmt.Errorf("未找到选择器 %v 对应的元素", selector)
	}
	if err := locator.Click(); err != nil {
		return fmt.Errorf("无法点击元素: %v", err)
	}
	return nil
}

// Keyboard 执行键盘操作
func (p *PageToolkit) Keyboard(key string) error {
	err := p.Page.Keyboard().Press("Enter")
	if err != nil {
		return fmt.Errorf("无法按键 %v: %v", key, err)
	}
	return nil
}

func (p *PageToolkit) ClickForWait(selector string) error {
	if err := p.WaitForElement(selector); err != nil {
		log.Printf("等待页面加载异常：%v", err)
		return err
	}
	return p.Click(selector)
}

func (p *PageToolkit) GetCurrentPage() (playwright.Page, error) {
	// 获取当前 context 下的所有页面
	pages := p.Context.Pages()
	// 假设最后一个页面是当前活动页面
	if len(pages) > 0 {
		return pages[len(pages)-1], nil
	}

	return nil, fmt.Errorf("没有找到当前活动页面")
}

// Screenshot 截取屏幕截图
func (p *PageToolkit) Screenshot(path string) error {
	if _, err := p.Page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(path),
	}); err != nil {
		return fmt.Errorf("无法截取屏幕截图: %v", err)
	}
	return nil
}

// DownloadFile 下载文件
func (p *PageToolkit) DownloadFile(selector, downloadPath string) error {
	locator := p.Page.Locator(selector).First()
	if count, _ := locator.Count(); count == 0 {
		return fmt.Errorf("未找到选择器 %v 对应的元素", selector)
	}
	url, _ := locator.GetAttribute("href")
	if url != "" {
		fmt.Printf("文件下载地址: %s\n", url)
		// 此处可以扩展为实际文件下载逻辑
	}
	return nil
}

// WaitForElement 等待元素出现
func (p *PageToolkit) WaitForElement(selector string) error {
	if err := p.Page.Locator(selector).First().WaitFor(playwright.LocatorWaitForOptions{Timeout: p.timeout}); err != nil {
		return fmt.Errorf("等待选择器 %v 对应的元素时发生错误: %v", selector, err)
	}
	return nil
}

// ScrollToBottom 滚动到页面底部
func (p *PageToolkit) ScrollToBottom() error {
	// 确保页面完全加载
	if _, err := p.Page.Evaluate(`document.readyState === "complete"`); err != nil {
		return fmt.Errorf("无法确认页面加载完成: %v", err)
	}

	// 使用 window.scrollTo 滚动到页面底部
	if _, err := p.Page.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`); err != nil {
		return fmt.Errorf("无法滚动到页面底部: %v", err)
	}

	return nil
}

// GetTextContent 获取指定选择器的文本内容
func (p *PageToolkit) GetTextContent(selector string) (string, error) {
	locator := p.Page.Locator(selector).First()

	if count, _ := locator.Count(); count == 0 {
		return "", fmt.Errorf("未找到选择器 %v 对应的元素", selector)
	}
	return p.htmlToPdf(locator), nil
}

func (p *PageToolkit) GetLinks(selector string) (map[string]string, error) {
	// 等待元素加载完成
	p.WaitForElement(selector)

	// 获取父元素的 Locator
	parentLocator := p.Page.Locator(selector).First()
	links, err := p.getRecursiveElementsSources(parentLocator, "a") // 递归获取链接
	if err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return nil, fmt.Errorf("未找到有效的链接")
	}
	if len(links) > int(p.PageLimit) {
		links = getLongestKeysMap(links, int(p.PageLimit))
	}
	p.links = links
	return links, nil
}

func getLongestKeysMap(m map[string]string, count int) map[string]string {
	// 创建一个切片来存储按长度排序的 keys
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	// 按 key 的长度降序排序
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if len(keys[i]) < len(keys[j]) {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	// 创建一个新的 map 以存储指定数量的最长 key
	newMap := make(map[string]string)
	for i := 0; i < count && i < len(keys); i++ {
		newMap[keys[i]] = m[keys[i]]
	}

	return newMap
}

// getRecursiveElementsSources 递归获取指定类型元素的所有源 URL
func (p *PageToolkit) getRecursiveElementsSources(locator playwright.Locator, elementType string) (map[string]string, error) {
	// 获取所有匹配的子元素定位器
	locatorAll, err := locator.Locator(elementType).All()
	if err != nil {
		return nil, err
	}

	// 创建一个用于存储文本和链接的 map
	sources := make(map[string]string)

	// 用于过滤的函数
	filterText := func(text string) bool {
		return utf8.RuneCountInString(text) >= int(p.linkText) && !strings.ContainsAny(text, "\n\r\t")
	}

	// 遍历所有匹配的子元素，查找指定类型的元素
	for i := 0; i < len(locatorAll); i++ {
		locatorItem := locatorAll[i]

		// 获取当前元素的文本内容
		text, err := locatorItem.InnerText()
		if err != nil {
			continue
		}
		name := strings.TrimSpace(text)

		// 如果文本不符合过滤条件，跳过
		if !filterText(name) {
			continue
		}

		// 获取当前元素的源属性
		var src string
		switch elementType {
		case "a":
			src, err = locatorItem.GetAttribute("href")
		case "video", "audio":
			src, err = locatorItem.GetAttribute("src")
		case "img":
			src, err = locatorItem.GetAttribute("src")
		default:
			return nil, fmt.Errorf("不支持的元素类型: %v", elementType)
		}

		if err != nil {
			return nil, fmt.Errorf("获取第 %d 个 %s 源时出错: %v", i, elementType, err)
		}

		// 如果获取到的源地址和文本有效，则加入 map
		if src != "" && name != "" {
			sources[name] = src
		}

		// 递归查找当前元素下的所有匹配的子元素
		subLocator := locatorItem.Locator(elementType)
		subSources, err := p.getRecursiveElementsSources(subLocator, elementType)
		if err != nil {
			return nil, fmt.Errorf("递归获取子元素时出错: %v", err)
		}

		// 合并子元素的结果到主 map 中
		for k, v := range subSources {
			sources[k] = v
		}
	}

	return sources, nil
}

// 获取链接下的内容
func (p *PageToolkit) GetLinkContent(selector string, elementType string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	i := 0
	for name, link := range p.links {
		if !strings.HasPrefix(link, "http:") && !strings.HasPrefix(link, "https:") {
			continue
		}

		i++
		if float64(i) > p.PageLimit {
			break
		}
		_, err := p.Page.Goto(link) // 打开链接
		if err != nil {
			return nil, fmt.Errorf("打开链接 %v 时出错: %v", link, err)
		}

		err = p.Page.Locator(selector).First().WaitFor(playwright.LocatorWaitForOptions{Timeout: p.timeout})
		if err != nil {
			fmt.Printf("等待页面加载超时：%v", err)
			continue
		}
		if elementType == "text" {
			text, err := p.GetTextContent(selector)
			if err != nil {
				return nil, fmt.Errorf("获取子页面链接时出错: %v", err)
			}
			result[name] = text
		} else {
			childLinks, err := p.getRecursiveElementsSources(p.Page.Locator(selector).First(), elementType)
			if err != nil {
				return nil, fmt.Errorf("获取子页面链接时出错: %v", err)
			}
			result[name] = childLinks
		}
	}
	p.links = nil
	return result, nil
}

// FormField 表单字段结构，包含英文字段名称和中文名称
type FormField struct {
	EnglishName string // 英文字段名称 (name属性)
	ChineseName string // 中文名称 (label文本)
	Value       string // 字段值
	IsRequired  bool   // 是否必填
}

// GetFormValues 获取表单中的所有字段及其值，包含中文名称和英文字段名称
func (p *PageToolkit) GetFormValues(selector string) (map[string]FormField, error) {
	form := p.Page.Locator(selector).First()
	if count, _ := form.Count(); count == 0 {
		return nil, fmt.Errorf("未找到选择器 %v 对应的表单", selector)
	}

	fields := form.Locator("input, select, textarea")
	count, err := fields.Count()
	if err != nil || count == 0 {
		return nil, fmt.Errorf("未找到表单中的任何字段")
	}

	values := make(map[string]FormField)
	for i := 0; i < count; i++ {
		field := fields.Nth(i)
		englishName, _ := field.GetAttribute("name")
		value, _ := field.InputValue()

		// 获取对应的中文名称
		labelLocator := p.Page.Locator(fmt.Sprintf("label[for='%s']", englishName))
		chineseName, _ := labelLocator.TextContent()

		// 如果找不到中文名称，则使用英文字段名称
		if chineseName == "" {
			chineseName = englishName
		}

		if englishName != "" {
			values[englishName] = FormField{
				EnglishName: englishName,
				ChineseName: chineseName,
				Value:       value,
			}
		}
	}
	return values, nil
}

// ValidateForm 验证表单，检查必填字段是否填写
func (p *PageToolkit) ValidateForm(formSelector string) error {
	formFields, err := p.GetFormValues(formSelector)
	if err != nil {
		return fmt.Errorf("获取表单值时出错: %v", err)
	}

	for _, field := range formFields {
		// 如果字段是必填且没有值，则返回错误
		if field.IsRequired && strings.TrimSpace(field.Value) == "" {
			return fmt.Errorf("字段 %s (%s) 为必填项，请填写", field.EnglishName, field.ChineseName)
		}
	}

	return nil
}

// GetTableContent 获取表格内容
func (p *PageToolkit) GetTableContent(selector string) ([][]string, error) {
	table := p.Page.Locator(selector).First()
	if count, _ := table.Count(); count == 0 {
		return nil, fmt.Errorf("未找到选择器 %v 对应的表格", selector)
	}

	rows := table.Locator("tr")
	rowCount, err := rows.Count()
	if err != nil || rowCount == 0 {
		return nil, fmt.Errorf("未找到表格中的任何行")
	}

	var tableContent [][]string
	for i := 0; i < rowCount; i++ {
		row := rows.Nth(i)
		cells := row.Locator("th, td")
		cellCount, _ := cells.Count()
		var rowContent []string
		for j := 0; j < cellCount; j++ {
			text, _ := cells.Nth(j).TextContent()
			rowContent = append(rowContent, text)
		}
		tableContent = append(tableContent, rowContent)
	}
	return tableContent, nil
}

// StartRecording 开始录制用户的操作
func (a *PageToolkit) StartRecording() {
	a.isRecording = true
	a.Page.On("click", func(event playwright.ElementHandle) {
		if a.isRecording {
			a.recordedActions = append(a.recordedActions, "await page.click('"+event.String()+"')")
		}
	})
	a.Page.On("keyboard", func(event playwright.ElementHandle) {
		fmt.Println("keyboard", event.String())
	})
}

// StopRecording 停止录制
func (a *PageToolkit) StopRecording() {
	a.isRecording = false
}

// RunScript 执行录制的脚本
func (a *PageToolkit) RunScript(script string) {
	// 执行录制的脚本
	_, err := a.Page.Evaluate(script)
	if err != nil {
	}
}

// html转换pdf
func (a *PageToolkit) htmlToPdf(locator playwright.Locator) string {
	html, err := locator.InnerHTML()
	if err != nil {
		log.Printf("创建PDF页面失败: %v\n", err)
		return ""
	}
	newPage, err := a.Browser.NewPage()
	if err != nil {
		log.Fatalf("创建PDF页面失败: %v\n", err)
		return ""
	}
	newPage.SetContent(html)
	//savePath := "temp/" + time.Now().Format("2006-01-02")
	//os.MkdirAll(savePath, 0666)
	options := playwright.PagePdfOptions{
		//Path:                playwright.String(path.Join(savePath, sno.GetString()+".pdf")), // 保存 PDF 文件
		Format:              playwright.String("A4"), // 页面格式
		DisplayHeaderFooter: playwright.Bool(false),  // 显示页眉和页脚
		PrintBackground:     playwright.Bool(false),  // 打印背景
	}
	data, err := newPage.PDF(options)
	if err != nil {
		fmt.Printf("转换生成文件PDF文件出错: %v\n", err)
		return ""
	}
	result, err := a.pdfToText(data)
	if err != nil {
		fmt.Printf("PDF 文本解析出错: %v\n", err)
		return ""
	}
	newPage.Close()
	return result
}

// pdf转换文本
func (a *PageToolkit) pdfToText(data []byte) (string, error) {
	r, _ := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	totalPage := r.NumPage()
	var mergedTexts []string
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		var mergedSentence string

		texts := p.Content().Text
		for _, text := range texts {
			if text.Y == lastTextStyle.Y {
				mergedSentence += text.S
			} else {
				if mergedSentence != "" {
					mergedTexts = append(mergedTexts, mergedSentence)
				}
				lastTextStyle = text
				mergedSentence = text.S
			}
		}

		if mergedSentence != "" {
			mergedTexts = append(mergedTexts, mergedSentence)
		}
	}

	mergedText := strings.Join(mergedTexts, "\n")
	return mergedText, nil
}
