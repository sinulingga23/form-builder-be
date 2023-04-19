# Form Builder Backend

## Run project (choose one below)

### Using Docker Compose
```
docker-compose up -d
```

### Using Docker only
Create an docker network:
```
docker network create dev
```

Run an redis server:
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

Run an nginx server:
```
docker run --name nginx -p "80:80" -v $(pwd)/default.conf:/etc/nginx/default.conf -d nginx:1.23.3
```

### Manually
cooming soon

## Metrics Monitoring
In this project using Prometheus for capture metrics from the program.
* Try hit the services from the available endpoints as much as you want.
* Open ```localhost:8087/metrics``` at your favorite browser
* You will see somethings like these:
```
# HELP form_builder_be_duration_request_endpoint It's show duration request for each endpoint
# TYPE form_builder_be_duration_request_endpoint histogram
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="200",service_name="m_form_service:add_form",le="0.1"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="200",service_name="m_form_service:add_form",le="0.15"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="200",service_name="m_form_service:add_form",le="0.2"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="200",service_name="m_form_service:add_form",le="0.25"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="200",service_name="m_form_service:add_form",le="0.3"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="200",service_name="m_form_service:add_form",le="+Inf"} 20
form_builder_be_duration_request_endpoint_sum{http_method="POST",http_status="200",service_name="m_form_service:add_form"} 2.72441375e+08
form_builder_be_duration_request_endpoint_count{http_method="POST",http_status="200",service_name="m_form_service:add_form"} 20
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="400",service_name="m_form_service:add_form",le="0.1"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="400",service_name="m_form_service:add_form",le="0.15"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="400",service_name="m_form_service:add_form",le="0.2"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="400",service_name="m_form_service:add_form",le="0.25"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="400",service_name="m_form_service:add_form",le="0.3"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="400",service_name="m_form_service:add_form",le="+Inf"} 14
form_builder_be_duration_request_endpoint_sum{http_method="POST",http_status="400",service_name="m_form_service:add_form"} 2.437624e+06
form_builder_be_duration_request_endpoint_count{http_method="POST",http_status="400",service_name="m_form_service:add_form"} 14
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="404",service_name="m_form_service:add_form",le="0.1"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="404",service_name="m_form_service:add_form",le="0.15"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="404",service_name="m_form_service:add_form",le="0.2"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="404",service_name="m_form_service:add_form",le="0.25"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="404",service_name="m_form_service:add_form",le="0.3"} 0
form_builder_be_duration_request_endpoint_bucket{http_method="POST",http_status="404",service_name="m_form_service:add_form",le="+Inf"} 7
form_builder_be_duration_request_endpoint_sum{http_method="POST",http_status="404",service_name="m_form_service:add_form"} 1.136209e+06
form_builder_be_duration_request_endpoint_count{http_method="POST",http_status="404",service_name="m_form_service:add_form"} 7
# HELP form_builder_be_total_request_endpoint It's show total request for each endpoint
# TYPE form_builder_be_total_request_endpoint counter
form_builder_be_total_request_endpoint{http_method="POST",http_status="200",service_name="m_form_service:add_form"} 20
form_builder_be_total_request_endpoint{http_method="POST",http_status="400",service_name="m_form_service:add_form"} 14
form_builder_be_total_request_endpoint{http_method="POST",http_status="404",service_name="m_form_service:add_form"} 7
```

