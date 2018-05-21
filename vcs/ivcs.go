package vcs

//IVCS is a VCS abstraction
type IVCS interface {
	Fetch(localpath string, remote string) (string, error)
}

//Fetch from a remote VCS
func Fetch(v IVCS, localpath string, remote string) (string, error) {
	return v.Fetch(localpath, remote)
}
