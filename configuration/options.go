package configuration

//Options ...
type Options struct {
	VCS         string
	TempVCSPath string
	APIVersion  string
	SSHKeyPath  string
	Parallel    bool
	DryRun      bool
	TryUpdate   bool
	Redeploy    bool
}
