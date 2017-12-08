# How to build a custom Fn server docker image containing the statistics API extension

This example shows you how to build a custom Fn server docker image containing the statistics API extension.

This will be of particular interest to end users and operators.

## Configure `ext.yaml`

You need just one file, `ext.yaml`, in which you must list the extensions you wish to include in your custom Fn server image. 
This directory contains a `ext.yaml` configured to include a single extension, the statistics API

If you require additional extensions, add a `name` element for each one.
The following would include the statistics API extension and a fictitious extension whose package name is `github.com/fnproject/ext-foo/foo`.


```yaml
extensions:
- name: github.com/fnproject/ext-metrics/statistics
- name: github.com/fnproject/ext-foo/foo
```

## Build your custom image

To build your custom image:


```sh
fn build-server -t imageuser/imagename
```

You can then use standard docker tools to install it in a docker image repository if required.


## Run your custom image


To run your custom image:

```sh
docker run --rm --name fnserver -it -v /var/run/docker.sock:/var/run/docker.sock -v $PWD/data:/app/data -p 8080:8080 -e FN_EXT_METRICS_PROM_HOST=prometheus --link prometheus imageuser/imagename
```

## Start Prometheus

Before you can use the statistics API you need to start Prometheus.

Now start Prometheus, specifying this config file `prometheus.yml` in this directory:
```
  docker run --name=prometheus -d -p 9090:9090 \
    -v ${GOPATH}/src/github.com/fnproject/fn/ext-metrics/examples/operators/prometheus.yml:/etc/prometheus/prometheus.yml \
    --link fnserver prom/prometheus
```
`prometheus.yml` configures Prometheus to scrape metrics from a Fn server running on `fnserver:8080`, where `fnserver` is an alias that is set in the command  above to refer to a container named `fnserver` in which the Fn server is expected to be running.

