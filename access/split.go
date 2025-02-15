package access

import (
	"regexp"
	"strings"
)

type BasicSplitProvider struct{}

func NewBasicSplitProvider() (*BasicSplitProvider, error) {
	return &BasicSplitProvider{}, nil
}

func (p *BasicSplitProvider) SplitText(text string) ([]string, error) {
	const maxLength = 210
	var res []string
	var temp string

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if len(temp) <= maxLength {
			if temp != "" {
				temp += "\n"
			}
			temp += line
		} else if temp != "" {
			res = append(res, temp)
		}
		temp = line
	}
	if temp != "" {
		if len(temp) < 300 && len(res) > 0 {
			res[len(res)-1] += temp
		} else {
			res = append(res, temp)
		}
	}
	return res, nil
}

type DefaultSplitProvider struct{}

func NewDefaultSplitProvider() (*DefaultSplitProvider, error) {
	return &DefaultSplitProvider{}, nil
}

func (p *DefaultSplitProvider) SplitText(text string) ([]string, error) {
	const maxLength = 210
	var sections []string
	var currentSection strings.Builder
	var codeBlock strings.Builder
	inCodeBlock := false
	codeBlockLines := 0

	lines := strings.Split(text, "\n")
	emptyLineCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			emptyLineCount++
			if emptyLineCount >= 4 && currentSection.Len() > 0 {
				sections = append(sections, currentSection.String())
				currentSection.Reset()
			}
			continue
		} else {
			emptyLineCount = 0
		}

		if line == "```" {
			if inCodeBlock {
				inCodeBlock = false
				if codeBlockLines >= 5 {
					if currentSection.Len() > 0 {
						sections = append(sections, currentSection.String())
						currentSection.Reset()
					}
					sections = append(sections, codeBlock.String())
				} else {
					currentSection.WriteString(codeBlock.String())
				}
				codeBlock.Reset()
				codeBlockLines = 0
			} else {
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			codeBlock.WriteString(line + "\n")
			codeBlockLines++
			if codeBlockLines >= 20 {
				if currentSection.Len() > 0 {
					sections = append(sections, currentSection.String())
					currentSection.Reset()
				}
				sections = append(sections, codeBlock.String())
				codeBlock.Reset()
				codeBlockLines = 0
			}
			continue
		}

		if isSectionSeparator(line) {
			if currentSection.Len() > 0 {
				sections = append(sections, currentSection.String())
				currentSection.Reset()
			}
			currentSection.WriteString(line + "\n")
			continue
		}

		if currentSection.Len() <= maxLength {
			if currentSection.Len() > 0 {
				currentSection.WriteString("\n")
			}
			currentSection.WriteString(line)
		} else {
			if currentSection.Len() > 0 {
				sections = append(sections, currentSection.String())
				currentSection.Reset()
			}
			currentSection.WriteString(line)
		}
	}

	if currentSection.Len() > 0 {
		sections = append(sections, currentSection.String())
	}

	return sections, nil
}

func isSectionSeparator(line string) bool {
	if strings.HasPrefix(line, "Chapter") || strings.HasPrefix(line, "Section") {
		return true
	}
	matched, _ := regexp.MatchString(`^\d+\.\s`, line)
	return matched
}

type SplitProvider interface {
	SplitText(text string) ([]string, error)
}

// QaSplitProvider structure
type QaSplitProvider struct{}

// NewQaSplitProvider creates a new instance of QaSplitProvider
func NewQaSplitProvider() (*QaSplitProvider, error) {
	return &QaSplitProvider{}, nil
}

// SplitText method splits the text into question-answer pairs
func (p *QaSplitProvider) SplitText(text string) ([]string, error) {
	var res []string
	var currentPair string
	var collectingAnswer bool

	// Split the text by lines
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Q:") {
			if currentPair != "" {
				// Save the previous question-answer pair
				res = append(res, currentPair)
			}
			currentPair = line
			collectingAnswer = false
		} else if strings.HasPrefix(line, "A:") {
			currentPair += "\n" + line
			collectingAnswer = true
		} else if collectingAnswer {
			currentPair += "\n" + line
		}
	}

	// Add the last question-answer pair
	if currentPair != "" {
		res = append(res, currentPair)
	}

	return res, nil
}

func GetSplitProvider(typ string) (SplitProvider, error) {
	var p SplitProvider
	var err error
	if typ == "Default" {
		p, err = NewDefaultSplitProvider()
	} else if typ == "QA" {
		p, err = NewQaSplitProvider()
	} else if typ == "Basic" {
		p, err = NewBasicSplitProvider()
	} else {
		p, err = NewDefaultSplitProvider()
	}

	if err != nil {
		return nil, err
	}
	return p, nil
}
