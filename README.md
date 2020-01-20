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

6) To run services with DockerMV, you need to create a Docker Swarm. Check the public IP address of your machine on your network. Notice that the IP address must be accessible from outside.
```
hostname -I | awk '{print $1}'
```
Then, you need to create an overlay network by running the following command:
```
sudo docker swarm init --advertise-addr HOST_IP
sudo docker network create -d overlay --attachable my-net
```

7) To run a command using DockerMV, first move to $HOME/DockerMV/go/src/github.com/docker/cli directory, and your command needs to start with ./build/docker

## Experiments

### TeaStore
The [TeaStore](https://github.com/DescartesResearch/TeaStore) application is a reference application for testing and benchmarking. You can find two version of its Recommender service on our [Docker Hub](https://hub.docker.com/u/sgholami) page. Also, we used the [teastore.jmx](teastore.jmx) to create a load on the system for our testing purposes.

### Znn
The [Znn](https://github.com/cmu-able/znn) application is used for testing and benchmarking of self-adaptive applications. We created two version of its content-providing component which are available on our [Docker Hub](https://hub.docker.com/u/alirezagoli) page. Also, we used the [znn.jmx](znn.jmx) to create a load on the system for our testing purposes.

#### How to Run the Znn with DockerMV
You can setup the Znn application with the following commands. Notice to replace the HOST_IP with your host IP address.
```
sudo docker run --network="my-net" -d -p 3306:3306 alirezagoli/znn-mysql:v1

./build/docker service create HOST_IP my-net my_znn 1081 $HOME/DockerMV/znn_sample_rule.txt alirezagoli/znn-text:v1 1 1g 1g 0.2 alirezagoli/znn-multimedia:v1 1 1g 1g 0.2
```

By running th following command you can see four containers are running.
```
sudo docker ps -a
```

Now, you can see the service working by running the following command. Notice that the NGINX port is randomly assigned and you can find it by the above command.
```
curl http://HOST_IP:NGINX_PORT/news.php
```

## How to remove containers?
In order to remove the containers you need to first stop them.
```
sudo docker stop CONTAINER_ID
```

Afterward, you can remove the containers 
```
sudo docker rm CONTAINER_ID
```

## Cite Us

The DockerMV was first published in Proceedings of the 2020 ACM/SPEC International Conference on Performance Engineering (ICPE '20). You can find the paper on the [ASGAARD lab website](https://www.google.com/url?q=http://asgaard.ece.ualberta.ca/publications/&sa=D&source=hangouts&ust=1579122442788000&usg=AFQjCNFElRVZ9AvFDUP-bTIoO4r5-XdNlg).
