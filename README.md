# gravitywell

![gravitywell](resources/bg.png)

_ITS LIKE HELM MEETS TERRAFORM_
- Deploy cluster into GCP from yaml
- Deploy manifests into those clusters from yaml

Pull all your Kubernetes deployment configuration into one place.

Run one command and one manifest to switch clusters, deploy services and be the boss of your infrastructure.

It's as easy as `gravitywell create -f ./`

![provision](resources/provision.png)

Or `gravitywell delete -f ./`

![deprovision](resources/deprovision.png)

## Installation

`go get github.com/AlexsJones/gravitywell`

## Requirements

- Golang 1.10

- `go get github.com/AlexsJones/vortex`
- `export GOOGLE_APPLICATION_CREDENTIALS=somegooglecloudserviceaccountfile.json` (_This needs to be set to a valid service account for the project you with to perform GCP operations on_)


## Example overview Manifest

There are two kinds of manifest.
`Cluster` and `Application`

```
APIVersion: "v1"
Kind: "Cluster"
Strategy:
  - Provider:
      Name: "Google Cloud Platform"
      Clusters:
        - Cluster:
            Name: "testclustera"
            Project: "beamery-trials"
            Region: "us-east4"
            Zones: ["us-east4-a"]
            InitialNodeCount: 1
            InitialNodeType: "n1-standard-1"
            OauthScopes: "https://www.googleapis.com/auth/monitoring.write,
            https://www.googleapis.com/auth/logging.write,
            https://www.googleapis.com/auth/trace.append,
            https://www.googleapis.com/auth/devstorage.full_control,
            https://www.googleapis.com/auth/compute"
            NodePools:
              - NodePool:
                  Name: "Pool-A"
                  Count: 3
                  NodeType: "n1-standard-1"
            PostInstallHook:
              - Execute:
                  Shell: "gcloud container clusters get-credentials TestClusterA --region=us-east4 --zone=a"
```

```
APIVersion: "v1"
Kind: "Application"
Strategy:
  - Cluster:
      Name: "minikube"
      Applications:
        - Application:
           Name: "kubernetes-nifi-cluster"
           Namespace: "nifi"
           Git: "git@github.com:AlexsJones/kubernetes-nifi-cluster.git"
           Action:
            - Execute:
               Shell: "ls -la"
               Kubectl:
                 Path: statefulset
        - Application:
            Name: "kubernetes-zookeeper-cluster"
            Namespace: "zk"
            Git: "git@github.com:AlexsJones/kubernetes-zookeeper-cluster.git"
            Action:
             - Execute:
                Shell: "./build_environment.sh small"
                Kubectl:
                  Path: deployment

```
Command output ...

```
go run main.go create -f examples/application/small.yaml
2019/01/02 12:19:08 Loading examples/application/small.yaml
Application kind found
WARN[0000] Switching to cluster: gke_beamery-trials_us-east4_testclustera
DEBU[0000] Loading deployment kubernetes-apache-tika
DEBU[0000] Fetching deployment kubernetes-apache-tika into .gravitywell/kubernetes-apache-tika
Enumerating objects: 45, done.
Total 45 (delta 0), reused 0 (delta 0), pack-reused 45
WARN[0001] Running shell command ./build_environment.sh default
Building for environment default
DEBU[0001] Successful
ERRO[0001] Could not read from file %s.gravitywell/kubernetes-apache-tika/deployment
ERRO[0001] Could not read from file %s.gravitywell/kubernetes-apache-tika/deployment/tika
INFO[0001] Decoded Kind: extensions/v1beta1, Kind=Deployment
INFO[0001] Decoded Kind: /v1, Kind=Namespace
INFO[0001] Decoded Kind: /v1, Kind=Service
Found Namespace resource
DEBU[0001] Namespace deployed
DEBU[0001] Found deployment resource
DEBU[0001] Deployment deployed
DEBU[0001] Found service resource
DEBU[0002] Service deployed
```

```
## Commands

_We support three kubectl commands currently_

```
delete
apply
create
```

### Supported Cloud providers

- [x] Google cloud platform 
- [ ] Amazon Web Services

### Supported Kubernetes resource definition types

- [x] ConfigMap
- [x] StatefulSet
- [x] Deployment
- [x] Service
- [x] PodDisruptionBudget
- [x] ServiceAccounts
- [x] RoleBinding
- [x] CronJob
- [ ] PersistantVolume
- [ ] PersistantVolumeClaim

### Roadmap

- [ ] Depends on Cluster from Application flag