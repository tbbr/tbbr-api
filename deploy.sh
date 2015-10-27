export GOOS=linux
export GOARCH=amd64
go build -v payup
scp payup maazali@23.92.222.14:~/payup-server
export GOOS=darwin
