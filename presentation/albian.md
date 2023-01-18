# Helm

- Package manager for K8s
- Automates maintenance of YAML manifests
- Has special directory
---

# Structure of helm charts


<img src="/helm_dir.png" class="" />

---

# Values file

- plain YAML files
- values inside the file can be accessed as an attribute <br>
<br>
<img src="/val_access.png" class="" />
<br>
<br>
<img src="/values.png" class="" />

---

# Demo

``` bash
helm template config/helm/file-distribution-operator/ > manifests.yaml

helm install <name of the release> config/helm/file-distribution-operator/ --namespace fdo-system --create-namespace
```