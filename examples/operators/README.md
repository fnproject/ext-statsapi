# How to build a custom Fn server docker image containing the statistics API extension

This example shows you how to build a custom Fn server docker image containing the statistics API extension.

This will be of particular interest to end users and operators.

## Configure `ext.yaml`

You need just one file, `ext.yaml`, in which you must list the extensions to be included in your custom Fn server image. 
This directory contains an example [ext.yaml](https://github.com/fnproject/ext-metrics/blob/master/examples/operators/ext.yaml) configured to include a single extension, the statistics API

```yaml
extensions:
- name: github.com/fnproject/ext-metrics/statistics
```

If you require additional extensions, add a `name` element for each one.


## Build your custom image

We will use the example [ext.yaml](https://github.com/fnproject/ext-metrics/blob/master/examples/operators/ext.yaml) in this directory:
```sh
cd $GOPATH/src/github.com/fnproject/ext-metrics/examples/operators
```

To build your custom image:
```sh
fn build-server -t imageuser/imagename
```

You can then use standard docker tools to deploy it in a docker image repository if required.

## Run your custom Fn image and Prometheus using Docker Compose

The quickest way to start your custom Fn server and Prometheus is to use Docker Compose. 
This takes care of configuring the two processes to connect to each other.

Install Docker Compose using [these instructions](https://docs.docker.com/compose/install/). 

We will use the [docker-compose.yml](https://github.com/fnproject/ext-metrics/blob/master/examples/operators/docker-compose.yml) in this directory.
You should change `imageuser/imagename` to whatever you specified when building your custom Fn image.

```yaml
version: '3'
services:
  fnserver:
    image: imageuser/imagename
    ports:
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - $PWD/data:/app/data
    environment:
    - FN_EXT_METRICS_PROM_HOST=prometheus
  prometheus:
    image: prom/prometheus
    restart: always
    ports:
      - "9090:9090"
    volumes:
      - ${GOPATH}/src/github.com/fnproject/ext-metrics/examples/operators/prometheus.yml:/etc/prometheus/prometheus.yml
```

Now start your custom Fn image and Prometheus

```
cd $GOPATH/src/github.com/fnproject/ext-metrics/examples/operators
```

```sh
docker-compose up
```

You can now deploy and run functions and try out the statistics API extension as described in the main [README](https://github.com/fnproject/ext-metrics/blob/master/README.md).

## Run your custom Fn image and Prometheus separately

Alternatively you can start your custom Fn image and Prometheus separately. 

### Run your custom image

The following command is used to run your custom image. Replace `<ip-address>` with the IP address on which the Fn server is listening:

```sh
docker run --rm --name fnserver -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $PWD/data:/app/data -p 8080:8080 \
  -e FN_EXT_METRICS_PROM_HOST=<ip-address> imageuser/imagename
```

* `FN_EXT_METRICS_PROM_HOST` is an environment variable which specifies the host on which the Prometheus server is running. 
The default is `localhost`, which doesn't work if the Fn server is running in docker .
* You can also use `FN_EXT_METRICS_PROM_PORT` to specify the port.

On Linux you can use
```sh
docker run --rm --name fnserver -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $PWD/data:/app/data -p 8080:8080 \
  -e FN_EXT_METRICS_PROM_HOST=`route | grep default | awk '{print $2}'` imageuser/imagename
```

### Start Prometheus

Before you can use the statistics API you need to start Prometheus.

Now start Prometheus, specifying the example Prometheus config file [prometheus.yml](https://github.com/fnproject/ext-metrics/blob/master/examples/operators/prometheus.yml) in this directory:
```
  docker run --name=prometheus -d -p 9090:9090 \
    -v ${GOPATH}/src/github.com/fnproject/ext-metrics/examples/operators/prometheus.yml:/etc/prometheus/prometheus.yml \
    --link fnserver prom/prometheus
```
[prometheus.yml](https://github.com/fnproject/ext-metrics/blob/master/examples/operators/prometheus.yml) configures Prometheus to scrape metrics from a Fn server running on `fnserver:8080`, where `fnserver` is an alias that is set in the command  above to refer to a container named `fnserver` in which the Fn server is expected to be running.

You can now deploy and run functions and try out the statistics API extension as described in the main [README](https://github.com/fnproject/ext-metrics/blob/master/README.md).

