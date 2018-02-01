WD=$PWD
cd $GOPATH/src/github.com/fnproject/ext-statsapi/test
bash run-cold-sync.bash
bash run-hot-sync.bash
bash run-cold-async.bash
bash run-hot-async.bash
cd $WD