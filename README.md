# Gravitywell

<img src="https://github.com/ashleymcnamara/gophers/blob/master/SPACEGIRL_GOPHER.png?raw=true" alt="drawing" width="200"/>

[![Maintainability](https://api.codeclimate.com/v1/badges/a6cd570642c5aeeedaf9/maintainability)](https://codeclimate.com/github/AlexsJones/gravitywell/maintainability)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Release](https://img.shields.io/github/release/AlexsJones/gravitywell.svg)

**Update** AWS API now in Alpha for GCP/AWS cluster builds

Gravitywell is designed to create kubernetes clusters and deploy your applications.
It uses YAML to store deployments and supports multiple versions of kubernetes resource definitions.
It lets you store your entire container infrastructure as code.

## How is gravitywell different and why is it useful?

- Can deploy across multiple cloud providers at the same time.
    - _All other tools we found would either do the multiple cloud deployment, without the apps or visa versa._
- Based on cloud providers own container API
    - _Ain't nobody got time to be deploying custom networking or policies when its done for free._
- Uses dynamic interpolation with [vortex](https://github.com/AlexsJones/vortex) so you aren't writing template files for days.
- Allows you to do more than just deploy clusters; it lets you bootstrap them with fully working dependant services.
    - _Think about getting mongodb,zookeeper,consul,nifi,nginx,api's and a bunch of other stuff going straight away_

## Getting Started

To get started you'll need golang installed or to fetch the binary from homebrew/releases page (OSX)

- Get with golang: 
    - `go get github.com/AlexsJones/gravitywell`
- Download with homebrew: `brew tap AlexsJones/homebrew-gravitywell && brew install gravitywell`
- Download as a cross-platform release: `https://github.com/AlexsJones/gravitywell/releases`

### Prerequisites

The current implementation works with Google Cloud Platform & Amazon web services.

_For Google Cloud Platform please set your service account for the right project_

`export GOOGLE_APPLICATION_CREDENTIALS=~/Downloads/alex-example-e28058e8985b.json`

_For Amazon Web Services please set the aws profile name and region_

`export AWS_DEFAULT_PROFILE=alexprod`
`export AWS_DEFAULT_REGION=us-west-2`

Aws also requires additional tools for authentication found [here](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html)


### Running

At this point you are ready to run gravitywell.

For working with templates as per the examples you'll also need [vortex](`go get github.com/AlexsJones/gravitywell`)
_This can be installed either via golang or as a binary also_


_Lets take it for a spin using the gcp example_

```
#If you've looked at the templates you'll see a helmesque style of interpolation
# "gke_{{.projectname}}_{{.projectregion}}_{{.clustername}}" we're going to override

vortex --output example-gcp/deployment --template example-gcp/templates \
--set "projectname=alex-example" --set "projectregion=us-east4" --set "clustername=testcluster"

# Now an examples/templates folder exists you simple run...

gravitywell create -f examples-gcp/deployment

# This will now start to provision any clusters that are required and deploy applications

```

### Example files

This is what an example cluster may look like:

```bash
APIVersion: "v1"
Kind: "Cluster"
Strategy:
  - Provider:
      Name: "Google Cloud Platform"
      Clusters:
        - Cluster:
            FullName: "gke_{{.projectname}}_{{.projectregion}}_{{.clustername}}"
            ShortName: "{{.clustername}}"
            Project: "{{.projectname}}"
            Region: "us-east4"
            Zones: ["us-east4-a"]
            Labels:
              type: "test"
            InitialNodeCount: 1
            InitialNodeType: "n1-standard-1"
            OauthScopes: "https://www.googleapis.com/auth/monitoring.write,
          https://www.googleapis.com/auth/logging.write,
          https://www.googleapis.com/auth/trace.append,
          https://www.googleapis.com/auth/devstorage.full_control,
          https://www.googleapis.com/auth/compute"
            NodePools:
              - NodePool:
                  Name: "np1"
                  Count: 3
                  NodeType: "n1-standard-1"
                  Labels:
                    k8s-node-type: "test"
            PostInstallHook:
              - Execute:
                  Path: "."
                  Shell: "gcloud container clusters get-credentials {{.clustername}} --region={{.projectregion}} --project={{.projectname}}"
            PostDeleteHook:
              - Execute:
                  Path: "."
                  Shell: "pwd"
```

And this is an example application

```bash
APIVersion: "v1"
Kind: "Application"
Strategy:
  - Cluster:
      FullName: "gke_{{.projectname}}_{{.projectregion}}_{{.clustername}}"
      ShortName: "{{.clustername}}"
      Applications:
        - Application:
            Name: "kubernetes-apache-tika"
            Namespace: "tika"
            Git: "git@github.com:AlexsJones/kubernetes-apache-tika.git"
            Action:
              - Execute:
                  Kind: "shell"
                  Configuration:
                    Command: pwd
                    Path: ../ #Optional value
              - Execute:
                  Kind: "shell"
                  Configuration:
                    Command: ./build_environment.sh default
              - Execute:
                  Kind: "kubernetes"
                  Configuration:
                    Path: deployment #Optional value
                    AwaitDeployment: true #Optional defaults to false
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
* kubicorn does alot of very cool stuff
* https://eksctl.io/ was a fantastic reference for golang AWS sdk API