# gravitywell

![gravitywell](resources/bg.png)


Pull all your Kubernetes deployment configuration into one place.

Run one command and one manifest to switch clusters, deploy services and be the boss of your infrastructure.

_It's a bit like docker-compose for Kubernetes deployments across clusters!_

![example](resources/output.gif)

_Or using --dryrun to test your deployment status_

![example2](resources/output2.gif)

## Installation

`go get github.com/AlexsJones/gravitywell`

## Requirements

`go get github.com/AlexsJones/vortex`

## Example overview Manifest

_Please see examples directory_

Example command: `gravitywell -config examples/small`

_Parallel deployments with --parallel flag_

Lets look at the small.yaml...

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

_We support three kubectl commands currently_

```
replace
apply
create
```

### Command Options

```
  -config string
    	Configuration path
  -dryrun bool
    	Run a dry run deployment to test what is deployment
  -tryupdate bool
    	Try to update the resource if possible
```

### Support APIResource types

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

- [x] Rationalise back into native API for manifest parsing
- [ ] Expand to deploy from in-memory git repo


### Example of a real production configuration across multiple clusters

```
APIVersion: "v1"
Strategy:
#alpha
  - Cluster:
      Name: "gke_MYCOMPANY-alpha_us-central1-a_MYCOMPANY-alpha-cluster"
      Deployments:
        - Deployment:
           Name: "MYCOMPANY-alpha"
           Namespace: "alpha"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/ingress/hostmap.yaml
            - Execute:
                Kubectl:
                  Command: apply
                  Path: alpha/account/serviceaccount.yaml
            - Execute:
                Shell: "vortex --varpath environments/alpha.yaml --template templates/accounts/devops.yaml --output alpha/account"
                Kubectl:
                  Command: replace
                  Path: alpha/account/devops.yaml
  - Cluster:
      Name: "gke_MYCOMPANY-alpha_us-central1-a_MYCOMPANY-services"
      Deployments:
        - Deployment:
           Name: "MYCOMPANY-queue"
           Namespace: "queue"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/ingress/hostmap.yaml
            - Execute:
                Shell: "vortex --varpath environments/alpha.yaml --template templates/accounts/devops.yaml --output alpha/account"
                Kubectl:
                  Command: replace
                  Path: alpha/account/devops.yaml
        - Deployment:
           Name: "MYCOMPANY-cron"
           Namespace: "service-cron"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/ingress/hostmap.yaml
        - Deployment:
           Name: "MYCOMPANY-devops"
           Namespace: "devops"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/ingress/hostmap.yaml
  - Cluster:
      Name: "gke_MYCOMPANY-alpha_us-central1-a_MYCOMPANY-foundation"
      Deployments:
        - Deployment:
           Name: "monstache"
           Namespace: "monstache"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: alpha/service

#beta
  - Cluster:
      Name: "gke_MYCOMPANY-beta_us-central1-a_MYCOMPANY-preproduction"
      Deployments:
        - Deployment:
           Name: "MYCOMPANY-beta"
           Namespace: "beta"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/ingress/hostmap.yaml
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/account/serviceaccount.yaml
            - Execute:
                Shell: "vortex --varpath environments/beta.yaml --template templates/accounts/devops.yaml --output beta/account"
                Kubectl:
                  Command: replace
                  Path: beta/account/devops.yaml
  - Cluster:
      Name: "gke_MYCOMPANY-beta_us-central1-a_MYCOMPANY-services"
      Deployments:
        - Deployment:
           Name: "MYCOMPANY-queue"
           Namespace: "queue"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/ingress/hostmap.yaml
            - Execute:
                Shell: "vortex --varpath environments/beta.yaml --template templates/accounts/devops.yaml --output beta/account"
                Kubectl:
                  Command: apply
                  Path: beta/account/devops.yaml
        - Deployment:
           Name: "MYCOMPANY-cron"
           Namespace: "service-cron"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/ingress/hostmap.yaml
        - Deployment:
           Name: "MYCOMPANY-devops"
           Namespace: "devops"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/ingress/hostmap.yaml
  - Cluster:
      Name: "gke_MYCOMPANY-beta_us-central1-a_MYCOMPANY-foundation"
      Deployments:
        - Deployment:
           Name: "monstache"
           Namespace: "monstache"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: beta/service
  - Cluster:
      Name: "gke_MYCOMPANY-production_us-east4-a_MYCOMPANY-production"
      Deployments:
        - Deployment:
           Name: "MYCOMPANY-production"
           Namespace: "production"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/ingress/hostmap.yaml
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/account/serviceaccount.yaml
  - Cluster:
      Name: "gke_MYCOMPANY-production_us-east4-a_MYCOMPANY-service-cluster"
      Deployments:
        - Deployment:
           Name: "MYCOMPANY-queue"
           Namespace: "queue"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/ingress/hostmap.yaml
        - Deployment:
           Name: "MYCOMPANY-cron"
           Namespace: "service-cron"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/ingress/hostmap.yaml
        - Deployment:
           Name: "MYCOMPANY-devops"
           Namespace: "devops"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/service
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/ingress/hostmap.yaml
  - Cluster:
      Name: "gke_MYCOMPANY-production_us-east4-a_MYCOMPANY-foundation"
      Deployments:
        - Deployment:
           Name: "monstache"
           Namespace: "monstache"
           Git: "git@github.com:MYCOMPANY/devops-kubernetes-configuration.git"
           Action:
            - Execute:
               Kubectl:
                 Command: apply
                 Path: production/service

```
