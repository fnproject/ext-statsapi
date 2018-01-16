# Fn server extension to provide a statistics API 

This is a Fn extension that extends the Fn API to provide some basic statistics about the operation of the Fn server. 

The API uses Prometheus internally which allows statistics to be combined from multiple Fn servers,
not just the one which is handling the request.

Since this is an extension it is not included in the core Fn server: 
to use it you need to build a custom Fn server, configured to include this extension.
This API also requires a Prometheus server to be running.

There are two examples which describe how to build a custom version of the Fn server: 

* [How to build a custom Fn server docker image](./examples/operators/README.md) (mainly of interest to end users and operators)
* [How to build a custom Fn server executable](./examples/developers/README.md) (mainly of interest to developers working on the statistics API extension itself)

These examples also describe how to start Prometheus 
and how to configure the custom Fn server and Prometheus to connect to one another.

## Try some API calls

If you have Prometheus and a custom Fn server running as described in these examples, 
you can then deploy and run some functions and obtain information about them using the statistics API extension.

### Statistics for all applications

The following API call requests metric values for the past five minutes, with an interval of 30s between values.

```sh
curl 'http://localhost:8080/v1/stats'
```

To specify a different time range and interval see [Time and step parameters](#time-and-step-parameters) below 

### Statistics for a single application

To obtain statistics for a single application `hello-async-a`:
```sh
curl 'http://localhost:8080/v1/apps/hello-async-a/stats'
```
### Statistics for a single route

To obtain statistics for a single route `hello-async-a1` in application `hello-async-a`:
```sh
curl 'http://localhost:8080/v1/apps/hello-async-a/routes/hello-async-a1/stats'
```

### Time and step parameters

The following API call requests metric values for the time period from `starttime` to `endtime`, with an interval of `step` between values. (You will need to replace the example values of `starttime` to `endtime` shown below with more recent times or you won't get any statistics.)

```sh
curl 'http://localhost:8080/v1/stats?starttime=2017-11-24T18:01:30.851Z&endtime=2017-11-24T18:11:30.849Z&step=30s'
```

`starttime` and `endtime` should be of the form `2017-11-24T18:01:30.851Z`

`step` should be a number followed by a time unit, such as `30s` or `5m`.

## Response format

Here is a sample response:

```json
{
  "status":"success",
  "data":{
    "calls":[
      {
        "time":1512416119,
        "value":12
      },
      {
        "time":1512416149,
        "value":20
      },
      {
        "time":1512416179,
        "value":34
      },
      {
        "time":1512416209,
        "value":41
      }
    ],
    "completed":[
      {
        "time":1512416119,
        "value":10
      },
      {
        "time":1512416149,
        "value":15
      },
      {
        "time":1512416179,
        "value":28
      },
      {
        "time":1512416209,
        "value":35
      }
    ],
    "errors":[
      {
        "time":1512416119,
        "value":1
      },
      {
        "time":1512416149,
        "value":3
      },
      {
        "time":1512416179,
        "value":3
      },
      {
        "time":1512416209,
        "value":3
      }
    ],   
    "timedout":[
      {
        "time":1512416119,
        "value":1
      },
      {
        "time":1512416149,
        "value":2
      },
      {
        "time":1512416179,
        "value":3
      },
      {
        "time":1512416209,
        "value":3
      }
    ],           
    "durations":[
       {
         "time":1512416119,
         "value":13.310655658088855
       },
       {
         "time":1512416149,
         "value":14.841753632280154
       },
       {
         "time":1512416179,
         "value":34.04480194936709
       },
       {
         "time":1512416209,
         "value":45.51680103676471
       }
     ],
    "failed":[]
  }
}
```

The `success` element will be set to `success` if the API call is successful. 
If the API call is unsuccessful then the `success` element will be set to `error` and an additional element `error` will contains a description of the failure.

The `data` element contains elements `completed`, `errors`, `timedout` and `durations`. 

* The `calls` element is an array of objects. Each object contains a single observation of the `fn_calls` counter metric at a specific time. This is a count of the number of function calls made since the server was started.

* The `completed` element is an array of objects. Each object contains a single observation of the `fn_completed` counter metric at a specific time. This is a count of the number of successful function calls since the server was started.

* The `errors` element is an array of objects. Each object contains a single observation of the `fn_errors` counter metric at a specific time.
This is a count of failed (but not timed out) function calls since the server was started.
If there were no errors the array may be empty.  

* The `timedout` element is an array of objects. Each object contains a single observation of the `fn_timedout` counter metric at a specific time.
This is a count of timed out function calls since the server was started.
If no function calls timed out function the array may be empty.  

* The `durations` element is an array of objects. Each object contains a single calculated value of the rolling mean `fn_span_agent_submit_duration_seconds` histogram metric, where the rolling mean is calculated over a period of one minute. 

In addition the `data` element contains the element `failed`. This is included for backward compatibility. It is deprecated and will be removed in the future.

* The `failed` element is an array of objects. Each object contains a single observation of the `fn_failed` counter metric at a specific time.
This is a count of failed or timed out function calls since the server was started.
If there were no failures the array may be empty.  




