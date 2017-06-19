// Package strumt provides a way to defines scenarios for prompting
// informations on command line
package strumt

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// NewPrompts creates a new prompt from stdin
func NewPrompts() Prompts {
	return Prompts{reader: bufio.NewReader(os.Stdin), prompts: map[string]Prompter{}}
}

// NewPromptsFromReader creates a new prompt from a given reader, useful for testing purpose
// for instance by providing a buffer
func NewPromptsFromReader(reader io.Reader) Prompts {
	return Prompts{reader: bufio.NewReader(reader), prompts: map[string]Prompter{}}
}

// Prompts stores all defined prompts and current
// running prompt
type Prompts struct {
	currentPrompt Prompter
	prompts       map[string]Prompter
	reader        *bufio.Reader
}

func (p *Prompts) parseMultipleLine(prompt MultilinePrompter) {
	inputs := []string{}

	for {
		input, err := p.reader.ReadString('\n')

		if err != nil {
			p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

			return
		}

		input = strings.TrimRight(input, "\n")
		inputs = append(inputs, input)

		bn, err := p.reader.ReadByte()

		if err == io.EOF {
			break
		}

		if bn == '\n' {
			break
		}

		if err := p.reader.UnreadByte(); err != nil {
			p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

			return
		}
	}

	if err := prompt.Parse(inputs); err != nil {
		p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

		return
	}

	if prompt.GetNextOnSuccess(inputs) == "" {
		p.currentPrompt = nil
		return
	}

	p.currentPrompt = p.prompts[prompt.GetNextOnSuccess(inputs)]
}

func (p *Prompts) parseLine(prompt LinePrompter) {
	input, err := p.reader.ReadString('\n')

	if err != nil {
		p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

		return
	}

	input = strings.TrimRight(input, "\n")

	if err := prompt.Parse(input); err != nil {

		p.currentPrompt = p.prompts[prompt.GetNextOnError(err)]

		return
	}

	if prompt.GetNextOnSuccess(input) == "" {
		p.currentPrompt = nil
		return
	}

	p.currentPrompt = p.prompts[prompt.GetNextOnSuccess(input)]
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
		prompt := p.currentPrompt

		switch lp := prompt.(type) {
		case LinePrompter:
			p.parseLine(lp)
		case MultilinePrompter:
			p.parseMultipleLine(lp)
		}

		if p.currentPrompt == nil {
			return
		}
	}
}
