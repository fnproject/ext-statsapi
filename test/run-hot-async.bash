#!/bin/bash
for value in {1..50}
do
  curl -X POST -d '{"name":"hello-hot-async-a/hello-hot-async-a1"}' http://localhost:8080/r/hello-hot-async-a/hello-hot-async-a1
done
