package strumt

// Prompter defines a generic common prompt
// ID returns a string id to identify prompter
// PromptString returns a string to diplay as prompt
// NextOnError is triggered when an error occurred during
// prompt sequence
type Prompter interface {
	ID() string
	PromptString() string
	NextOnError(error) string
}

// LinePrompter defines a one line prompter
// that will ask only for one user input.
// To mark prompter as the last prompter
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
// To mark prompter as the last prompter
// NextOnSucces must returns an empty string
type MultilinePrompter interface {
	Prompter
	NextOnSuccess([]string) string
	Parse([]string) error
}

// PromptRenderer can be implemented to customize
// the way prompt is rendered, PromptString result is given as parameter
type PromptRenderer interface {
	PrintPrompt(string)
}

// ErrorRenderer can be implemented to customize
// the way an error returned by Parse is rendered
type ErrorRenderer interface {
	PrintError(err error)
}

// SeparatorRenderer can be implemented to customize
// the way a prompt is separated from another, default
// is to add a new line
type SeparatorRenderer interface {
	PrintSeparator()
}
