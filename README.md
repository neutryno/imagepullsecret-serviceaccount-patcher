# ImagePullSecret Service Account Patcher

Simple Go application that takes to the Kubernetes API to add (multiple) `ImagePullSecrets` to all 
ServiceAccounts in the cluster. 
This project was started because of the issue that credentials to private container registry cannot be
provided on a clusterwide level (cf. [stackoverflow issue](https://stackoverflow.com/questions/52320090/automatically-add-imagepullsecrets-to-a-serviceaccount)).
It was inspired by [titansoft-pte-ltd/imagepullsecret-patcher](https://github.com/titansoft-pte-ltd/imagepullsecret-patcher) which, 
however, only allows to add one private container registry secret to the cluster's service accounts.

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
kubectl apply -f deployment/deployment.yaml
kubectl apply -f deployment/rbac.yaml
```

# Development

## Build
```bash
GOOS=linux go build -o ./dist/app .
docker build . -t neutryno/serviceaccount-patcher
```