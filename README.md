# Golang Operator

Go (golang) operator for Kubernetes.

## Installation

Here's what you need to do if you just want to install this operator.

### Build and run the operator

Before running the operator, the CRD must be registered with the Kubernetes apiserver:

```sh
kubectl apply -f deploy/crds/golang_v1alpha1_golang_crd.yaml
```

```sh
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/operator.yaml
```

Verify that the `golang-operator` is up and running:

```sh
> kubectl get deploy
NAME              DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
golang-operator   1         1         1            1           1m
```

### Create a Golang CR

Create the example `Golang` CR that was generated at `deploy/crds/golang_v1alpha1_golang_cr.yaml`:

```sh
kubectl apply -f deploy/crd/golang_v1alpha1_golang_cr.yaml
```

Verify that the `golang-operator` creates the deployment for the CR:

```sh
> kubectl get deploy
NAME              DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
go-hello-world    1         1         1            1           3m
golang-operator   1         1         1            1           4m
```

Check the pods and CR status to confirm the status is updated with the `golang` pod names:

```sh
> kubectl get pods
NAME                               READY     STATUS    RESTARTS   AGE
go-hello-world-7c67d65476-9l6lk    1/1       Running   0          5m
golang-operator-5f4c5c675b-j7ff7   1/1       Running   1          6m
```

#### Customing the Golang CR

You can copy or customise the example CR to suit your own application.

| Custom Property      |      Description                          | Required |
|----------------------|:-----------------------------------------:|---------:|
| `spec.size`          | The number of pods you wish to be created | Yes      |
| `spec.image`         | Docker image of your Go application       | Yes      |
| `spec.containerPort` | Port that your Docker application runs on | No       |

## Development setup

Install the Operator SDK CLI:

```sh
mkdir -p $GOPATH/src/github.com/operator-framework
cd $GOPATH/src/github.com/operator-framework
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
make dep
make install
```

Refer to the [Operator SDK User Guide][operator-sdk-user-guide] for instructions on how to use and develop this project.

## Contributing

1. Fork it (<https://github.com/craicoverflow/golang-operator/fork>)
2. Create your feature branch (`git checkout -b fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin fooBar`)
5. Create a new Pull Request

<!-- Markdown link & img dfn's -->
[npm-image]: https://img.shields.io/npm/v/datadog-metrics.svg?style=flat-square
[npm-url]: https://npmjs.org/package/datadog-metrics
[npm-downloads]: https://img.shields.io/npm/dm/datadog-metrics.svg?style=flat-square
[travis-image]: https://img.shields.io/travis/dbader/node-datadog-metrics/master.svg?style=flat-square
[travis-url]: https://travis-ci.org/dbader/node-datadog-metrics
[wiki]: https://github.com/yourname/yourproject/wiki
[operator-sdk-user-guide]: https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md