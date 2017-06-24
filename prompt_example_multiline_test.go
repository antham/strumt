package strumt_test

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/antham/strumt"
)

type SlicePrompt struct {
	datas *[]string
}

func (s *SlicePrompt) PromptString() string {
	return "Give several input"
}

func (s *SlicePrompt) Parse(values []string) error {
	*(s.datas) = values

	return nil
}

func (s *SlicePrompt) NextOnSuccess(values []string) string {
	return ""
}

func (s *SlicePrompt) NextOnError(err error) string {
	return "sliceprompt"
}

func Example_multilinePrompt() {
	var datas []string
	buf := "test1\ntest2\ntest3\ntest4\n\n"

	p := strumt.NewPromptsFromReaderAndWriter(bytes.NewBufferString(buf), ioutil.Discard)
	p.AddMultilinePrompter("sliceprompt", &SlicePrompt{&datas})
	p.SetFirst("sliceprompt")
	p.Run()

	fmt.Println(datas)
	// Output:
	// [test1 test2 test3 test4]
}
