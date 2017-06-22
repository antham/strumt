package strumt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Datas struct {
	Db struct {
		Username string
		Password string
		Port     int
	}
	Hosts map[string]string
	Ips   []string
}

type StringPrompt struct {
	valuePtr          *string
	prompt            string
	nextPrompt        string
	nextPromptOnError string
}

func (s *StringPrompt) GetPromptString() string {
	return s.prompt
}

func (s *StringPrompt) Parse(value string) error {
	if value == "" {
		return fmt.Errorf("Empty value given")
	}

	*(s.valuePtr) = value
	return nil
}

func (s *StringPrompt) GetNextOnSuccess(value string) string {
	return s.nextPrompt
}

func (s *StringPrompt) GetNextOnError(err error) string {
	return s.nextPromptOnError
}

type IntPrompt struct {
	valuePtr          *int
	prompt            string
	nextPrompt        string
	nextPromptOnError string
}

func (s *IntPrompt) GetPromptString() string {
	return s.prompt
}

func (s *IntPrompt) Parse(value string) error {
	v, err := strconv.Atoi(value)

	if err != nil {
		return fmt.Errorf("Provide a numerical value")
	}

	*(s.valuePtr) = v

	return nil
}

func (s *IntPrompt) GetNextOnSuccess(value string) string {
	return s.nextPrompt
}

func (s *IntPrompt) GetNextOnError(err error) string {
	return s.nextPromptOnError
}

type IpsPrompt struct {
	valuePtr          *[]string
	prompt            string
	nextPrompt        string
	nextPromptOnError string
}

func (s *IpsPrompt) GetPromptString() string {
	return s.prompt
}

func (s *IpsPrompt) Parse(values []string) error {
	for _, value := range values {
		if net.ParseIP(value) == nil {
			return fmt.Errorf("%s is not a valid IP", value)
		}
	}

	(*s.valuePtr) = values

	return nil
}

func (s *IpsPrompt) GetNextOnError(err error) string {
	return s.nextPromptOnError
}

func (s *IpsPrompt) GetNextOnSuccess(value []string) string {
	return s.nextPrompt
}

type MapPrompt struct {
	valuePtr          *map[string]string
	prompt            string
	nextPrompt        string
	nextPromptOnError string
}

func (m *MapPrompt) GetPromptString() string {
	return m.prompt
}

func (m *MapPrompt) Parse(values []string) error {
	for _, data := range values {
		keyValue := strings.Split(data, ":")

		if len(keyValue) != 2 {
			return fmt.Errorf("Check %s is a valid couple key:value", data)
		}

		(*m.valuePtr)[keyValue[0]] = keyValue[1]
	}

	return nil
}

func (m *MapPrompt) GetNextOnError(err error) string {
	return m.nextPromptOnError
}

func (m *MapPrompt) GetNextOnSuccess(value []string) string {
	return m.nextPrompt
}

func TestNewPrompts(t *testing.T) {
	p := NewPrompts()

	assert.Equal(t, bufio.NewReader(os.Stdin), p.reader)
	assert.Equal(t, os.Stdout, p.writer)
}

func TestPromptsRun(t *testing.T) {
	actual := &Datas{
		Hosts: map[string]string{},
		Ips:   []string{},
	}

	expected := &Datas{}
	expected.Db.Username = "user"
	expected.Db.Password = "password"
	expected.Db.Port = 10000
	expected.Hosts = map[string]string{
		"localhost": "127.0.0.1",
		"myIp":      "1.2.3.4",
	}
	expected.Ips = []string{
		"127.0.0.1",
		"1.2.3.4",
		"8.9.10.11",
	}

	buf := "\nuser\n\npassword\ntest\n10000\n127.0.0.1\ntest\n1.2.3.4\n8.9.10.11\n\n127.0.0.1\n1.2.3.4\n8.9.10.11\n\nlocalhost:127.0.0.1\ntest\nmyIp:1.2.3.4\n\nlocalhost:127.0.0.1\nmyIp:1.2.3.4\n\n"

	var actualStdout bytes.Buffer

	p := NewPromptsFromReader(bytes.NewBufferString(buf), &actualStdout)

	p.AddLinePrompter("username", &StringPrompt{&actual.Db.Username, "Give a username", "password", "username"})
	p.AddLinePrompter("password", &StringPrompt{&actual.Db.Password, "Give a password", "port", "password"})
	p.AddLinePrompter("port", &IntPrompt{&actual.Db.Port, "Give a port", "ips", "port"})
	p.AddMultilinePrompter("ips", &IpsPrompt{&actual.Ips, "Give some ips", "hosts", "ips"})
	p.AddMultilinePrompter("hosts", &MapPrompt{&actual.Hosts, "Give some host/ip couples", "", "hosts"})

	p.SetFirst("username")
	p.Run()

	expectedStdout := "Give a username : \nEmpty value given\nGive a username : \n" +
		"Give a password : \nEmpty value given\nGive a password : \n" +
		"Give a port : \nProvide a numerical value\nGive a port : \n" +
		"Give some ips : \ntest is not a valid IP\nGive some ips : \n" +
		"Give some host/ip couples : \nCheck test is a valid couple key:value\n" +
		"Give some host/ip couples : \n"

	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedStdout, actualStdout.String())
}

type StringWithCustomRendererPrompt struct {
	valuePtr          *string
	prompt            string
	nextPrompt        string
	nextPromptOnError string
	writer            io.Writer
}

func (s *StringWithCustomRendererPrompt) GetPromptString() string {
	return s.prompt
}

func (s *StringWithCustomRendererPrompt) Parse(value string) error {
	if value == "" {
		return fmt.Errorf("empty value given")
	}

	*(s.valuePtr) = value

	return nil
}

func (s *StringWithCustomRendererPrompt) GetNextOnSuccess(value string) string {
	return s.nextPrompt
}

func (s *StringWithCustomRendererPrompt) GetNextOnError(err error) string {
	return s.nextPromptOnError
}

func (s *StringWithCustomRendererPrompt) PrintPrompt(prompt string) {
	fmt.Fprintf(s.writer, "==> %s : \n", prompt)
}

func (s *StringWithCustomRendererPrompt) PrintError(err error) {
	fmt.Fprintf(s.writer, "==> Something went wrong : %s\n", err.Error())
}

func TestPromptRunWithCustomRenderer(t *testing.T) {
	var actualStdout bytes.Buffer

	p := NewPromptsFromReader(bytes.NewBufferString("\ntest\n"), ioutil.Discard)

	var value string

	p.AddLinePrompter("test", &StringWithCustomRendererPrompt{&value, "Give a value", "", "test", &actualStdout})

	p.SetFirst("test")
	p.Run()

	assert.Equal(t, "test", value)
	assert.Equal(t, "==> Give a value : \n==> Something went wrong : empty value given\n==> Give a value : \n", actualStdout.String())
}
