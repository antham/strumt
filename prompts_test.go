package strumt

import (
	"bytes"
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
		return err
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

type SlicePrompt struct {
	valuePtr          *[]string
	prompt            string
	nextPrompt        string
	nextPromptOnError string
}

func (s *SlicePrompt) GetPromptString() string {
	return s.prompt
}

func (s *SlicePrompt) Parse(values []string) error {
	(*s.valuePtr) = values

	return nil
}

func (s *SlicePrompt) GetNextOnError(err error) string {
	return s.nextPromptOnError
}

func (s *SlicePrompt) GetNextOnSuccess(value []string) string {
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

	buf := []byte("user\npassword\n10000\n127.0.0.1\n1.2.3.4\n8.9.10.11\n\nlocalhost:127.0.0.1\nmyIp:1.2.3.4\n\n")

	p := NewPromptsFromReader(bytes.NewBuffer(buf))

	p.AddLinePrompter("username", &StringPrompt{&actual.Db.Username, "Give a username", "password", "username"})
	p.AddLinePrompter("password", &StringPrompt{&actual.Db.Password, "Give a password", "port", "password"})
	p.AddLinePrompter("port", &IntPrompt{&actual.Db.Port, "Give a port", "ips", "port"})
	p.AddMultilinePrompter("ips", &SlicePrompt{&actual.Ips, "Give some ips", "hosts", "ips"})
	p.AddMultilinePrompter("hosts", &MapPrompt{&actual.Hosts, "Give some hosts", "", "hosts"})

	p.SetFirst("username")
	p.Run()

	assert.Equal(t, expected, actual)
}
