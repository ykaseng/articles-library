# articles-library

## Introduction
This is a simple API to create and retrieve articles through using RESTful HTTP implemented in Go.
To run, simply use the command:
```
make serve
```

Make serve will build an image with a base Go image and get all the required packages specified in the module and serve the service at port `8080`. Make serve will also spawn a PostgreSQL database, populate a `library` database with tables from `start.sh`.

To run the tests, use the command:
```
make test
```

Make test runs run `go test -v ./...` and will similarly spawn a PostgreSQL database but will unmount the database volume at the end of each test.

## API Interface
### Create Article
- Method: `POST`
- Path: `/articles`
- Request Body:
```JSON
{
    "title": "Hello World",
    "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
    "author": "John"
}
```
- Response Header: `HTTP 201`
- Response Body:
```JSON
{
    "status": 201,
    "message": "Success",
    "data": {
      "id": <article_id>
    }
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```JSON
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```

Sample Request:
```cURL
curl -X POST \
  http://localhost:8080/articles \
  -H 'Accept: */*' \
  -H 'Accept-Encoding: gzip, deflate' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Length: 517' \
  -H 'Content-Type: text/plain' \
  -H 'Host: localhost:8080' \
  -H 'Postman-Token: 2ac0009b-3fe7-43ad-93a5-80d7690f7095,a2e286aa-9711-4742-92a1-fc848b997714' \
  -H 'User-Agent: PostmanRuntime/7.18.0' \
  -H 'cache-control: no-cache' \
  -d '{
    "title": "Hello World",
    "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
    "author": "John",
}'
```

### Get Article by ID
- Method: `GET`
- Path: `articles/<article_id>`
- Response Header: `HTTP 200`
- Response Body:
```JSON
{
    "status": 200,
    "message": "Success",
    "data": [
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      }
    ]
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```JSON
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```

Sample Request:
```cURL
curl -X GET \
  http://localhost:8080/articles/1 \
  -H 'Accept: */*' \
  -H 'Accept-Encoding: gzip, deflate' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Host: localhost:8080' \
  -H 'Postman-Token: 0959d617-30c5-4dc7-bcac-0958e937a771,b730da7e-53c7-4fed-b73f-56934f713ed9' \
  -H 'User-Agent: PostmanRuntime/7.18.0' \
  -H 'cache-control: no-cache'
```

### Get All Articles
- Method: `GET`
- Path: `/articles`
- Response Header: `HTTP 200`
- Response Body:
```JSON
{
    "status": 200,
    "message": "Success",
    "data": [
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      },
      {
        "id": <article_id>,
        "title":<article_title>,
        "content":<article_content>,
        "author":<article_author>,
      }
    ]
}
```
or
- Response Header: `HTTP <HTTP_CODE>`
- Response Body:
```JSON
{
    "status": <HTTP_CODE>,
    "message": <ERROR_DESCRIPTION>,
    "data": null
}
```

Sample Request:
```cURL
curl -X GET \
  http://localhost:8080/articles \
  -H 'Accept: */*' \
  -H 'Accept-Encoding: gzip, deflate' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Host: localhost:8080' \
  -H 'Postman-Token: d98d29a2-586f-4b99-b908-b85c0819e174,6cfa3758-3d23-4c08-b706-9c039d1d05a0' \
  -H 'User-Agent: PostmanRuntime/7.18.0' \
  -H 'cache-control: no-cache'
```