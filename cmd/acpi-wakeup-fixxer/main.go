package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"code.tokarch.uk/mainnika/acpi-wakeup-fixxer/pkg/wakeup"
)

var (
	devicesToDisableFlag []string
)

var rootCmd = &cobra.Command{
	Use:   "acpi-wakeup-fixxer",
	Short: "A CLI tool to fix ACPI wakeup issues",
	Args:  cobra.NoArgs,
	Long:  "acpi-wakeup-fixxer is a command-line tool to fix ACPI wakeup issues on your system.",
	Run:   rootCmdRun,
}

func init() {
	rootCmd.PersistentFlags().StringArrayVar(&devicesToDisableFlag, "disable", []string{"LID0", "XHC1"}, "Devices to disable")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	wakeUpController := wakeup.NewWakeupController()

	devices, err := wakeUpController.GetWakeupDevices(wakeup.StatusEnabled)
	if err != nil {
		slog.Error("failed to get wakeup devices", "error", err)
		os.Exit(1)
	}

	slog.Info("Found devices", "devices", devices)
	slog.Info("Disabling devices", "devices", devicesToDisableFlag)

	var disabledDevices []string
	for _, device := range devices {
		needDisable := false
		for _, deviceToDisable := range devicesToDisableFlag {
			if device == deviceToDisable {
				needDisable = true
				break
			}
		}

		if needDisable {
			slog.Info("Disabling device", "device", device)
			if err := wakeUpController.ToggleWakeupDevice(device); err != nil {
				slog.Error("failed to disable device", "device", device, "error", err)
				continue
			}

			disabledDevices = append(disabledDevices, device)
		}
	}

	if len(disabledDevices) > 0 {
		slog.Info("Disabled devices", "devices", disabledDevices)
	} else {
		slog.Info("No devices were disabled")
	}
}
