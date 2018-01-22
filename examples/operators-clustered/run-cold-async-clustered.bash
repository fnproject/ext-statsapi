#!/bin/bash

# Execute a variety of cold async functions, some on the Fn server listening on 8080 and some on the Fn server listening on 8081
for value in {1..10}
do
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a2
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a3
  
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a2 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1

  curl localhost:8081/r/hello-cold-async-a/hello-cold-async-a1 
  curl localhost:8081/r/hello-cold-async-a/hello-cold-async-a2 
  curl localhost:8081/r/hello-cold-async-a/hello-cold-async-a1
done