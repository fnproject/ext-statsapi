#!/bin/bash
for value in {1..10}
do
  echo $value
  sleep 1
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a2
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a3
  
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a2 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1
 
  curl localhost:8080/r/hello-cold-async-b/hello-cold-async-b1 
  curl localhost:8080/r/hello-cold-async-b/hello-cold-async-b2
  curl localhost:8080/r/hello-cold-async-b/hello-cold-async-b3
done