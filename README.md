# Podstalk

Containers that display information about the pods in the namespace they're running in.


## Usage

0. Clone this repository
0. `helm upgrade -i podstalk ./chart`


### Configuration

Podstalk can be configured using the following environment variables:

* `NAMESPACE` - (_default_: whatever is in `/var/run/secrets/kubernetes.io/serviceaccount/namespace`) - Namespace to list pods in. Set to `""` (empty string) to watch pods in ALL namespaces (given the ServiceAccount Podstalk runs under is privileged enough).
* `DEBUG` - Set to a non-empty value to enable additional logging output.
* `BASE_PATH` - If the app is not served from the host's root (`/`), set this to the subdir it is served from (eg. `/podstalk`).

If you use the Helm chart, see `values.yaml`

## Development

For out-of-cluster development I use the following setup:

* [minikube][] - local K8s clusters
* [gin][] - A live reload utility for Go web applications
* [yarn][] - Node.js package manager

```sh
# Start minikube
minikube -p podstalk start

# This should now point to your minikube cluster
kubectl config get-contexts

# Start the Server
make dev-server

# In a separate window, start the client dev server
make dev-client
```

## License

MIT (see `LICENSE`)

---
> [Manuel Hutter](https://hutter.io/) -
> GitHub [@mhutter](https://github.com/mhutter) -
> Twitter [@dratir](https://twitter.com/dratir)

[minikube]: https://github.com/kubernetes/minikube
[gin]: https://github.com/codegangsta/gin
[yarn]: https://yarnpkg.com/
