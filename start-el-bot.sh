# workdir=$(cd $(dirname $0); pwd)
cd src/main
# export GOPATH=$workdir:$GOPATH
echo $GOPATH
export QQ=$1
export AUTHKEY=$2
go run main.go