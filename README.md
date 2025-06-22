# Censys Take Home

## Problem Statement
Build a simple decomposed Key-Value store by implementing two services which communicate over gRPC.

The first service should implement a basic JSON Rest API to serve as the primary public interface. This service should then externally communicate with a second service over gRPC, which implement a basic Key-Value store service that can:

- Store a value at a given key
- Retrieve the value for a given key
- Delete a given key

The JSON interface should at a minimum be able to expose and implement these three functions.

You can write this in whichever languages you choose, however Go would be preferred. The final result should be built into two separate Docker containers which can be used to run each service independently.

Please upload the code to a publicly accessible GitHub, GitLab or other public code repository account. A README file should be provided, briefly documenting what you are delivering. Like our own code, we expect testing instructions: whether it’s an automated test framework, or simple manual steps.

## Design Decisions
Cache:
1. Use `map[string][]byte` 
   1. `key string` - don't think we will need to have a `key` with other data types, could use `comparable` in that case, but don't think it is needed.
   2. `value []byte` - `[]byte` is advantageous for use with marshaling and unmarshaling into structs. Initially tried to use `any`, but then you get into some oddness with gRPC, so I pivoted. Could have also used `interface{}`.
2. Create cache interface `CacheService` which can be also passed to other internal services, or to a different API server than gRPC.
3. No unit or race condition testing added, didn't have time.

gRPC:
1. Simple contract to pass caching to public API.
2. No unit or integration testing added, didn't have time.

API:
1. Simple API which implements create, read (no update) and delete endpoints, which calls the gRPC cache server. 
2. Minor logging added for troubleshooting live deploy.
3. Didn't separate out handler or routes from main, as this is a simple API implementation.
4. No integration testing added, no time to implement.

## Setup
1. Install `git` (if not already installed).
   1. [Git download](https://git-scm.com/downloads)
   2. Homebrew
  ```sh
  brew install git
  ```
  
2. Install Docker Client (if not already installed), and start the client.
   1. [Docker Desktop Client download](https://www.docker.com/products/docker-desktop/)
   2. Homebrew
   ```sh
   brew install docker
   ```
   
3. Clone the repository, into a directory of your choosing.
```sh
git clone https://github.com/tjroeder/censys_take_home.git && cd censys_take_home
```

4. In the terminal, run:
```sh
docker-compose up --build
```
> Note: to close the docker containers, either use the docker client or `ctrl+c` in the terminal window.

5. You now may make HTTP requests to the API server hosted on the `api-gateway` container running on `localhost:8080`, or gRPC request to the gRPC server hosted on `grpc-server` container running on `localhost:50051`.

## Testing
There is no automated testing, to test the API, run curl commands from the terminal or other HTTP request client, e.g. Postman, etc

1. Test `GET` key-value record before storing a record. Tests `404 NOT FOUND` response.
```sh
# GET /v1/keyvalues/{key}

# Example 
curl -i http://localhost:8080/v1/keyvalues/user_1

# Expected Output
# HTTP/1.1 404 Not Found
# Content-Type: text/plain; charset=utf-8
# X-Content-Type-Options: nosniff
# Date: Sun, 22 Jun 2025 21:49:46 GMT
# Content-Length: 10
```

2. Test `CREATE` key-value record, but with invalid JSON. Tests `400 BAD REQUEST` response.
```sh
# POST /v1/keyvalues
# Header: Content-Type=application/json
# Body: {"key":string, "value":string}

# Example
curl -X POST http://localhost:8080/v1/keyvalues \
  -H "Content-Type: application/json" \
  -d '{"key":"foo","value":"{"id":123,"username":"tim"}"'

# Output
# HTTP/1.1 400 Bad Request
# Content-Type: text/plain; charset=utf-8
# X-Content-Type-Options: nosniff
# Date: Sun, 22 Jun 2025 21:22:48 GMT
# Content-Length: 12
```

3. Test `CREATE` key-value valid record. Tests `201 CREATED` response.
```sh
# POST /v1/keyvalues 
# Header: Content-Type=application/json
# Body: {"key":string, "value":string}

# Example
curl -i -X POST http://localhost:8080/v1/keyvalues \
  -H "Content-Type: application/json" \
  -d '{"key":"user_1","value":"{\"id\":1,\"username\":\"tim\",\"Email\":{\"email\":\"tim@email.com\",\"active\":true}}"}'

# Output
# HTTP/1.1 201 Created
# Date: Sun, 22 Jun 2025 20:08:01 GMT
# Content-Length: 0
```

This test is interesting because it has embedded json in value, which can be unmarshaled into the User struct below.
```go
type Email struct {
  Email  string `json:"email"`
  Active bool   `json:"active"`
}

type User struct {
  ID       int    `json:"id"`
  Username string `json:"username"`
  Email    Email
}
```

1. Test `GET` key-value record after storing a record. Tests `200 OK` response.
```sh
# GET /v1/keyvalues/{key}
# Example 
curl -i http://localhost:8080/v1/keyvalues/user_1

# Expected Output
# HTTP/1.1 200 OK
# Content-Type: application/json
# Date: Sun, 22 Jun 2025 20:09:11 GMT
# Content-Length: 100

# {"value":"{\"id\":1,\"username\":\"tim\",\"Email\":{\"email\":\"tim@email.com\",\"active\":true}}"}
```

5. Test `DELETE` key-value record after storing a record. Tests `204 NO CONTENT` response.
```sh
# DELETE /v1/keyvalues/{key}
# Example 
curl -i -X DELETE http://localhost:8080/v1/keyvalues/user_1

# Expected Output
# HTTP/1.1 204 No Content
# Date: Sun, 22 Jun 2025 20:09:37 GMT
```

6. Test `GET` key-value record after deleting the stored record. Tests `404 NOT FOUND` response.
```sh
# GET /v1/keyvalues/{key}
# Example 
curl -i http://localhost:8080/v1/keyvalues/user_1

# Expected Output
# HTTP/1.1 404 Not Found
# Content-Type: text/plain; charset=utf-8
# X-Content-Type-Options: nosniff
# Date: Sun, 22 Jun 2025 21:49:46 GMT
# Content-Length: 10
```

### Testing of gRPC caching service failure
1. Terminate the gRPC container
2. Run curl commands for:
   1. `GET /v1/keyvalues/{key}`
   2. `POST /v1/keyvalues`
   3. `DELETE /v1/keyvalues/{key}`
3. Verify the return is `500 Internal Server Error`

## Future Development Ideas
### Internal Cache
- Add Cache Persistance
- Add default or optional TTL setting
- Add updating cache value TTL 
- Add setting "`-`" for cache value to be recognized as a `nil` value for when a GET request has already been attempted and there was a cache miss.
- Separate out Get and GetOk functions
- Add "Add" function 
- Create unit tests

### Public API
- Separate out handler and router from Main
- Create handler tests
- Create integration tests
- Implement logging/monitoring service for internal server errors and TraceID

### gRPC Server
- Create gRPC server tests
- Create in-depth testing instructions for testing the service using gRPC client

### Overall
- Create makefile for testing and building
- Testcontainers implementation
