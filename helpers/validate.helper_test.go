package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidVersionForm(t *testing.T) {
	type args struct {
		version string
		field   string
	}
	tests := []struct {
		name        string
		args        args
		status      bool
		WantMessage bool
	}{
		{
			name:        "Valid version",
			args:        args{version: "1.0.0"},
			status:      true,
			WantMessage: false,
		},
		{
			name:        "Empty string",
			args:        args{version: ""},
			status:      false,
			WantMessage: true,
		},
		{
			name:        "Wrong form",
			args:        args{version: "1.2"},
			status:      false,
			WantMessage: true,
		},
		{
			name:        "Non integer version",
			args:        args{version: "1.0.x"},
			status:      false,
			WantMessage: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, message := IsValidVersionFormat(&tt.args.version, tt.args.field)
			assert.Equal(t, tt.status, status)
			assert.Equal(t, tt.WantMessage, message != nil)
		})
	}
}

func TestIsValidVersionUpdate(t *testing.T) {
	type args struct {
		updatingVersion string
		currentVersion  string
	}
	tests := []struct {
		name      string
		args      args
		WantError bool
	}{
		{
			name: "Not updating",
			args: args{
				currentVersion:  "1.0.0",
				updatingVersion: "1.0.0",
			},
			WantError: true,
		},
		{
			name: "Patching",
			args: args{
				currentVersion:  "1.0.0",
				updatingVersion: "1.0.1",
			},
			WantError: false,
		},
		{
			name: "Patching with more than one version",
			args: args{
				currentVersion:  "1.0.1",
				updatingVersion: "1.0.4",
			},
			WantError: true,
		},
		{
			name: "Minor update",
			args: args{
				currentVersion:  "1.1.5",
				updatingVersion: "1.2.0",
			},
			WantError: false,
		},
		{
			name: "Minor update without rounding up patching",
			args: args{
				currentVersion:  "1.1.5",
				updatingVersion: "1.2.5",
			},
			WantError: true,
		},
		{
			name: "Minor with more than one version",
			args: args{
				currentVersion:  "1.1.5",
				updatingVersion: "1.3.0",
			},
			WantError: true,
		},
		{
			name: "Major update",
			args: args{
				currentVersion:  "1.1.5",
				updatingVersion: "2.0.0",
			},
			WantError: false,
		},
		{
			name: "Major update without rounding up patching",
			args: args{
				currentVersion:  "1.1.5",
				updatingVersion: "2.0.5",
			},
			WantError: true,
		},
		{
			name: "Major update without rounding up minor",
			args: args{
				currentVersion:  "1.1.5",
				updatingVersion: "2.1.0",
			},
			WantError: true,
		},
		{
			name: "Major update without rounding up",
			args: args{
				currentVersion:  "1.1.5",
				updatingVersion: "2.1.0",
			},
			WantError: true,
		},
		{
			name: "Major with more than one version",
			args: args{
				currentVersion:  "1.1.5",
				updatingVersion: "3.0.0",
			},
			WantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, message := IsValidVersionUpdate(&tt.args.updatingVersion, tt.args.currentVersion, "")
			assert.Equal(t, tt.WantError, message != nil)
		})
	}
}
