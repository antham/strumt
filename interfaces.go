package strumt

import (
	"io"
)

// Prompter defines a generic common prompt.
//
// ID returns a string ID to identify prompter and to let other prompter call it.
//
// PromptString returns a string to be displayed as prompt.
//
// NextOnError is triggered when an error occurred during
// prompt sequence, it must returns the ID of the prompt
// to be called when an error occurred, most of the time it would
// be the ID of the current prompt to loop on it
type Prompter interface {
	ID() string
	PromptString() string
	NextOnError(error) string
}

// LinePrompter defines a one line prompter
// that will ask only for a single line user input.
//
// NextOnSuccess must returns the ID of the next prompt
// to be called. To mark prompter as the last prompter,
// NextOnSuccess must returns an empty string
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
// NextOnSuccess must returns the ID of the next prompt
// to be called. To mark prompter as the last prompter,
// NextOnSuccess must returns an empty string
type MultilinePrompter interface {
	Prompter
	NextOnSuccess([]string) string
	Parse([]string) error
}

// PromptRenderer can be implemented to customize
// the way prompt is rendered, original PromptString is given
// as second parameter
type PromptRenderer interface {
	PrintPrompt(io.Writer, string)
}

// ErrorRenderer can be implemented to customize
// the way an error returned by Parse is rendered
type ErrorRenderer interface {
	PrintError(io.Writer, error)
}

// SeparatorRenderer can be implemented to customize
// the way a prompt is separated from another. When
// this interface is not implemented, the default behaviour
// is to define a new line as separator
type SeparatorRenderer interface {
	PrintSeparator(io.Writer)
}
