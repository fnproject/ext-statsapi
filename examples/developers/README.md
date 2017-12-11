# How to build a custom Fn server executable containing the statistics API extension

This example shows you how to build a custom Fn server executable containing the statistics API extension.

This will be of particular interest to developers working on the statistics API extension itself.

## Create `main.go`

You need to set up a go project with a copy of the [core Fn server `main.go`](https://github.com/fnproject/fn/blob/master/main.go).

Then add the following line to your `main.go` before the call to `funcServer.Start(ctx)`:
```
funcServer.AddExtensionByName(statistics.StatisticsExtensionName())
```
You will need to add the following import statement:
```
"github.com/fnproject/ext-metrics/statistics"
```
There is an [example `main.go`](../../main.go) in the root directory of this project.

You then need to install dependencies and build your executable.

## Build your custom Fn server

These steps describe how to build a custom Fn server using the example [main.go](../../main.go).

You need to work in the root directory of this project
```
cd $GOPATH/src/github.com/fnproject/ext-metrics
```

Install dependencies:

```sh
glide install
```

Build the executable:

```sh
go build
```

## Run your custom Fn server

```sh
./ext-metrics
```

By default, the metrics API will fetch data from a Prometheus server listening at `localhost:9090`. If a different host or port is required, set the following before starting your custom Fn server:
```
export FN_EXT_METRICS_PROM_HOST=<host>
export FN_EXT_METRICS_PROM_PORT=<port>
```

## Start Prometheus

Before you can use the statistics API you need to start Prometheus.

The following command starts prometheus in a Docker container, using the config file `prometheus.yml` in this directory.
Replace `<ip-address>` with the IP address on which the Fn server is listening:
```
docker run --name=prometheus -d -p 9090:9090 \
  -v ${GOPATH}/src/github.com/fnproject/ext-metrics/examples/developers/prometheus.yml:/etc/prometheus/prometheus.yml \
  --add-host="fnserver:<ip-address>" prom/prometheus
```    
On Linux you can use the following:
```
docker run --name=prometheus -d -p 9090:9090 \
  -v ${GOPATH}/src/github.com/fnproject/ext-metrics/examples/developers/prometheus.yml:/etc/prometheus/prometheus.yml \
  --add-host="fnserver:`route | grep default | awk '{print $2}'`" prom/prometheus
```

You can now deploy and run functions and try out the statistics API extension as described in the main [README](https://github.com/fnproject/ext-metrics/blob/master/README.md).
