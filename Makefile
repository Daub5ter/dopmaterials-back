openapi_run:
	docker run --name openapi -d -p 80:8080 -v ./api-gateway/api/api-gateway-openapi.yaml:/usr/share/nginx/html/openapi.yaml swaggerapi/swagger-ui

openapi_stop:
	docker rm -f openapi