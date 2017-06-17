package strumt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Datas struct {
	Db struct {
		Username string `strumtP:"Enter a valid username"`
		Password string `strumtP:"Enter a valid password"`
		Port     int `strumtP:"Enter a valid port"`

	}
	Hosts map[string]string `strumtP:"Enter host/ip couples"`
	Ips []string `strumtP:"Enter ips"`
}

func TestPopulate(t *testing.T) {
	actual := &Datas{}
	Run(actual)

	expected := &Datas{}
	expected.Db.Username = "user"
	expected.Db.Password = "password"
	expected.Db.Port = 10000
	expected.Hosts = map[string]string{
		"localhost": "127.0.0.1",
			"myIp": "1.2.3.4",
	}
	expected.Ips = []string{
		"127.0.0.1",
		"1.2.3.4",
		"8.9.10.11",
	}

	assert.Equal(t, expected, actual)
}
