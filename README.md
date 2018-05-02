# Podstalk

Container that displays informations about the Pod it's running on.

## Usage

```sh
kubectl create namespace podstalk && \
kubectl -n podstalk apply -f https://github.com/mhutter/podstalk/raw/master/kube/daemonset.yml
```

Will create a daemonset running `mhutter/podstalk` and a service named
`podstalk`, along with a Role and a Rolebinding required to list pods in the
namespace.

### Customization

The following env vars can be set to customize Podstalk:

* `TITLE` (_default: `Podstalk`_) - `<title>` of the HTML page

## Development

To regenerate bindata.go, do:

```sh
go-bindata -pkg podstalk -prefix data/ data/*
```

For local development I recommend using [gin][]:

```sh
gin --build cmd/podstalk
```

This will run podstalk on http://localhost:3000/ and reload each time a go file
is changed.

To use with [minikube][] you can do the following:

```sh
# delploy podstalk into minikube
kubectl create namespace podstalk
kubectl -n podstalk apply -f kube/deployment.yml
kubectl -n podstalk expose svc/podstalk --type NodePort --name podstalk-ext

# set `imagePullPolicy` to `IfNotPresent`
kubectl -n podstalk edit deployment podstalk

# Connect to the minikube docker daemon
eval $(minikube docker-env)

# build the image (each time you changed the code)
docker build -t mhutter/podstalk .

# restart all pods (each time you built a new image)
kubectl -n podstalk delete pod -l app=podstalk

# access podstalk in the browser
minikube -n podstalk service podstalk-ext
```

[gin]: https://github.com/codegangsta/gin
[minikube]: https://kubernetes.io/docs/getting-started-guides/minikube/
