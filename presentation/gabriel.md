# What does our operator do?

<img src="/fdo_architecture_old.png" class="" />


---

# What does our operator do?

<img src="/fdo_architecture.png" class="" />

---

# Use Cases

- Distributing SSH keys to nodes.
- Distributing SSL certificates to nodes.
- Distributing configuration files to nodes.
- Ensuring AppArmor profiles are up to date on nodes.

--- 

# operator sdk

- uses official kubernetes controller-runtime library
- makes writing of operators easier by providing
    - high level APIs and abstractions
    - Tools for scaffolding and code generation
    - Extensions to cover common Operator use cases

<br>
<br>
<br>

``` bash
mkdir memcached-operator
cd memcached-operator
operator-sdk init --domain example.com --repo github.com/example/memcached-operator
operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller
```
