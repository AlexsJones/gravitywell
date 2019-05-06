package kinds

//- Execute:
//Kind: "shell"
//Configuration:
//Command: pwd
//Path: ../ #Optional value
//- Execute:
//Kind: "shell"
//Configuration:
//Command: ./build_environment.sh default
//- Execute:
//Kind: "kubernetes"
//Configuration:
//Path: deployment #Optional value
//AwaitDeployment: true #Optional defaults to false
type Execute struct {
	Kind          string            `yaml:"Kind"`
	Configuration map[string]string `yaml:"Configuration"`
}

type ActionList struct {
	APIVersion string    `yaml:"APIVersion"`
	Kind       string    `yaml:"Kind"`
	Executions []Execute `yaml:"Executions"`
	LocalPath  string    `yaml:"LocalPath"`
	RemotePath string    `yaml:"RemotePath"`
}
