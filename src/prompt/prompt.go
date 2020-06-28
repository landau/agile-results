package prompt

import (
	"bufio"
)

// StringPrompter - Use this to interact with users from input to output
type StringPrompter struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

// Prompt - Prompts for a string
func (p *StringPrompter) Prompt(s string) (res string, err error) {
	p.writer.WriteString(s)
	p.writer.Flush()
	return p.reader.ReadString('\n')
}

// NewStringPrompter - A prompter for strings
func NewStringPrompter(reader *bufio.Reader, writer *bufio.Writer) StringPrompter {
	return StringPrompter{reader: reader, writer: writer}
}
