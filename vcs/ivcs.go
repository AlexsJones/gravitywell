package vcs

//IVCS is a VCS abstraction
type IVCS interface {
	Fetch(localpath string, remote string, reference string, keypath string) (string, error)
	Update(localpath string) (string, error)
}

//Fetch from a remote VCS
func Fetch(v IVCS, localpath string, remote string, reference string, keypath string) (string, error) {
	return v.Fetch(localpath, remote, reference, keypath)
}

//Update a locally checked out repo
func Update(v IVCS, localpath string) (string, error) {
	return v.Update(localpath)
}
