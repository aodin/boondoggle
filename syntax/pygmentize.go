package syntax

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

// Pygmentize uses the external "pygmentize" command. Pygmentize must
// be installed on the executing machine:
// http://pygments.org/
// pip install Pygments
type Pygmentize struct{}

var _ Highlighter = Pygmentize{}

func (h Pygmentize) Highlight(text []byte, lang string) ([]byte, error) {
	// Always output as HTML
	// TODO Add line numbers? "-O", "style=colorful,linenos=table"
	args := []string{"-f", "html"}

	// If a language was provided, add it as an arg
	if lang != "" {
		args = append(args, "-l", lang)
	}

	// Prepare the command to read from std in and write to std out
	pygmentize := exec.Command("pygmentize", args...)

	input, err := pygmentize.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("StdinPipe error: %s", err)
	}

	output, err := pygmentize.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("StdoutPipe error: %s", err)
	}

	// Begin the command
	if err = pygmentize.Start(); err != nil {
		return nil, fmt.Errorf("Start error: %s", err)
	}

	if _, err = input.Write(text); err != nil {
		return nil, fmt.Errorf("Write error: %s", err)
	}

	if err = input.Close(); err != nil {
		return nil, fmt.Errorf("Close error: %s", err)
	}

	// Read the generated HTML from stdout
	b, err := ioutil.ReadAll(output)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	// End the command
	if err = pygmentize.Wait(); err != nil {
		return nil, fmt.Errorf("Wait error: %s", err)
	}

	return b, nil
}

// TODO add a constructor that checks for presence of "pygmentize"?
