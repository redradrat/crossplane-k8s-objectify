# Crossplane K8s Objectify

Small tool to convert kubernetes resources manifests into a bunch of `Object` resources to use within Crossplane `Compositions`.

See: [Crossplane - Kubernetes Provider](https://github.com/crossplane-contrib/provider-kubernetes)

## Usage
```
converter -i manifests.yaml -o objects.yaml
```
