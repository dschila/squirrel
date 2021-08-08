<h1 align="center">
squirrel - share-backend-api
</h1>

<p align="center">REST API <b>backend</b> (Golang) for <a href="https://github.com/beyeja/peacock/">peacock</a> to store files into a minIO-Bucket</p>


## ⚡️ Quick start 

[Download](https://golang.org/dl/) and install **Go**. Version `1.16` or higher is required.

[Download](https://docs.docker.com/get-docker/) and install **Docker**. Make sure you have installed Docker Compose as well

Start required services **mongo-db** and **minIO** by using docker-compose:

```bash
docker-compose up -d
```

Let's start the backend service
```bash
go run main.go
```
