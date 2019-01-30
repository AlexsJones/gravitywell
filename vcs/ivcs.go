package vcs

//IVCS is a VCS abstraction
type IVCS interface {
	Fetch(localpath string, remote string, keypath string, branch string) (string, error)
	Add(localpath string, remote string, keypath string, files []string) error
	Commit(localpath string, remote string, keypath string, message string) error
	Push(localpath string, remote string, keypath string) error
	Update(localpath string) (string, error)
}

//Fetch from a remote VCS
func Fetch(v IVCS, localpath string, remote string, keypath string, branch string) (string, error) {
	return v.Fetch(localpath, remote, keypath, branch)
}

//Add from a remote VCS
func Add(v IVCS, localpath string, remote string, keypath string, files []string) error {
	return v.Add(localpath, remote, keypath, files)
}

//Commit from a remote VCS
func Commit(v IVCS, localpath string, remote string, keypath string, message string) error {
	return v.Commit(localpath, remote, keypath, message)
}

//Push from a remote VCS
func Push(v IVCS, localpath string, remote string, keypath string) error {
	return v.Push(localpath, remote, keypath)
}

//Update a locally checked out repo
func Update(v IVCS, localpath string) (string, error) {
	return v.Update(localpath)
}
