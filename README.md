Form Builder Backend

# Build Image from Dockerfile
```
docker build -f Dockerfile -t github.com/sinulingga23/form-builder-be .
```

# Run Container from an Image
```
docker run -e DB_HOST=host.docker.internal -e DB_USER=dennyrezky -e DB_PASSWORD= -e DB_NAME=form_builder -e DB_SSL_MODE=disable -e PORT=8087 -e REDIS_SERVER_NETWORK=localhost:6379 -e REDIS_SERVER_PASSWORD= -p 8087:8087 --expose 8087 --name form-builder-be -d github.com/sinulingga23/form-builder-be
```
