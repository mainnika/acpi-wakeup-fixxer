package wakeup

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type procfsMock struct{}

func (p *procfsMock) ACPIWakeup() (io.ReadCloser, error) {
	const testFileRaw = `Device	S-state	  Status   Sysfs node
PEG0	  S3	*disabled
EC	  S4	*disabled  platform:PNP0C09:00
HDEF	  S3	*disabled  pci:0000:00:1b.0
RP01	  S3	*enabled   pci:0000:00:1c.0
RP02	  S3	*enabled   pci:0000:00:1c.1
RP03	  S3	*enabled   pci:0000:00:1c.2
ARPT	  S4	*enabled   pci:0000:03:00.0
RP05	  S3	*enabled   pci:0000:00:1c.4
RP06	  S3	*enabled   pci:0000:00:1c.5
SPIT	  S3	*disabled  spi:spi-APP000D:00
XHC1	  S3	*disabled  pci:0000:00:14.0
ADP1	  S4	*disabled  platform:ACPI0003:00
LID0	  S4	*disabled  platform:PNP0C0D:00
`

	return io.NopCloser(strings.NewReader(testFileRaw)), nil
}

func (p *procfsMock) ACPIWakeupWrite() (io.WriteCloser, error) {
	return nil, fmt.Errorf("not implemented")
}

func TestSomething(t *testing.T) {

	tt := []struct {
		name        string
		wantStatus  Status
		wantDevices []string
	}{
		{
			name:       "Read and parse successfully",
			wantStatus: StatusEnabled,
			wantDevices: []string{
				"RP01",
				"RP02",
				"RP03",
				"RP05",
				"RP06",
				"ARPT",
			},
		},
		{
			name:       "Read and parse disabled successfully",
			wantStatus: StatusDisabled,
			wantDevices: []string{
				"PEG0",
				"EC",
				"HDEF",
				"SPIT",
				"XHC1",
				"ADP1",
				"LID0",
			},
		},
		{
			name:       "Read and parse all successfully",
			wantStatus: StatusAll,
			wantDevices: []string{
				"RP01",
				"RP02",
				"RP03",
				"RP05",
				"RP06",
				"ARPT",
				"PEG0",
				"EC",
				"HDEF",
				"SPIT",
				"XHC1",
				"ADP1",
				"LID0",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := NewWakeupController()
			c.ProcfsProvider = &procfsMock{}

			devices, err := c.GetWakeupDevices(tc.wantStatus)
			assert.NoError(t, err)

			assert.Len(t, devices, len(tc.wantDevices))
			for i, wantDevice := range tc.wantDevices {
				assert.Contains(t, devices, wantDevice, "device %d not found", i)
			}
		})
	}
}
