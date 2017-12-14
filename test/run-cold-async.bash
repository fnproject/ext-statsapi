#!/bin/bash
for value in {1..10}
do
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a2
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a3
  
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a2 
  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a1
 
  curl localhost:8080/r/hello-cold-async-b/hello-cold-async-b1 
  curl localhost:8080/r/hello-cold-async-b/hello-cold-async-b2
  #don't call this last one so we can check the stats of a function that has never been called
  #curl localhost:8080/r/hello-cold-async-b/hello-cold-async-b3
done