# Gravitywell

[![Maintainability](https://api.codeclimate.com/v1/badges/a6cd570642c5aeeedaf9/maintainability)](https://codeclimate.com/github/AlexsJones/gravitywell/maintainability)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Release](https://img.shields.io/github/release/AlexsJones/gravitywell.svg)

Gravitywell is designed to deploy kubernetes clusters and their applications.
It uses YAML to store deployments and supports multiple versions of kubernetes resource definitions.
Operating on a simple concurrency model it allows you to deploy in parallel to save time.
## Getting Started

To get started you'll need golang installed or to fetch the binary from homebrew/releases page (OSX)

- Get with golang: 
    - `go get github.com/AlexsJones/gravitywell`
- Download with homebrew: `brew tap AlexsJones/homebrew-gravitywell && brew install gravitywell`
- Download as a cross-platform release: `https://github.com/AlexsJones/gravitywell/releases`

### Prerequisites

The current implementation works with Google Cloud Platform.
This means you'll need a service account with at least Kubernetes cluster admin scopes.

- See more here:`https://cloud.google.com/iam/docs/creating-managing-service-accounts`

For working with templates as per the examples you'll also need [vortex](`go get github.com/AlexsJones/gravitywell`)
_This can be installed either via golang or as a binary also_

### Installing

Once you have the service account you'll want to export the path:

`export GOOGLE_APPLICATION_CREDENTIALS=~/Downloads/alex-example-e28058e8985b.json`


At this point you are ready to run gravitywell.
Remember to have [vortex](`go get github.com/AlexsJones/gravitywell`) installed if you want to avoid hand writing manifests.

_Lets take it for a spin_

```
#If you've looked at the templates you'll see a helmesque style of interpolation
# "gke_{{.projectname}}_{{.projectregion}}_{{.clustername}}" we're going to override

vortex --output examples/deployment --template examples/templates \
--set "projectname=alex-example" --set "projectregion=us-east4" --set "clustername=testcluster"

# Now an examples/templates folder exists you simple run...

gravitywell create -f examples/deployment

# This will now start to provision any clusters that are required and deploy applications

```

### Flags

```go
DryRun     bool   `short:"d" long:"dryrun" description:"Performs a dryrun."`
FileName   string `short:"f" long:"filename" description:"filename to execute, also accepts a path."`
SSHKeyPath string `short:"s" long:"sshkeypath" description:"Custom ssh key path."`
MaxTimeout string `short:"m" long:"maxtimeout" description:"Max rollout time e.g. 60s or 1m"`
```

## Running the tests

`go test ./... -v` from the gravitywell directory on your gopath


## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/AlexsJones/gravitywell/tags). 

## Authors

* **Alex Jones** - *Initial work* 

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Helm & terraform both great projects
