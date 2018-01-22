# How to build a custom Fn server docker image containing the statistics API extension

This example shows you how to build a custom Fn server docker image containing the statistics API extension.

This will be of particular interest to end users and operators.

It describes two alternative ways to run the custom Fn server docker image:

* [Run your custom Fn image and Prometheus using Docker Compose](/examples/operators/README.md#run-your-custom-fn-image-and-prometheus-using-docker-compose)
* [Run your custom Fn image and Prometheus separately](/examples/operators/README.md#run-your-custom-fn-image-and-prometheus-separately)

## Configure `ext.yaml`

You need just one file, `ext.yaml`, in which you must list the extensions to be included in your custom Fn server image. 
This directory contains an example [ext.yaml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators/ext.yaml) configured to include a single extension, the statistics API

```yaml
extensions:
- name: github.com/fnproject/ext-statsapi/stats
```

If you require additional extensions, add a `name` element for each one.


## Build your custom image

We will use the example [ext.yaml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators/ext.yaml) in this directory:
```sh
cd $GOPATH/src/github.com/fnproject/ext-statsapi/examples/operators
```

To build your custom image:
```sh
fn build-server -t imageuser/fn-ext-statsapi
```
If you intend to deploy the image to a docker image repository you will need to change `imageuser` to something suitable such as your repository username. If you are not planning to do this you can leave it unchanged.

## Run your custom Fn image and Prometheus using Docker Compose

The quickest way to start your custom Fn server and Prometheus is to use Docker Compose. 
This takes care of configuring the two processes to connect to each other.

Install Docker Compose using [these instructions](https://docs.docker.com/compose/install/). 

We will use the [docker-compose.yml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators/docker-compose.yml) in this directory.
You should change `imageuser` to whatever you specified when building your custom Fn image.

```yaml
version: '3'
services:
  fnserver:
    image: imageuser/fn-ext-statsapi
    ports:
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - $PWD/data:/app/data
    environment:
    - FN_EXT_STATS_PROM_HOST=prometheus
  prometheus:
    image: prom/prometheus
    restart: always
    ports:
      - "9090:9090"
    volumes:
      - ${GOPATH}/src/github.com/fnproject/ext-statsapi/examples/operators/prometheus.yml:/etc/prometheus/prometheus.yml
```

This starts your custom Fn image. The environment variable `FN_EXT_STATS_PROM_HOST` is used to specify that the Fn server should fetch
statistics from a Prometheus server running on the `prometheus:9090`, where   `prometheus` is defined in 
[docker-compose.yml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators/docker-compose.yml)
to refer to the Prometheus server.

It also starts Prometheus using the config file [prometheus.yml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators/prometheus.yml) 
which configures Prometheus to scrape metrics from a Fn server running on `fnserver:8080`, where `fnserver` is defined in
[docker-compose.yml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators/docker-compose.yml)
to refer to the Fn server.

Now start your custom Fn image and Prometheus

```sh
docker-compose up
```

You can now deploy and run functions and try out the statistics API extension as described in the main [README](https://github.com/fnproject/ext-statsapi/blob/master/README.md).

## Run your custom Fn image and Prometheus separately

Alternatively you can start your custom Fn image and Prometheus separately. 

### Run your custom image

The following command is used to run your custom image. Replace `<ip-address>` with the IP address on which the Fn server is listening:

```sh
docker run --rm --name fnserver -d -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $PWD/data:/app/data -p 8080:8080 \
  -e FN_EXT_STATS_PROM_HOST=<ip-address> imageuser/imagename
```

* `FN_EXT_STATS_PROM_HOST` is an environment variable which specifies the host on which the Prometheus server is running. 
The default is `localhost`, which doesn't work if the Fn server is running in docker .
* You can also use `FN_EXT_STATS_PROM_PORT` to specify the port.

On Linux you can use
```sh
docker run --rm --name fnserver -d -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $PWD/data:/app/data -p 8080:8080 \
  -e FN_EXT_STATS_PROM_HOST=`route | grep default | awk '{print $2}'` imageuser/imagename
```
### Start Prometheus

The following command starts prometheus in a Docker container, using the config file `prometheus.yml` in this directory.
Replace `<ip-address>` with the IP address on which the Fn server is listening:
```
docker run --name=prometheus -d -p 9090:9090 \
  -v ${GOPATH}/src/github.com/fnproject/ext-statsapi/examples/developers/prometheus.yml:/etc/prometheus/prometheus.yml \
  --add-host="fnserver:<ip-address>" prom/prometheus
```    
On Linux you can use the following:
```
docker run --name=prometheus -d -p 9090:9090 \
  -v ${GOPATH}/src/github.com/fnproject/ext-statsapi/examples/developers/prometheus.yml:/etc/prometheus/prometheus.yml \
  --add-host="fnserver:`route | grep default | awk '{print $2}'`" prom/prometheus
```

You can now deploy and run functions and try out the statistics API extension as described in the main [README](https://github.com/fnproject/ext-statsapi/blob/master/README.md).