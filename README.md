# Strumt [![codecov](https://codecov.io/gh/antham/strumt/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/strumt) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/strumt)](https://goreportcard.com/report/github.com/antham/strumt) [![GoDoc](https://godoc.org/github.com/antham/strumt?status.svg)](http://godoc.org/github.com/antham/strumt) [![GitHub tag](https://img.shields.io/github/tag/antham/strumt.svg)]()

Strumt is a library to create prompt chain. It provides multiline prompt, input validation, retry on error, ability to create typesafe prompt, ability to customize prompt and error display, a recording of prompt session and the ability to easily test your prompt scenario.

## Example

Checkout godoc to have more examples : [https://godoc.org/github.com/antham/strumt](https://godoc.org/github.com/antham/strumt)

---

[![asciicast](https://asciinema.org/a/126121.png)](https://asciinema.org/a/126121)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"

    "github.com/antham/strumt"
)

func main() {
    user := User{}

    p := strumt.NewPromptsFromReaderAndWriter(bufio.NewReader(os.Stdin), os.Stdout)
    p.AddLinePrompter(&StringPrompt{&user.FirstName, "Enter your first name", "userName", "lastName", "userName"})
    p.AddLinePrompter(&StringPrompt{&user.LastName, "Enter your last name", "lastName", "age", "lastName"})
    p.AddLinePrompter(&IntPrompt{&user.Age, "Enter your age", "age", "", "age"})
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

}

type StringPrompt struct {
    store             *string
    prompt            string
    currentID         string
    nextPrompt        string
    nextPromptOnError string
}

func (s *StringPrompt) ID() string {
	return s.currentID
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
    currentID         string
    nextPrompt        string
    nextPromptOnError string
}

func (i *IntPrompt) ID() string {
	return i.currentID
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
```
