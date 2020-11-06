package plog

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Exec makes a Plog which runs and connects to a binary running an IO Plog.
func Exec(exe string, args ...string) *Plog {
	debug("host %q %v", exe, args)

	p := empty()

	p.openFn = func() error {
		cmd := exec.Command(exe, args...)

		cmd.Stderr = os.Stderr

		cmdIn, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("stdin pipe: %w", err)
		}

		debug("made stdin pipe")

		cmdOut, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("stdout pipe: %w", err)
		}

		debug("made stdout pipe")

		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("start: %w", err)
		}

		debug("started cmd %q", exe)

		p.closeFn = func() {
			debug("killing %q", exe)
			cmd.Process.Kill()
		}

		p.mes = ioMessenger{
			Decoder: json.NewDecoder(cmdOut),
			Encoder: json.NewEncoder(cmdIn),
		}
		close(p.mesReady)

		return nil
	}

	return p
}
