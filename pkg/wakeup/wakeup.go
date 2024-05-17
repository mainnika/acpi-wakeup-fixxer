package wakeup

import (
	"bufio"
	"fmt"
	"log/slog"
	"strings"

	"code.tokarch.uk/mainnika/acpi-wakeup-fixxer/pkg/procfs"
)

const (
	deviceHeader      = "Device"
	columnDeviceIndex = 0
	columnStatusIndex = 2
	columnIndexMax    = 3
)

type Status string

const StatusAll Status = ""
const StatusEnabled Status = "*enabled"
const StatusDisabled Status = "*disabled"

type WakeupController struct {
	ProcfsProvider procfs.Procfs
}

func NewWakeupController() *WakeupController {
	return &WakeupController{ProcfsProvider: &procfs.ProcfsDefaultPath{}}
}

func (w *WakeupController) GetWakeupDevices(withStatus Status) ([]string, error) {
	wakeupFile, err := w.ProcfsProvider.ACPIWakeup()
	if err != nil {
		return nil, fmt.Errorf("failed to get wakeup file: %w", err)
	}
	defer func() { _ = wakeupFile.Close() }()

	scanner := bufio.NewScanner(wakeupFile)
	var devices []string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < columnIndexMax {
			slog.Warn("unexpected number of columns", "line", line)
			continue
		}

		device := fields[columnDeviceIndex]
		if device == deviceHeader {
			continue
		}

		status := Status(fields[columnStatusIndex])
		if withStatus != StatusAll && status != withStatus {
			continue
		}

		devices = append(devices, device)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan wakeup file: %w", err)
	}

	return devices, nil
}

func (w *WakeupController) ToggleWakeupDevice(device string) error {
	wakeupFile, err := w.ProcfsProvider.ACPIWakeupWrite()
	if err != nil {
		return fmt.Errorf("failed to open wakeup file for writing: %w", err)
	}

	_, err = fmt.Fprintf(wakeupFile, "%s\n", device)
	if err != nil {
		_ = wakeupFile.Close()
		return fmt.Errorf("failed to write to wakeup file: %w", err)
	}

	err = wakeupFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close wakeup file: %w", err)
	}

	return nil
}
