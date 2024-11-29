
#  Key Value Store service


This is a key-value store to save and retrieve data.

The project has two services:
- API service
- Key-value store service


The key-value store service is implemented using an in-memory store, which means that the data is stored in memory and is lost when the service is restarted.

The API service consists RESTful API's that allows clients to interact with the key-value store. The API has endpoints for setting, getting, and deleting key-value pairs.




## Environment Variables

To run this project, you will need to add/modify the following environment variables to your .env file. This file has to be located at the root of the project.

`KVSTORE_HOST`

`KVSTORE_PORT`

`API_HOST`

`API_PORT`

Current configuration

```bash
KVSTORE_HOST=kvstore-grpc
KVSTORE_PORT=50510

API_HOST=rest-api
API_PORT=8081

```

## Deployment


From the root of the project

First generate the proto files necessary for grpc communication

```bash
$ make proto
```

Run the docker compose file to spin up the service and the key value store.

```bash
$ docker-compose up --build
```

Run this to tear down the service

```bash
$ docker-compose down
```
## Usage/Examples



##### Create the key value pair that has the key "test" and value "test"


```bash
curl --location 'localhost:8081/store' \
--header 'Content-Type: application/json' \
--data '{
    "key": "test",
    "value": "test"
}'
```



##### Retrieve the value of key value pair that has the key "test"


```bash
curl --location 'localhost:8081/store?key=test'
```

##### Delete the key value pair that has the key "test"


```bash
curl --location --request DELETE 'localhost:8081/store/test'
```


## Running Tests

From the root of the project run

```bash
  go test ./...
```


## API Reference

### Store key-value pair

```bash
  POST /store
```
Body:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `key` | `string` | **Required**. Key of the key-value pair|
| `value` | `string` | **Required**. Value of the key-value pair|

Response:

| Parameter | Type     | Description                             |
| :-------- | :------- |:----------------------------------------|
| `success` | `bool` | Status of the operation.|


### Get value of the key-value pair

```bash
  GET /store?key={key}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `key` | `string` | **Required**. Key of the key-value pair|


Response:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `value` | `string` |  Value of the key-value pair|


### Delete key-value pair

```bash
  DELETE /store/{key}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `key` | `string` | **Required**. Key of the key-value pair|

Response:

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `success` | `bool` | Status of the operation. |




