# acpi-wakeup-fixxer

This repository contains a tool to disable the wakeup ability on your macbook system that runs linux, which can be particularly useful to prevent immediate wakeup from suspend mode. This tool is compatible with various systems, including the MacBook Air 2013/2015.

# features

- disables wakeup by echoing a device name to /proc/acpi/wakeup
- default list of devices to disable is LID0 XHC1
- one can adjust the list by using `-disable` flag
- systemd unit runs the app on load and resume

# installation

## rhel/fedora systems

- take latest artifacts from https://github.com/mainnika/acpi-wakeup-fixxer/actions/workflows/build-rpm.yaml
- unzip and install `fedora-40-x86_64/result/golang-code-tokarch-mainnika-acpi-wakeup-fixxer-0-0.1.fc40.x86_64.rpm`

## build from sources
- `go build -o acpi-wakeup-fixxer ./cmd/acpi-wakeup-fixxer`

# systemd

```
[Unit]
Description=Fix ACPI wakeup issues
DefaultDependencies=no
Conflicts=shutdown.target
Before=sysinit.target shutdown.target
ConditionPathExists=|/proc/acpi/wakeup

[Service]
ExecStart=/usr/local/bin/acpi-wakeup-fixxer $ACPI_WAKEUP_FIXXER_ARGS
Type=oneshot
TimeoutSec=0
RemainAfterExit=yes
EnvironmentFile=-/etc/default/acpi-wakeup-fixxer

[Install]
WantedBy=sysinit.target
```