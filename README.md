## Introduction
airy is a lightweight distributed game server, using grpc to do communication between processes, in line with the design principle of convenient expansion module, a type of service is called a module, you can like building blocks easily for your service module loading different components, you can complete the design of a module, of course, you can also expand the components you need.

## Getting Started
### Prerequisites
* [Go](https://golang.org/) >= 1.18
* [etcd](https://github.com/coreos/etcd) (used for service discovery)
* [redis](https://github.com/redis/redis) (used for storage)

### Installing
clone the repo
```
git clone https://github.com/lopies/airy.git
```

### Configuration
currently only a distributed version is available, so if you want to start the service, you need to complete the following preparations
```
cd config/
```
edit configuration(etcd / redis)

### Start
once you have completed the above steps, you are ready to start the service
```
cd cmd && go run logic.go
cd cmd && go run gate.go
```

## ## Authors
* **Franco Zhou**

## Test
robot test program case is provided, you can open the client test(you can also edit config file or behavior file)
```
cd client && go run main.go
```


