# DockerMV
DockerMv is an extended version of Docker which supports software multi-versioning in services. By using DockerMV, we can create Docker services with more than one image.

## How to use?
In order to run this program, follow these steps:
1) Install the [go programming language](https://golang.org/dl/) as following:

    1.1) Run the following commands:
    ```
    wget https://dl.google.com/go/go1.13.6.linux-amd64.tar.gz
    tar -xvf go1.13.6.linux-amd64.tar.gz
    sudo mv go /usr/local
    ```
    1.2) Now set the GOPATH variable, which tells GO where to look for its files:
    ```
    sudo nano ~/.profile
    ```
    At the end of the file add this line:
    ```
    export GOPATH=$HOME/DockerMV/go
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    ```
    Now save and close the file, and refresh your profile.
    ```
    source ~/.profile
    ```

2) Install [docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/) as follows:

    2.1) Run the following commands:
    ```
    sudo apt install apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
    apt-cache policy docker-ce
    ```
    At the end you should see output similar to this:
    ```
    Installed: (none)
    Candidate: 5:19.03.5~3-0~ubuntu-bionic
    Version table:
        5:19.03.5~3-0~ubuntu-bionic 500
            500 https://download.docker.com/linux/ubuntu bionic/stable amd64 Packages
    ```
    2.2) Finally install docker
    ```
    sudo apt install docker-ce
    ```

3) Clone this project and put in you $HOME directory.

4) Run the following command from your $HOME directory.
```
chmod -R 777 DockerMV
```

5) Move to $HOME/DockerMV/go/src/github.com/docker/cli directory and run the following command to build the project.
```
make binary
```
You will see a similar message when the program is successfuly compiled:
```
WARNING: you are not in a container.
Use "make -f docker.Makefile binary" or set
DISABLE_WARN_OUTSIDE_CONTAINER=1 to disable this warning.

Press Ctrl+C now to abort.

WARNING: binary creates a Linux executable. Use cross for macOS or Windows.
./scripts/build/binary
Building statically linked build/docker-linux-amd64
```

6) To run a command using DockerMV, first move to $HOME/DockerMV/go/src/github.com/docker/cli directory, and your command needs to start with ./build/docker

## Experiments

### TeaStore
The [TeaStore](https://github.com/DescartesResearch/TeaStore) application is a reference application for testing and benchmarking. You can find two version of its Recommender service on our [Docker Hub](https://hub.docker.com/u/sgholami) page. Also, we used the [teastore.jmx](teastore.jmx) to create a load on the system for our testing purposes.

#### How to Run the TeaStore with DockerMV
The TeaStore's [wiki page](https://github.com/DescartesResearch/TeaStore/wiki/Getting-Started#run-teastore-containers-using-docker) has a complete instruction on how to install the TeaStore on Docker. The followings instruct on how to use the TeaStore with two different Recommender Services. Remember to replace your IP addresses in the following commands.

```
docker run -e "HOST_NAME=10.2.5.26" -e "SERVICE_PORT=10000" -p 10000:8080 -d descartesresearch/teastore-registry

docker run -p 3306:3306 -d descartesresearch/teastore-db

docker run -e "REGISTRY_HOST=10.2.5.26" -e "REGISTRY_PORT=10000" -e "HOST_NAME=10.2.5.26" -e "SERVICE_PORT=1111" -e "DB_HOST=10.2.5.26" -e "DB_PORT=3306" -p 1111:8080 -d descartesresearch/teastore-persistence

docker run -e "REGISTRY_HOST=10.2.5.26" -e "REGISTRY_PORT=10000" -e "HOST_NAME=10.2.5.26" -e "SERVICE_PORT=2222" -p 2222:8080 -d descartesresearch/teastore-auth

docker run -e "REGISTRY_HOST=10.2.5.26" -e "REGISTRY_PORT=10000" -e "HOST_NAME=10.2.5.26" -e "SERVICE_PORT=4444" -p 4444:8080 -d descartesresearch/teastore-image

docker run -e "REGISTRY_HOST=10.2.5.26" -e "REGISTRY_PORT=10000" -e "HOST_NAME=10.2.5.26" -e "SERVICE_PORT=8080" -p 8080:8080 -d descartesresearch/teastore-webui

./build/docker service create e REGISTRY_HOST=10.2.5.26 e REGISTRY_PORT=10000 e HOST_NAME=10.2.5.26 e SERVICE_PORT=3333 10.2.5.26 my-net my_recommender 8080 my_rule.txt sgholami/teastore-recommender:MultipleTrain 1 sgholami/teastore-recommender:SingleTrain 1

```

### Znn
The [Znn](https://github.com/cmu-able/znn) application is used for testing and benchmarking of self-adaptive applications. We created two version of its content-providing component which are available on our [Docker Hub](https://hub.docker.com/u/alirezagoli) page. Also, we used the [znn.jmx](znn.jmx) to create a load on the system for our testing purposes.

## Cite Us

The DockerMV was first published in Proceedings of the 2020 ACM/SPEC International Conference on Performance Engineering (ICPE '20). You can find the paper on the [ASGAARD lab website](https://www.google.com/url?q=http://asgaard.ece.ualberta.ca/publications/&sa=D&source=hangouts&ust=1579122442788000&usg=AFQjCNFElRVZ9AvFDUP-bTIoO4r5-XdNlg).
