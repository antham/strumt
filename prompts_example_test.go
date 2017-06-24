package strumt_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/antham/strumt"
)

type StringPrompt struct {
	store             *string
	prompt            string
	nextPrompt        string
	nextPromptOnError string
}

func (s *StringPrompt) PromptString() string {
	return s.prompt
}

func (s *StringPrompt) Parse(value string) error {
	if value == "" {
		return fmt.Errorf("Empty value given")
	}

	*(s.store) = value

	return nil
}

func (s *StringPrompt) NextOnSuccess(value string) string {
	return s.nextPrompt
}

func (s *StringPrompt) NextOnError(err error) string {
	return s.nextPromptOnError
}

type IntPrompt struct {
	store             *int
	prompt            string
	nextPrompt        string
	nextPromptOnError string
}

func (i *IntPrompt) PromptString() string {
	return i.prompt
}

func (i *IntPrompt) Parse(value string) error {
	age, err := strconv.Atoi(value)

	if err != nil {
		return fmt.Errorf("%s is not a valid number", value)
	}

	if age <= 0 {
		return fmt.Errorf("Give a valid age")
	}

	*(i.store) = age

	return nil
}

func (i *IntPrompt) NextOnSuccess(value string) string {
	return i.nextPrompt
}

func (i *IntPrompt) NextOnError(err error) string {
	return i.nextPromptOnError
}

type User struct {
	FirstName string
	LastName  string
	Age       int
}

func Example() {
	user := User{}

	buf := "\nBrad\n\nBlanton\nwhatever\n0\n31\n"

	p := strumt.NewPromptsFromReaderAndWriter(bytes.NewBufferString(buf), ioutil.Discard)
	p.AddLinePrompter("userName", &StringPrompt{&user.FirstName, "Enter your first name", "lastName", "userName"})
	p.AddLinePrompter("lastName", &StringPrompt{&user.LastName, "Enter your last name", "age", "lastName"})
	p.AddLinePrompter("age", &IntPrompt{&user.Age, "Enter your age", "", "age"})
	p.SetFirst("userName")
	p.Run()

	for _, step := range p.Scenario() {
		fmt.Println(step.PromptString())
		fmt.Println(step.Inputs()[0])

		if step.Error() != nil {
			fmt.Println(step.Error())
		}
	}

	fmt.Println()
	fmt.Printf("User datas : %#v", user)

	// Output:
	// Enter your first name
	//
	// Empty value given
	// Enter your first name
	// Brad
	// Enter your last name
	//
	// Empty value given
	// Enter your last name
	// Blanton
	// Enter your age
	// whatever
	// whatever is not a valid number
	// Enter your age
	// 0
	// Give a valid age
	// Enter your age
	// 31
	//
	// User datas : strumt_test.User{FirstName:"Brad", LastName:"Blanton", Age:31}
}
