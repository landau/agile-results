package prompt

import (
	"bufio"
	"strings"
)

// Prompter - Use this to interact with users from input to output
type Prompter struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

// Prompt - Prompts for a string
func (p *Prompter) Prompt(s string) (string, error) {
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
func New(reader *bufio.Reader, writer *bufio.Writer) Prompter {
	return Prompter{reader: reader, writer: writer}
}
