package prompt

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestIOPrompter_Prompt(t *testing.T) {
	type prompterFields struct {
		reader *bufio.Reader
		writer *bufio.Writer
	}

	type args struct {
		s string
	}

	readBuffer := &bytes.Buffer{}
	readBuffer.WriteString("foo\n")

	tests := []struct {
		name    string
		fields  prompterFields
		args    args
		wantRes string
		wantErr bool
	}{
		{
			"Writes to a string writer and reads from a reader",
			prompterFields{
				reader: bufio.NewReader(readBuffer),
				writer: bufio.NewWriter(&bytes.Buffer{}),
			},
			args{s: "foobar"},
			"foo",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &IOPrompter{
				reader: tt.fields.reader,
				writer: tt.fields.writer,
			}

			gotRes, err := p.Prompt(tt.args.s)

			if (err != nil) != tt.wantErr {
				t.Errorf("StringPrompter.Prompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotRes != tt.wantRes {
				t.Errorf("StringPrompter.Prompt() = %v, want %v", gotRes, tt.wantRes)
			}

			// TODO: verify that the writer was written to with expected text
		})
	}
}

func TestIOPrompter_PromptList(t *testing.T) {
	type prompterFields struct {
		reader *bufio.Reader
		writer *bufio.Writer
	}

	type args struct {
		s string
	}

	readBuffer := &bytes.Buffer{}
	readBuffer.WriteString("foo,bar\n")

	tests := []struct {
		name    string
		fields  prompterFields
		args    args
		wantRes []string
		wantErr bool
	}{
		{
			"Writes to a string writer and reads from a reader",
			prompterFields{
				reader: bufio.NewReader(readBuffer),
				writer: bufio.NewWriter(&bytes.Buffer{}),
			},
			args{s: "foobar"},
			[]string{"foo", "bar"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &IOPrompter{
				reader: tt.fields.reader,
				writer: tt.fields.writer,
			}

			gotRes, err := p.PromptList(tt.args.s)

			if (err != nil) != tt.wantErr {
				t.Errorf("StringPrompter.PromptList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("StringPrompter.PromptList() = %v, want %v", gotRes, tt.wantRes)
			}

			// TODO: verify that the writer was written to with expected text
		})
	}
}
