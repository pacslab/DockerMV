# DockerMV
DockerMv is an extended version of Docker which supports software multi-versioning in services. By using DockerMV, we can create Docker services with more than one image.

## How to use?
In order to run this program, follow these steps:
1) Install [go programming language](https://golang.org/dl/).

2) Install [docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/) 

3) Download this project and put in you GO home directory, e.g., go/src/github/ directory. 

4) Move to docker/cli/build directory and run the following command to build the project.
```
make binary
```

5) To run a command using DockerMV, your command needs to start with ./build/docker
