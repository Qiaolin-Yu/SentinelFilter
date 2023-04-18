# Sentinel: Bloom Filter-as-a-Service

This is a simple HTTP server that provides an interface to use a Bloom filter. A Bloom filter is a probabilistic data structure that allows you to quickly check if an item is a member of a set or not. It now also supports the gRPC protocol.

The server supports the following HTTP endpoints:

- `GET /?key=<value>`: Returns a boolean indicating whether the given value may be in the Bloom filter or not.
- `POST /`: Adds a value to the Bloom filter. The value should be provided as a JSON payload with the key "key".

## How to use
Note that the server uses a mutex to ensure that the Bloom filter can be used concurrently by multiple HTTP clients.
To use the HTTP interface for the Bloom filter, you can follow these steps:

### 1. Run the server with the following command:
```bash
go run ./...
```
By default, the server listens on port 7070.


### 2. Make a GET request to check if a value is in the Bloom filter:
```bash
curl http://localhost:7070/?key=value
```
If the value is in the Bloom filter, the server will return "true". Otherwise, it will return "false".

### 3. Add a value to the Bloom filter using a POST request:
```bash
curl -d '{"key":"value"}' -H "Content-Type: application/json" -X POST http://localhost:7070/
```
