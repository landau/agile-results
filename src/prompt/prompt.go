package prompt

import (
	"bufio"
	"strings"
)

// Prompter -
type Prompter interface {
	Prompt(s string) (string, error)
	PromptList(s string) ([]string, error)
}

// IOPrompter - Use this to interact with users from input to output
type IOPrompter struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

// Prompt - Prompts for a string
func (p *IOPrompter) Prompt(s string) (string, error) {
	p.writer.WriteString(s)
	p.writer.Flush()
	res, err := p.reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimRight(res, "\n"), nil
}

// PromptList -
func (p *IOPrompter) PromptList(s string) ([]string, error) {
	res, err := p.Prompt(s)

	list := make([]string, 0)

	if err != nil {
		return list, err
	}

	return strings.Split(res, ","), nil
}

type (
	// MockPrompter -
	MockPrompter struct {
		PromptCalls           []string
		PromptReturnValue     MockPrompterPromptReturnValue
		PromptListCalls       []string
		PromptListReturnValue MockPrompterPromptListReturnValue
	}

	// MockPrompterPromptReturnValue -
	MockPrompterPromptReturnValue struct {
		S   string
		Err error
	}

	// MockPrompterPromptListReturnValue -
	MockPrompterPromptListReturnValue struct {
		Items []string
		Err   error
	}
)

// Prompt -
func (p *MockPrompter) Prompt(s string) (string, error) {
	p.PromptCalls = append(p.PromptCalls, s)
	return p.PromptReturnValue.S, p.PromptReturnValue.Err
}

// PromptList -
func (p *MockPrompter) PromptList(s string) ([]string, error) {
	p.PromptListCalls = append(p.PromptListCalls, s)
	return p.PromptListReturnValue.Items, p.PromptListReturnValue.Err
}

// TODO: shortcut for writing to stdin/out would be NewStdInOutPrompter

// New - A prompter for strings
func New(reader *bufio.Reader, writer *bufio.Writer) *IOPrompter {
	return &IOPrompter{reader: reader, writer: writer}
}
