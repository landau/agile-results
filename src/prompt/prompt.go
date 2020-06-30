package prompt

import (
	"bufio"
	"strings"
)

// Prompter -
type Prompter interface {
	Prompt(s string) (string, error)
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
	// FIXME: pull out result and strip trailing \n
	res, err := p.reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimRight(res, "\n"), nil
}

// TODO: shortcut for writing to stdin/out would be NewStdInOutPrompter

// New - A prompter for strings
func New(reader *bufio.Reader, writer *bufio.Writer) *IOPrompter {
	return &IOPrompter{reader: reader, writer: writer}
}
