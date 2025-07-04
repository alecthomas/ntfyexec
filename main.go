package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
)

var cli struct {
	Title string `help:"Ntfy notification title, if any." env:"NTFY_TITLE"`
	Token string `help:"Ntfy access token, if any." env:"NTFY_TOKEN"`

	Topic   string   `arg:"" required:"" help:"Ntfy topic."`
	Command []string `arg:"" required:"" help:"Command to execute." passthrough:"partial"`
}

const description = `
Execute a command, and send a notification to ntfy.sh if it fails.

The combined stdout and stderr will be sent as the notification body.
`

func main() {
	kctx := kong.Parse(&cli, kong.Description(description))

	err := execute()
	kctx.FatalIfErrorf(err)
}

func execute() error {
	w, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(w.Name())

	cmd := exec.Command(cli.Command[0], cli.Command[1:]...)
	cmd.Stdout = w
	cmd.Stderr = w
	err = cmd.Run()
	if err == nil {
		return nil
	}

	// Construct notification title.
	title := cli.Title
	if title == "" {
		title = fmt.Sprintf("%q failed", cli.Command[0])
	}

	if eerr := (&exec.ExitError{}); errors.As(err, &eerr) {
		title += fmt.Sprintf(" (exit status %d)", eerr.ExitCode())
	}

	_, err = w.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek temp file: %w", err)
	}
	return notify(title, w)
}

func notify(title string, output *os.File) error {
	log.Printf("Notifying: %s", title)
	req, _ := http.NewRequest("PUT", "https://ntfy.sh/"+cli.Topic, output)
	req.Header.Set("Authorization", "Bearer "+cli.Token)
	req.Header.Set("Title", title)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return fmt.Errorf("notification failed: %s", body)
	}
	return err
}
