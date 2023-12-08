# Starship

Using starship for a better development experience, where one can create a 
complex testing and qa environment in a multi-node, multichain environment.

## Installation
Run `make setup-deps` to install and verify the dependencies installed.
Checkout docs: https://starship.cosmology.tech/get-started/step-1

### Connect to a k8s cluster
#### For linux
Run `make setup-kind`, to spin up a local k8s cluster using kind, which uses docker to spin up
For more details, checkout: https://starship.cosmology.tech/get-started/step-2#211-setup-with-kind-cluster

#### For mac
Enable `kubernetes` from docker-desktop
For more details, checkout: https://starship.cosmology.tech/get-started/step-2#212-setup-with-docker-desktop

#### Remote k8s cluster
If you have a remote k8s cluster, then you can use that as well, directly.
All you would need is to be able to access it via `kubectl`.

### Check connection to k8s cluster
Run `kubectl get nodes`, to check the connection to the k8s cluster

## Deploy Starship
We use helm-charts for packaging all k8s based setup. This is controlled by the helm chart versions. Run the following to fetch the helm chart from the Makefile variable HELM_VERSION.
```bash
# Fetch the helm chart with correct version
make setup-helm
```

Now you can spin up Starship locally with:
```bash
# deploy configs/local.yaml
make install-local
# OR
# deploy configs/devnet.yaml if you have more resources available, or a remote k8s cluster
make install-devnet
```

> Note: In the CI we use the config file `configs/ci.yaml` to deploy the starship, with limited resources

### Check the status of the deployment
Run `kubectl get pods` to check the status of the pods, once the pods are in `Running` state, you can proceed with port-forwarding.
```bash
make port-forward-local
# OR
make port-forward-devnet
```

Now you can access the starship explorer UI at http://localhost:8080

## Run tests
Tests are designed such that one can re-run the same tests against an already running infra. This will save the cost of initialization of the infra.
Startup the cluster
```bash
make install

## check status of the pods
kubectl get pods

## Once the pods are up run:
make port-forward

## Run tests, can run this now multiple times as long as tests are running
make test

## Cleanup
make stop
```

## Teardown
To teardown the starship deployment, run:
```bash
make stop
```

## Troubleshooting local setup

Currently, there seems to be some issues when running starship on a local system. This section will help clear out some of them

### Not starting
If all or some of the pods are in `Pending` state, then it means that resources for docker containers are not enough.
There can be 2 ways around this:

1. Increase the resources for your local kubernetes cluster.
* Docker Desktop: Go to `Settings` > `Resources`, increase CPU and memory
2. Reduce the resources for each of the nodes in `configs/local.yaml` file. You can look at `configs/ci.yaml` to understand the `resource` directive in the chains
* `configs/ci.yaml` uses very little resources, so should be able to run locally

> NOTE: When resoureces are reduced or if the devnet has been running for a longer time, then the pods seem to die out or keep restarting. This is due to memory overflow. Will be fixed soon. For now
> one can just run
```bash
make delete

## wait for nodes to die out, check with
kubectl get pods

## restart
make install
```

## Future work
* Add more e2e tests specific to comdex modules
* Add chain upgrade tests, as part of a seperate CI/CD process
