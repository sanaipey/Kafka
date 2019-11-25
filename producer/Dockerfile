#build from this base image
FROM abhirockzz/kafka_go_client
#set the working directory of image to be go/src/app
WORKDIR /go/src/app
#copy main.go producer file to current directory of docker image
COPY main.go .
#sets the entry point of the application when it runs| species how to run a container based off of this image
CMD ["go", "run", "main.go"]

