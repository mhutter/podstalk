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
