# File Distribution Operator

## Motivation

This proposal is to introduce a new operator to the Kubernetes API. This operator will be responsible for distributing files to nodes. It does this by mounting a Secret to a defined path on the node.

The operator will be configurable by a CustomResource which will define the following:

- `secretName` - The name of the secret where the files are stored
- `fileName` - The name of the field in the secret where the file is stored (if not specified, all files in the secret will be mounted)
- `destination` - The destination path on the node where the files will be mounted. This path will be relative to the root of the node's filesystem.
- `interval` - The interval at which the files will be checked for updates. This interval will be in seconds.
- `mode` - The mode of the file(s) to be distributed. 

For security reasons, only CRs and secrets in the same namespace that the operator is deployed to will be allowed to be processed.

For better understanding, here is a sample architecture of the operator:

![Architecture](fdo_architecture.png)

## Implementation

The operator will be implemented as a Kubernetes controller. It will be written in Go and will use the Operator SDK. The operator itself should be distributed either as a yaml to be applied to a cluster or as a helm package. 

For the sake of convenience, a CI pipeline will be created to build and push the operator to a public registry and run linting and unit tests.

## Tasks and Milestones

- [ ] Create a new operator that can be deployed to a cluster.
- [ ] Create a new CustomResourceDefinition that can be used to configure the operator. 
- [ ] Setup GitHub Action CI 
- [ ] Create a new helm chart that can be used to deploy the operator to a cluster.
- [ ] Create documentation how to install and use the operator.

## Responsibilities

- Krasniqi A. (@waodim) - helm chart
- Krenn G. (@gkrenn) - operator code
- Pepryk K. (@kpe09) - GitHub CI and documentation/samples

Please keep in mind that this is a proposal and is subject to change, especially if there are any issues with the implementation responsibilities have to be reassigned.

## Use Cases

- Distributing SSH keys to nodes.
- Distributing SSL certificates to nodes.
- Distributing configuration files to nodes.
- Ensuring AppArmor profiles are up to date on nodes.

