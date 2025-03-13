IMAGE_NAME = notiffy-service
CONTAINER_NAME = notiffy-container
# if commands arent working u just need to add sudo to prefix each command
build:
	docker build -t $(IMAGE_NAME) .
run:
#	In the -p XX:YY flag, the XX is the host port, while YY is the port within the container. 
	docker run --name $(CONTAINER_NAME) -d -p 50051:50051 -p 9999:9999 $(IMAGE_NAME)
stop:
	docker stop $(CONTAINER_NAME)
start:
	docker start $(CONTAINER_NAME)
restart:
	docker restart $(CONTAINER_NAME)

logs:
	docker logs $(CONTAINER_NAME)

clean:
	docker rm -f $(CONTAINER_NAME)
	docker rmi -f $(IMAGE_NAME)
# Notes
#to go inside docker container docker exec -it CONTAINER_ID /bin/sh
#docker stats
