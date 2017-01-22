# tbbr-api
[![Build Status](https://travis-ci.org/tbbr/tbbr-api.svg?branch=master)](https://travis-ci.org/tbbr/tbbr-api)

The api that serves data to the tbbr web application and the tbbr android application.

## Setup

#### 1) Clone tbbr-api into your GOPATH

```
$ cd $GOPATH/src/github.com/tbbr/
$ git clone https://github.com/tbbr/tbbr-api
```

#### 2) Install glide (dependency manager)
On Mac OS X you can install the latest release via Homebrew:
```
$ brew install glide
```

On Ubuntu Precise(12.04), Trusty (14.04), Wily (15.10) or Xenial (16.04) you can install from our PPA:

```
sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update
sudo apt-get install glide
```

For more information on glide checkout: https://github.com/Masterminds/glide

#### 3) Get all dependencies

```
$ glide install
```

#### 4) Install postgreSQL

#### 5) Create tbbr and tbbr_test database

#### 6) Install migrate into your GOPATH

```
$ go get github.com/mattes/migrate
```

#### 7) Run migrations

```
$ migrate -url postgres://username@host:port/test_database_name?sslmode=disable -path ./migrations up
```

Concrete example
```
$ migrate -url postgres://maazali@localhost:5432/tbbr_test?sslmode=disable -path ./migrations up
```


#### 8) Export tbbr specific database environment variables
Update ~/.bash_profile to insert the following
```
export TBBR_DB_NAME=tbbr
export TBBR_DB_USER=maazali
export TBBR_DB_PASSWORD=(your username's password)
```

## Build and Run

```
$ go build
$ ./tbbr-api
```

## Testing

```
$ go test
```
