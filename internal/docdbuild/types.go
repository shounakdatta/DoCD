package docdbuild

import (
	"os"
	"os/exec"
)

// CmdReference : Storage struct for cmd instances
type CmdReference struct {
	Cmd     *exec.Cmd
	LogFile *os.File
}
