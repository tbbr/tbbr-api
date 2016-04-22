export GOOS=linux
export GOARCH=amd64
go build -v payup
ssh root@159.203.33.146 service tbbr stop
scp payup root@159.203.33.146:/home/tbbr/tbbr-server/tbbr
ssh root@159.203.33.146 service tbbr start
export GOOS=darwin
