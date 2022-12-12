# install-plan-approver-operator
A Kubernetes Operator to approve InstallPlans created by Operator Lifecycle Manager (OLM).

## Description
Operator Lifecycle Manager (OLM) provides manual upgrade strategy for Subscriptions to avoid installation of a specific version with bugs, having potential to break things on cluster. An InstallPlan is created every time there is a new version available for a specific installed operator. Those InstallPlans require manual approver after logging into the cluster, which is against the core concept of GitOps. Install-Plan-Approver-Operator is designed to remove the need to login to cluster to approve the InstallPlan. It checks for the StartingCSV field in Subscription spec and if there is an InstallPlan available against that specific ClusterServiceVersion (CSV), then it approves that InstallPlan.

## Getting Started
You’ll need a Kubernetes cluster, with Operator Lifecycle Manager (OLM) installed, to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
Clone the repo and run

```sh
make deploy
```

This'll create all the required resources to

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

