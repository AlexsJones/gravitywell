# gravitywell

![gravitywell](resources/bg.png)


Pull all your Kubernetes deployment configuration into one place.

Run one command and one manifest to switch clusters, deploy services and be the boss of your infrastructure.

## Example overview Manifest
```
APIVersion: "v1"
Strategy:
  - Cluster:
      Name: "minikube"
      Deployments:
        - Deployment:
           Name: "kubernetes-nifi-cluster"
           Git: "github.com/AlexsJones/kubernetes-nifi-cluster.git"
           Action:
            - Execute:
               shell: -|
                ./build_environment.sh default
               kubectl:
                 create: deployment
        - Deployment:
            Name: kubernetes-zookeeper-cluster
            Git: github.com/AlexsJones/kubernetes-zookeeper-cluster.git
            Action:
             - Execute:
                shell: -|
                 ./build_environment.sh default
                kubectl:
                  create: deployment

````


Cool eh?
