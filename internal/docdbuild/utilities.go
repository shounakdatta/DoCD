package docdbuild

import (
	"os"
	"os/exec"
)

// cmdReference : Storage struct for cmd instances
type cmdReference struct {
	Cmd        *exec.Cmd
	LogFile    *os.File
	Terminated bool
}
