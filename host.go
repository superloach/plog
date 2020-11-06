package plug

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Host makes a Plug which runs and connects to the stdin/stdout of a binary which serves a Guest.
func Host(exe string, args ...string) *Plug {
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

		// swap out and in
		p.dec = json.NewDecoder(cmdOut)
		p.enc = json.NewEncoder(cmdIn)

		p.closeFn = func() {
			debug("killing %q", exe)
			cmd.Process.Kill()
		}

		close(p.ioReady)

		return nil
	}

	return p
}
