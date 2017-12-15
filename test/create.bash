WD=$PWD
cd $GOPATH/src/github.com/fnproject/ext-metrics/test/hello-cold-sync-a
fn deploy --all --local
cd $GOPATH/src/github.com/fnproject/ext-metrics/test/hello-cold-async-a
fn deploy --all --local
cd $GOPATH/src/github.com/fnproject/ext-metrics/test/hello-cold-async-b
fn deploy --all --local
cd $GOPATH/src/github.com/fnproject/ext-metrics/test/hello-hot-sync-a/
fn deploy --all --local
cd $GOPATH/src/github.com/fnproject/ext-metrics/test/hello-hot-async-a/
fn deploy --all --local
cd $WD
