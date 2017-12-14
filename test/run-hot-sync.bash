#!/bin/bash
for value in {1..10}
do
  curl -X POST -d '{"name":"hello-hot-sync-a/hello-hot-sync-a1"}' http://localhost:8080/r/hello-hot-sync-a/hello-hot-sync-a1
done
