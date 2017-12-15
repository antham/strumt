package strumt

import (
	"io"
)

// Prompter defines a generic common prompt.
//
// ID returns a string id to identify prompter and to let other prompter call it.
//
// PromptString returns a string to be diplayed as prompt.
//
// NextOnError is triggered when an error occurred during
// prompt sequence, it must returns the id of the prompt
// to be called when an error occured, most of the time it would
// be the id of the current prompt to loop on it
type Prompter interface {
	ID() string
	PromptString() string
	NextOnError(error) string
}

// LinePrompter defines a one line prompter
// that will ask only for one user input.
//
// NextOnSuccess must returns the id of the next prompt
// to be called. To mark prompter as the last prompter
// NextOnSucces must returns an empty string
type LinePrompter interface {
	Prompter
	NextOnSuccess(string) string
	Parse(string) error
}

// MultilinePrompter defines a mutiline prompter
// that will let the possibility to the user to
// provide several input, result is provided as
// an input slice.
//
// NextOnSuccess must returns the id of the next prompt
// to be called. To mark prompter as the last prompter
// NextOnSucces must returns an empty string
type MultilinePrompter interface {
	Prompter
	NextOnSuccess([]string) string
	Parse([]string) error
}

// PromptRenderer can be implemented to customize
// the way prompt is rendered, PromptString result is given as parameter
type PromptRenderer interface {
	PrintPrompt(io.Writer, string)
}

// ErrorRenderer can be implemented to customize
// the way an error returned by Parse is rendered
type ErrorRenderer interface {
	PrintError(io.Writer, error)
}

// SeparatorRenderer can be implemented to customize
// the way a prompt is separated from another, default
// is to add a new line
type SeparatorRenderer interface {
	PrintSeparator(io.Writer)
}
