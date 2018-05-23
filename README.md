# gravitywell

![gravitywell](resources/bg.png)


Pull all your Kubernetes deployment configuration into one place.

Run one command and one manifest to switch clusters, deploy services and be the boss of your infrastructure.

_It's a bit like docker-compose for Kubernetes deployments across clusters!_

![example](resources/output.gif)

## Requirements

`go get github.com/AlexsJones/vortex`

## Example overview Manifest

_Please see examples directory_

Example command: `gravitywell -config deploy-nifi.yaml`

Lets look at the deploy-nifi.yaml...

```
APIVersion: "v1"
Strategy:
  - Cluster:
      Name: "minikube"
      Deployments:
        - Deployment:
           Name: "kubernetes-nifi-cluster"
           Namespace: "nifi"
           Git: "https://github.com/AlexsJones/kubernetes-nifi-cluster.git"
           Action:
            - Execute:
               Shell: "ls -la"
               Kubectl:
                 Command: replace
                 Path: statefulset
        - Deployment:
            Name: "kubernetes-zookeeper-cluster"
            Namespace: "zk"
            Git: "https://github.com/AlexsJones/kubernetes-zookeeper-cluster.git"
            Action:
             - Execute:
                Shell: "./build_environment.sh small"
                Kubectl:
                  Path: deployment
                  Command: replace
````

### Support APIResource types

- [x] ConfigMap
- [x] StatefulSet
- [x] Deployment
- [x] Service
- [ ] CronJob
- [ ] PersistantVolume
- [ ] PersistantVolumeClaim

### Roadmap

- [ ] Parallel cluster Deployments
- [x] Rationalise back into native API for manifest parsing
- [ ] Expand to deploy from in-memory git repo
- [ ] Support additional VCS (SVN etc.)
