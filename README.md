# fake-rest-server

Simple Rest Api to test Rest client.
Server always respond with CORS headers allowing all requests.

##### Get dependencies
``
go get
``

##### Build

``
go build
``

##### Run

``
./fake-rest-server -bind=:8080
``

##### Usage
``` 
Usage of ./fake-rest-server:
  -bind string
        http bind port (default ":8080")
  -data string
        Filename to load data. Data must be named *.yaml or *.json format, see data.example.yaml
```

### Examples

- Add resource data

```
curl --request POST \
  --url http://localhost:8080/resources/1 \
  --header 'content-type: application/json' \
  --data '{"toto":"tata"}'
```
  
- Get resource data 

```
curl --request GET \
  --url http://localhost:8080/resources/1
```

- delete resource 

```
curl --request DELETE \
  --url http://localhost:8080/resources/1
```

### Docker

- build image

```
docker build -t fake-rest-server .
```

- run image

```
docker run -p 8080:8080 fake-rest-server
```

- Run image loading data file


```
docker run -p 8081:8080 -v $PWD:/mnt -it -e DATA=/mnt/data.example.yaml fake-rest-server
```
