# How to use the statistics API extension with a cluster of custom Fn servers

This example demonstrates the use of the statistics API to obtain aggregated statistics from a cluster of Fn servers.

This will be of particular interest to end users and operators.

This example uses Docker Compose.

## Build your custom image

If you have not already done so you need to build a custom Fn server docker image containing the statistics API extension. 
This is described in  [How to build a custom Fn server docker image](../operators/README.md).
For convenience the same instructions are repeated below:

```sh
cd $GOPATH/src/github.com/fnproject/ext-statsapi/examples/operators
fn build-server -t imageuser/imagename
```
If you intend to deploy the image to a docker image repository you will need to change `imageuser` to something suitable such as your repository username. If you are not planning to do this you can leave it unchanged.

## Run a cluster of two custom Fn images and Prometheus using Docker Compose

The quickest way to start a cluster of custom Fn servers and Prometheus is to use Docker Compose. 
This takes care of configuring the various processes to connect to each other.

Install Docker Compose using [these instructions](https://docs.docker.com/compose/install/). 

We will use the [docker-compose.yml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators-clustered/docker-compose.yml) in this directory.
You should change `imageuser` to whatever you specified when building your custom Fn image:

```yaml
version: '3'
services:
  logstore:
    hostname: logstore
    image: minio/minio
    ports:
      - "9091:9000"
    environment:
      - MINIO_ACCESS_KEY=admin
      - MINIO_SECRET_KEY=password
    volumes:
      - ./data/logstore:/data
    command: server /data
  db:
    image: "mysql"
    restart: always
    ports:
      - "3306:3306"
    environment:
      - "MYSQL_DATABASE=funcs"
      - "MYSQL_ROOT_PASSWORD=root"
    volumes:
      - ./data/mysql:/var/lib/mysql
  mq:
    image: "redis"
    restart: always
    ports:
      - "6379:6379"
  fnserver0:
    image: imageuser/fn-ext-statsapi
    restart: always
    depends_on:
      - mq
      - db
    ports:
      - "8080:8080"
    environment:
      FN_PORT: "8080"
      FN_EXT_STATS_PROM_HOST: "prometheus" 
      FN_DB_URL: "mysql://root:root@tcp(db:3306)/funcs"
      FN_MQ_URL: "redis://mq:6379/"
      FN_LOGSTORE_URL: "s3://admin:password@logstore:9000/us-east-1/fnlogs"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
  fnserver1:
    image: imageuser/fn-ext-statsapi
    restart: always
    depends_on:
      - mq
      - db
    ports:
      - "8081:8081"
    environment:
      FN_PORT: "8081"
      FN_EXT_STATS_PROM_HOST: "prometheus" 
      FN_DB_URL: "mysql://root:root@tcp(db:3306)/funcs"
      FN_MQ_URL: "redis://mq:6379/"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
  grafana:
    image: grafana/grafana
    restart: always
    ports:
      - "3000:3000"
    links:
      - fnserver0
      - fnserver1
      - prometheus
    depends_on:
      - fnserver0
      - fnserver1
      - prometheus
  prometheus:
    image: prom/prometheus
    restart: always
    depends_on:
      - fnserver0
      - fnserver1
    ports:
      - "9090:9090"
    links:
      - fnserver0
      - fnserver1
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
```

This will start two Fn servers using the custom Fn image you created above. 
One listens on port 8080 and the other listens on port 8081.
For both Fn server, the environment variable `FN_EXT_STATS_PROM_HOST` is used to specify that the Fn server should fetch
statistics from a Prometheus server running on the `prometheus:9090`, where   `prometheus` is defined to refer to the Prometheus server.

It will also start Prometheus using the config file [prometheus.yml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators-clustered/prometheus.yml):
```
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: 'fn-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's the Fn server
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'functions'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      # Specify all the fn servers from which metrics will be scraped
      - targets: ['fnserver0:8080','fnserver1:8081'] # Uses /metrics by default      
```
Note the last line: this configures Prometheus to scrape metrics from two Fn servers running on `fnserver0:8080` and `fnserver1:8081`, 
where `fnserver0` and `fnserver0` are defined in [docker-compose.yml](https://github.com/fnproject/ext-statsapi/blob/master/examples/operators-clustered/docker-compose.yml)
to refer to the two Fn servers.

Now start your custom Fn images and Prometheus

```sh
cd $GOPATH/src/github.com/fnproject/ext-statsapi/examples/operators-clustered
docker-compose up
```

You can now deploy and run functions and try out the statistics API extension.

## Trying out the statistics API with a cluster of custom Fn servers

Create some simple cold async functions
```
cd $GOPATH/src/github.com/fnproject/ext-statsapi/test/hello-cold-async-a
fn deploy --all --local
```
Now run these async functions, some on one Fn server and some on the other. You may want to run this several times to generate plenty of data.
Note that each run performs 90 function calls.
```
cd $GOPATH/src/github.com/fnproject/ext-statsapi/examples/operators-clustered
bash run-cold-async-clustered.bash
```
You can now use the statistics API to obtain aggregated metrics across both Fn servers. If you ran the above script once you should see jhe number of calls grow to 90.  
```
curl localhost:8080/v1/apps/hello-cold-async-a/stats
```
Alternatively open a browser at (localhost:8080/v1/apps/hello-cold-async-a/stats)

Note that you can use this API on either of the two custom Fn servers and get the same result. Try the following and compare:
```
curl localhost:8081/v1/apps/hello-cold-async-a/stats
```
Alternatively open a browser at (localhost:8081/v1/apps/hello-cold-async-a/stats)
