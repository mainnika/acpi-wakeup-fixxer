package procfs

import (
	"fmt"
	"io"
	"os"
)

const (
	acpiWakeupDefaultPath = "/proc/acpi/wakeup"
)

type Procfs interface {
	ACPIWakeup() (io.ReadCloser, error)
	ACPIWakeupWrite() (io.WriteCloser, error)
}

type ProcfsDefaultPath struct{}

func (p *ProcfsDefaultPath) ACPIWakeup() (io.ReadCloser, error) {
	f, err := os.Open(acpiWakeupDefaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", acpiWakeupDefaultPath, err)
	}

	return f, nil
}

func (p *ProcfsDefaultPath) ACPIWakeupWrite() (io.WriteCloser, error) {
	f, err := os.OpenFile(acpiWakeupDefaultPath, os.O_WRONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open for writing %s: %w", acpiWakeupDefaultPath, err)
	}

	return f, nil
}
