package plog

import (
	"fmt"
	"os"
	"os/exec"
)

// Exec returns an Opener which runs and connects to a binary that serves a Plog on stdio.
func Exec(exe string, args ...string) Opener {
	return func() (Messenger, error) {
		cmd := exec.Command(exe, args...)

		cmd.Stderr = os.Stderr

		cmdIn, err := cmd.StdinPipe()
		if err != nil {
			return nil, fmt.Errorf("stdin pipe: %w", err)
		}

		debug("made stdin pipe")

		cmdOut, err := cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("stdout pipe: %w", err)
		}

		debug("made stdout pipe")

		err = cmd.Start()
		if err != nil {
			return nil, fmt.Errorf("start: %w", err)
		}

		debug("started cmd %q", exe)

		io, err := IO(cmdOut, cmdIn)()
		if err != nil {
			return nil, fmt.Errorf("io: %w", err)
		}

		return execMes{
			Cmd: cmd,
			IO:  io,
		}, nil
	}
}

type execMes struct {
	Cmd *exec.Cmd
	IO  Messenger
}

func (e execMes) Recv() (*Msg, error) {
	return e.IO.Recv()
}

func (e execMes) Send(msg *Msg) error {
	return e.IO.Send(msg)
}

func (e execMes) Close() error {
	debug("killing %s", e.Cmd)

	err := e.Cmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("kill: %w", err)
	}

	return e.IO.Close()
}
