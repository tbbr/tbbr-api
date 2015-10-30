export GOOS=linux
export GOARCH=amd64
go build -v payup
ssh root@23.92.222.14 service payup stop
scp payup root@23.92.222.14:/home/maazali/payup-server
ssh root@23.92.222.14 service payup start
export GOOS=darwin
