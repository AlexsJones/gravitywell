package configuration

import "time"

//Options ...
type Options struct {
	VCS                string
	TempVCSPath        string
	APIVersion         string
	SSHKeyPath         string
	DryRun             bool
	MaxBackOffDuration time.Duration
}
