# Notice: Not Actively Maintained
This project is no longer actively maintained.
No further updates, bug fixes, or support are planned.
Please consider this when deciding to use or contribute to this project.
Feel free to fork & use this code as you see fit.

# ImagePullSecret Service Account Patcher

Simple Go application that takes to the Kubernetes API to add (multiple) `ImagePullSecrets` to all 
ServiceAccounts in the cluster. 

## Motivation
This project was started because of the issue that credentials to private container registry cannot be
provided on a clusterwide level (cf. [stackoverflow issue](https://stackoverflow.com/questions/52320090/automatically-add-imagepullsecrets-to-a-serviceaccount)).
Others suggested manually pulling images to your nodes, patching Service Accounts manually or adapting the `docker/config.json`
of each cluster's node (cf. [here](https://stackoverflow.com/a/55230340/5930295)).

This project was inspired by [titansoft-pte-ltd/imagepullsecret-patcher](https://github.com/titansoft-pte-ltd/imagepullsecret-patcher) 
which, however, only allows to add one private container registry secret to the cluster's service accounts.

## Usage
It is at best used in conjunction with [mittwald/kubernetes-replicator](https://github.com/mittwald/kubernetes-replicator).
Thus this is the complete approach:

1. Install [mittwald/kubernetes-replicator](https://github.com/mittwald/kubernetes-replicator)
2. Create container registry secrets in the `kube-system` namespace
```bash
kubectl -n kube-system create secret docker-registry <SECRET_NAME_1> --docker-server=<registry.server.de> --docker-username=<username> --docker-password=<password>
kubectl -n kube-system create secret docker-registry <SECRET_NAME_2> --docker-server=<registry.server.de> --docker-username=<username> --docker-password=<password>
```
3. Patch secrets to make them replicable by [mittwald/kubernetes-replicator](https://github.com/mittwald/kubernetes-replicator)
```bash
kubectl -n kube-system patch secret <SECRET_NAME_1> -p '{"metadata": {"annotations": {"replicator.v1.mittwald.de/replicate-to": ".*"}}}'
kubectl -n kube-system patch secret <SECRET_NAME_2> -p '{"metadata": {"annotations": {"replicator.v1.mittwald.de/replicate-to": ".*"}}}'
```
4. Add your secrets' names to the `REGISTRY_SECRET_NAMES` environment variable in `deployment/deployment.yaml`. 
5. Install neutryno/serviceaccount-patcher
```bash
kubectl apply -f https://raw.githubusercontent.com/neutryno/imagepullsecret-serviceaccount-patcher/master/deployment/rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/neutryno/imagepullsecret-serviceaccount-patcher/master/deployment/deployment.yaml
```

## Build
```bash
docker buildx build . -t neutryno/imagepullsecret-serviceaccount-patcher --platform linux/amd64,linux/arm64 --no-cache
```
