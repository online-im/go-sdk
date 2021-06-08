export GOOS="linux"
export GOARCH="amd64"
go build -o online-im-go-client .
docker build --no-cache -t online-im-go-client-image  .
rm ./online-im-go-client
kubectl create namespace glory
kubectl delete -f ./client.yaml
kubectl create -f ./client.yaml