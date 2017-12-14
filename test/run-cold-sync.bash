#!/bin/bash
for value in {1..10}
do
  curl localhost:8080/r/hello-cold-sync-a/hello-cold-sync-a1 
done