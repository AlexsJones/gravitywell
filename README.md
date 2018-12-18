# gravitywell

![gravitywell](resources/bg.png)

- Deploy cluster into GCP from yaml
- Deploy manifests into those clusters from yaml

Pull all your Kubernetes deployment configuration into one place.

Run one command and one manifest to switch clusters, deploy services and be the boss of your infrastructure.

It's as easy as `gravitywell create -f ./`

## Installation

`go get github.com/AlexsJones/gravitywell`

## Requirements

`go get github.com/AlexsJones/vortex`
`export GOOGLE_APPLICATION_CREDENTIALS=` This needs to be set to a valid service account for the project you with to perform GCP operations on
## Example overview Manifest

There are two kinds of manifest.
`Cluster` and `Application`

```bash
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

```bash
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