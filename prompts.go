// Package strumt provides a way to defines scenarios for prompting
// informations on command line
package strumt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// NewPrompts creates a new prompt from stdin
func NewPrompts() Prompts {
	return Prompts{reader: bufio.NewReader(os.Stdin), writer: os.Stdout, prompts: map[string]Prompter{}}
}

// NewPromptsFromReader creates a new prompt from a given reader and writer
// , useful for testing purpose for instance by providing a buffer
func NewPromptsFromReader(reader io.Reader, writer io.Writer) Prompts {
	return Prompts{reader: bufio.NewReader(reader), writer: writer, prompts: map[string]Prompter{}}
}

// Prompts stores all defined prompts and current
// running prompt
type Prompts struct {
	currentPrompt Prompter
	prompts       map[string]Prompter
	reader        *bufio.Reader
	writer        io.Writer
}

func (p *Prompts) renderPrompt(prompt Prompter) {
	switch pr := prompt.(type) {
	default:
		fmt.Fprintf(p.writer, "%s : \n", prompt.GetPromptString())
	}
}
func (p *Prompts) renderError(prompt Prompter, err error) {
	switch pr := prompt.(type) {
	default:
		fmt.Fprintf(p.writer, "%s\n", err.Error())
	}
}

func (p *Prompts) isMultilineEnd() (bool, error) {
	bn, err := p.reader.ReadByte()

	if err == io.EOF {
		return true, nil
	}

	if bn == '\n' {
		return true, nil
	}

	if err := p.reader.UnreadByte(); err != nil {
		return false, err
	}

	return false, nil
}

func (p *Prompts) parseMultipleLine(prompt MultilinePrompter) ([]string, error) {
	inputs := []string{}

	for {
		input, err := p.reader.ReadString('\n')
		input = strings.TrimRight(input, "\n")

		if err != nil {
			p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

			return []string{}, err
		}

		inputs = append(inputs, input)

		end, err := p.isMultilineEnd()

		if err != nil {
			p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

			return []string{}, err
		}

		if end {
			break
		}
	}

	if err := prompt.Parse(inputs); err != nil {
		p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

		return inputs, err
	}

	if prompt.GetNextOnSuccess(inputs) == "" {
		p.currentPrompt = nil

		return inputs, nil
	}

	p.currentPrompt = p.prompts[prompt.GetNextOnSuccess(inputs)]

	return inputs, nil
}

func (p *Prompts) parseLine(prompt LinePrompter) (string, error) {
	input, err := p.reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")

	if err != nil {
		p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

		return "", err
	}

	if err := prompt.Parse(input); err != nil {
		p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

		return input, err
	}

	if prompt.GetNextOnSuccess(input) == "" {
		p.currentPrompt = nil

		return input, nil
	}

	p.currentPrompt = p.prompts[prompt.GetNextOnSuccess(input)]

	return input, nil
}

// AddLinePrompter add a new LinePrompter mapped to a given id
func (p *Prompts) AddLinePrompter(id string, prompt LinePrompter) {
	p.prompts[id] = prompt
}

// AddMultilinePrompter add a new MultilinePrompter mapped to a given id
func (p *Prompts) AddMultilinePrompter(id string, prompt MultilinePrompter) {
	p.prompts[id] = prompt
}

// SetFirst defines from which prompt, the prompt sequence has to start
func (p *Prompts) SetFirst(id string) {
	p.currentPrompt = p.prompts[id]
}

// Run executes prompt sequence
func (p *Prompts) Run() {
	for {
		var err error

		prompt := p.currentPrompt
		p.renderPrompt(prompt)

		switch lp := prompt.(type) {
		case LinePrompter:
			_, err = p.parseLine(lp)
		case MultilinePrompter:
			_, err = p.parseMultipleLine(lp)
		}

		if err != nil {
			p.renderError(prompt, err)
		}

		if p.currentPrompt == nil {
			return
		}
	}
}
