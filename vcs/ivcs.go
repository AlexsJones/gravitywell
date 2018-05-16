package vcs

//IVCS is a VCS abstraction
type IVCS interface {
	Fetch(remote string) (string, error)
}

//Fetch from a remote VCS
func Fetch(v IVCS, remote string) (string, error) {
	return v.Fetch(remote)
}
