Form Builder Backend

# Run this project.
Create an docker network:
```
docker network create dev
```

Run an redis server using docker:
```
docker run --name redisserver --network dev -p "127.0.0.1:6379:6379" -d redis:6.2.11 redis-server
```

Build an image from this Dockerfile:
```
docker build -f Dockerfile -t github.com/sinulingga23/form-builder-be .
```

Run an container using latest image of this project:
```
docker run -e DB_HOST=host.docker.internal -e DB_USER=dennyrezky -e DB_PASSWORD= -e DB_NAME=form_builder -e DB_SSL_MODE=disable -e PORT=8087 -e REDIS_SERVER_NETWORK=redisserver -e REDIS_SERVER_PASSWORD= -p 8087:8087 --expose 8087 --name form-builder-be --network dev -d github.com/sinulingga23/form-builder-be
```

