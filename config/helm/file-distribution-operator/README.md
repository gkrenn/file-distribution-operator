# File Distribution Operator Helm Chart

## Prerequisites

In order to use this helm chart it goes without saying that the helm cli is needed.
Details on how to install can be found [here](https://github.com/helm/helm)
Helm is very useful for packaging Kubernetes resources and makes the installation of Kubernetes applications easy and straight forward.

## Installation

In order to apply a specific configuration one has to fill out the values in the provided ```values.yaml```.

There are several ways to use this chart. One can for example render the whole manifest for the Deployment of the operator and apply it by hand.

```sh
helm template config/helm/file-distribution-operator/ > manifests.yaml
```

This manifest file can then be applied to the Kubernetes cluster and this will lead to the operator to be deployed.

```sh
kubectl create namespace fdo-system
kubectl apply -f manifests.yaml
```

Another, more convenient way to use this chart is to use the ```install``` command, where one can choose a name for the release:

```sh
helm install <name of the release> config/helm/file-distribution-operator/ --namespace fdo-system --create-namespace
```

## Uninstallation

If the helm chart was applied through the rendered manifest the uninstallation looks as follows:

```sh
kubectl delete -f manifests.yaml
```

If the helm chart was installed through the ``ìnstall`` command the uninstallation looks as follows:

```shell
helm uninstall <name of the release> --namespace fdo-system
```
