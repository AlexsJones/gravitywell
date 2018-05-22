# gravitywell

![gravitywell](resources/bg.png)


Pull all your Kubernetes deployment configuration into one place.

Run one command and one manifest to switch clusters, deploy services and be the boss of your infrastructure.

![example](resources/output.gif)

## Example overview Manifest

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
                 Type: statefulset
                 Path: statefulset
        - Deployment:
            Name: "kubernetes-zookeeper-cluster"
            Namespace: "zk"
            Git: "https://github.com/AlexsJones/kubernetes-zookeeper-cluster.git"
            Action:
             - Execute:
                Shell: "ls -la"
                Kubectl:
                  Path: micro
                  Type: statefulset
                  Command: replace
````

### Roadmap

- [] Parallel cluster Deployments
- [] Rationalise back into native API for manifest parsing
- [] Expand to deploy from in-memory git repo
- [] Support additional VCS (SVN etc.)
