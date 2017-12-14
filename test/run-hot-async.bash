  curl localhost:8080/r/hello-cold-async-a/hello-cold-async-a2
#!/bin/bash
for value in {1..20}
do
  echo $value
  sleep 1
  curl -X POST -d '{"name":"hello-hot-async-a/hello-hot-async-a1"}' http://localhost:8080/r/hello-hot-async-a/hello-hot-async-a1
done
