# Metrics

This is a Fn extension to extend the Fn API to provide statistical metrics. 

There are two alternative ways to use it
* Add the metrics API to to your own custom Fn server
* Use the example extended Fn server

In both cases you need to have a Prometheus server running as well.

## Add the metrics API to to your custom Fn server

If you intend to add the metrics API extension to your own custom version of the Fn server, add the following to your `main.go`:

```
// Add the metrics API extension before you call funcServer.Start(ctx)
handlers.AddEndpoints(funcServer)
```
You will need to add an import statement for `github.com/fnproject/ext-metrics/statistics`. See `main.go` in this directory for an example.

You also need to start a Prometheus server.

## Use the example extended Fn server

This project provides an example of a Fn server which has been extended to include the metrics API. 
See `main.go` in this directory.
You can use this to try out the metrics API.

### Build the example extended Fn server

If you intend to use the example extended Fn server, you first need to  build it:
```sh
glide install
```

```sh
go build
```
### Start the example extended Fn server 

By default, the metrics API will fetch data from a Prometheus server listening at `localhost:9090`. If a different host or port is required, set
```
export FN_EXT_METRICS_PROM_HOST=<host>
export FN_EXT_METRICS_PROM_PORT=<port>
```
Now start the example extended Fn server 

```sh
./ext-metrics
```

## Start Prometheus

You need to configure Prometheus to scrape data from the Fn server. 

A suitable Prometheus configuration file `prometheus.yml` is provided in this project. This contains the following configuration settings:

* Metrics are scraped from a Fn server listening on `fnserver:8080`. `fnserver` is an alias. When Prometheus is started in docker the `--add-host` parameter is used to specify the actual hostname or IP address to which it corresponds.
* Some Prometheus metrics are relabelled to support the queries that are required.

Start Prometheus, replacing `<ip-address>` with the IP address on which the extended Fn server is listening:
```
  docker run --name=prometheus -d -p 9090:9090 \
    -v ${GOPATH}/src/github.com/fnproject/fn/examples/grafana/prometheus.yml:/etc/prometheus/prometheus.yml \
    --add-host="fnserver:<ip-address>" prom/prometheus
```    
On Linux you can do:
```
  docker run --name=prometheus -d -p 9090:9090 \
    -v ${GOPATH}/src/github.com/fnproject/ext-metrics/prometheus.yml:/etc/prometheus/prometheus.yml \
    --add-host="fnserver:`route | grep default | awk '{print $2}'`" prom/prometheus
```

## Try some API calls

After starting your custom Fn server or the example extended Fn server, you must create some functions. If you don't then you won't get any statistics!

### Statistics for all applications

The following API call requests metric values for the past five minutes, with an interval of 30s between values.

```sh
curl 'http://localhost:8080/v1/statistics'
```

To specify a different time range and interval see [Time and step parameters](#time-and-step-parameters) below 

### Statistics for a single application

To obtain statistics for a single application `hello-async-a`:
```sh
curl 'http://localhost:8080/v1/apps/hello-async-a/statistics'
```
### Statistics for a single route

To obtain statistics for a single route `hello-async-a1` in application `hello-async-a`:
```sh
curl 'http://localhost:8080/v1/apps/hello-async-a/routes/hello-async-a1/statistics'
```

### Time and step parameters

The following API call requests metric values for the time period from `starttime` to `endtime`, with an interval of `step` between values. (You will need to replace the example values of `starttime` to `endtime` shown below with more recent times or you won't get any statistics.)

```sh
curl 'http://localhost:8080/v1/statistics?starttime=2017-11-24T18:01:30.851Z&endtime=2017-11-24T18:11:30.849Z&step=30s'
```

`starttime` and `endtime` should be of the form `2017-11-24T18:01:30.851Z`

`step` should be a number followed by a time unit, such as `30s` or `5m`.

## Response format

Here is a sample response:

```json
{
  "status": "success",
  "data": {
    "completed": [
      {
        "Time": 1511546850,
        "Value": 18
      },
      {
        "Time": 1511546880,
        "Value": 32
      },
      {
        "Time": 1511546910,
        "Value": 48
      },
      {
        "Time": 1511546940,
        "Value": 67
      },
      {
        "Time": 1511546970,
        "Value": 89
      },
      {
        "Time": 1511547000,
        "Value": 104
      },
      {
        "Time": 1511547030,
        "Value": 121
      },
      {
        "Time": 1511547060,
        "Value": 139
      }
    ],
    "failed": []
  }
}
```

The `success` element will be set to success` if the API call is successful. 
If the API call is unsuccessful then `success` element will be set to `error` and an additional element `error` will contains a description of the failure.

The `data` element contains elements `completed` and `failed`. 

* The `completed` element contains an array of objects. Each object contains a single observation of the `fn_api_completed` counter metric at a specific time. This is a count of the number of successful function calls since the server was started.
* The `failed` element is an array of objects. Each object contains a single observation of the `fn_api_failed` counter metric at a specific time.
This is a count of failed (or timed out) function calls since the server was started.
If there were no failures the array may be empty.  

## Design notes

* We use the element names `completed` and `failed` for consistency with the existing Prometheus metrics from which they are obtained. 

## Still to be done

* The `completed` metric will be replaced a new metric `calls` which will be a count of all completed calls, including failed calls.

* ~Currently it is assumed that the Prometheus server is on `localhost:9090`. This needs to be configurable.~

* ~Per-application metrics~

* ~Per-function (route) metrics~

* Proper tests

* Replace `completed` and `failed` with new metrics that conform to the requirements

* Duration metrics. This will require a new tracing span that works for both cold and hot functions. 

## Running the tests

To run the tests
```sh
go test ./...
```
