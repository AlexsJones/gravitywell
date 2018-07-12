package configuration

//Options ...
type Options struct {
	VCS         string
	TempVCSPath string
	APIVersion  string
	SSHKeyPath  string
	IgnoreList  []string
	DryRun      bool
	TryUpdate   bool
	Redeploy    bool
}
