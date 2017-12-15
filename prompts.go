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

// Step represents a scenario step which is
// the result of a prompt execution. We store
// the prompt string, inputs that the user has given,
// and the prompt error if one occurred
type Step struct {
	prompt string
	inputs []string
	err    error
}

// PromptString returns prompt string displayed by
// the prompt string
func (s Step) PromptString() string {
	return s.prompt
}

// Inputs retrieves all inputs given by user,
func (s Step) Inputs() []string {
	return s.inputs
}

// Error returns the error, if any,
// triggered on a prompt error
func (s Step) Error() error {
	return s.err
}

// NewPrompts creates a new prompt from stdin
func NewPrompts() Prompts {
	return Prompts{reader: bufio.NewReader(os.Stdin), writer: os.Stdout, prompts: map[string]Prompter{}}
}

// NewPromptsFromReaderAndWriter creates a new prompt from a given reader and writer, useful for testing purpose
func NewPromptsFromReaderAndWriter(reader io.Reader, writer io.Writer) Prompts {
	return Prompts{reader: bufio.NewReader(reader), writer: writer, prompts: map[string]Prompter{}}
}

// Prompts stores all defined prompts and current
// running prompt
type Prompts struct {
	currentPrompt Prompter
	prompts       map[string]Prompter
	reader        *bufio.Reader
	writer        io.Writer
	scenario      []Step
}

func (p *Prompts) parse() ([]string, Prompter, error) {
	var nextPrompt Prompter
	var inputs []string
	var err error

	switch prompt := p.currentPrompt.(type) {
	case LinePrompter:
		var input string

		input, err = parseLine(p.reader, prompt)

		if prompt.NextOnSuccess(input) != "" {
			nextPrompt = p.prompts[prompt.NextOnSuccess(input)]
		}

		inputs = append(inputs, input)
	case MultilinePrompter:
		inputs, err = parseMultipleLine(p.reader, prompt)

		if prompt.NextOnSuccess(inputs) != "" {
			nextPrompt = p.prompts[prompt.NextOnSuccess(inputs)]
		}
	}

	if err != nil {
		nextPrompt = p.prompts[p.currentPrompt.NextOnError(err)]
	}

	return inputs, nextPrompt, err
}

func (p *Prompts) appendScenario(promptString string, inputs []string, err error) {
	p.scenario = append(
		p.scenario,
		Step{
			promptString,
			inputs,
			err,
		},
	)
}

// AddLinePrompter add a new LinePrompter mapped to a given id
func (p *Prompts) AddLinePrompter(prompt LinePrompter) {
	p.prompts[prompt.ID()] = prompt
}

// AddMultilinePrompter add a new MultilinePrompter mapped to a given id
func (p *Prompts) AddMultilinePrompter(prompt MultilinePrompter) {
	p.prompts[prompt.ID()] = prompt
}

// SetFirst defines from which prompt the prompt sequence has to start
func (p *Prompts) SetFirst(id string) {
	p.currentPrompt = p.prompts[id]
}

// Scenario retrieves all steps done during
// a prompt sequence
func (p *Prompts) Scenario() []Step {
	return p.scenario
}

// Run executes a prompt sequence
func (p *Prompts) Run() {
	p.scenario = []Step{}

	for {
		prompt := p.currentPrompt
		renderPrompt(p.writer, prompt)

		inputs, nextPrompt, err := p.parse()

		if err != nil {
			renderError(p.writer, prompt, err)
		}

		p.appendScenario(prompt.PromptString(), inputs, err)

		if nextPrompt == nil {
			return
		}

		renderSeparator(p.writer, prompt)

		p.currentPrompt = nextPrompt
	}
}

func isMultilineEnd(reader io.ByteScanner) (bool, error) {
	bn, err := reader.ReadByte()

	if err == io.EOF {
		return true, nil
	}

	if bn == '\n' {
		return true, nil
	}

	if err := reader.UnreadByte(); err != nil {
		return false, err
	}

	return false, nil
}

func parseMultipleLine(reader *bufio.Reader, prompt MultilinePrompter) ([]string, error) {
	inputs := []string{}

	for {
		input, err := reader.ReadString('\n')
		input = strings.TrimRight(input, "\n")

		if err != nil {
			return []string{}, err
		}

		inputs = append(inputs, input)

		end, err := isMultilineEnd(reader)

		if err != nil {
			return []string{}, err
		}

		if end {
			break
		}
	}

	return inputs, prompt.Parse(inputs)
}

func parseLine(reader *bufio.Reader, prompt LinePrompter) (string, error) {
	input, err := reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")

	if err != nil {
		return "", err
	}

	return input, prompt.Parse(input)
}

func renderPrompt(writer io.Writer, prompt Prompter) {
	switch pr := prompt.(type) {
	case PromptRenderer:
		pr.PrintPrompt(prompt.PromptString())
	default:
		fmt.Fprintf(writer, "%s\n", prompt.PromptString())
	}
}

func renderError(writer io.Writer, prompt Prompter, err error) {
	switch pr := prompt.(type) {
	case ErrorRenderer:
		pr.PrintError(err)
	default:
		fmt.Fprintf(writer, "%s\n", err.Error())
	}
}

func renderSeparator(writer io.Writer, prompt Prompter) {
	switch pr := prompt.(type) {
	case SeparatorRenderer:
		pr.PrintSeparator()
	default:
		fmt.Fprintf(writer, "\n")
	}
}
