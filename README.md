# Lucas

<a href="https://www.youtube.com/watch?v=VrS6akzR3sk"><img src="https://cdn.davidwolfe.com/wp-content/uploads/2017/11/spider-video-FI.jpg"/></a>

Lucas is a web crawler built using [Go](https://golang.org/) and the [Colly](https://github.com/gocolly/colly) library.

It is currently setup to crawl floryday -> it will write its output to the connected psql DB and output the results of its latest crawl to the console

## Running Locally

### DB Setup

```
docker-compose up -d
psql -h localhost -U user lucas_db -f dbsetup.sql
```

### Run Lucas

`go run lucas.go`

## Docker

To run Lucas in a Docker container

```
# build docker image
docker build .

# run docker container and portforward port 3000
docker run -ti -p 8000:8000 --network="host" <docker-image-id>

# publish docker image to docker hub
docker push <docker-repo>
```
