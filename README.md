# DockerMV
DockerMv is an extended version of Docker which supports software multi-versioning in services. By using DockerMV, we can create Docker services with more than one image.

## How to use?
In order to run this program, follow these steps:
1) Install the [go programming language](https://golang.org/dl/) as following:
1.1) Run the following commands:
```
sudo apt-get update
sudo apt-get -y upgrade
wget https://dl.google.com/go/go1.13.6.linux-amd64.tar.gz
sudo tar -xvf go1.13.6.linux-amd64.tar.gz
sudo mv go /usr/local
```
1.2) Now set the GOPATH variable, which tells GO where to look for its files:
```
sudo nano ~/.profile
```
At the end of the file add this line:
```
export GOPATH=$HOME/work
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```
Now save and close the file, and refresh your profile.
```
source ~/.profile
```

2) Install [docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/) 

3) Download this project and put in you GO home directory, e.g., go/src/github/ directory. 

4) Move to docker/cli directory and run the following command to build the project.
```
make binary
```

5) To run a command using DockerMV, your command needs to start with ./build/docker

## Experiments

### TeaStore
The [TeaStore](https://github.com/DescartesResearch/TeaStore) application is a reference application for testing and benchmarking. You can find two version of its Recommender service on our [Docker Hub](https://hub.docker.com/u/sgholami) page. Also, we used the [teastore.jmx](teastore.jmx) to create a load on the system for our testing purposes.

### Znn
The [Znn](https://github.com/cmu-able/znn) application is used for testing and benchmarking of self-adaptive applications. We created two version of its content-providing component which are available on our [Docker Hub](https://hub.docker.com/u/alirezagoli) page. Also, we used the [znn.jmx](znn.jmx) to create a load on the system for our testing purposes.
