# Go payment-api service

# Installation

## Clone the repo
```bash
git clone git@github.com:akorol1998/headway-test.git
```

### Setup configurations
To configure this service you must have `docker` installed, and be able to run `Make` rules.


# How to run
Make sure you are in the root directory of your project
### First, build the image.
```bash
$ make build
```
### Then you can run it in isolation
It will bring up the database and service.

```bash
$ make up
```

This will also output logs in the format of:
`{"ID": <product-id>, "Name": <product-name>}`
You should use `<product-id>` in order to perform api calls to the service. *Note there is product which name is `InvalidProvider`, you should use it's `product-id` in order to achieve Unexpected bahavior on the service*

### Test the via Postman/Curl
Service is at `0.0.0.0:8080`.
```bash
curl -v  http://localhost:8080/api/v1/payment/url?productID=<product-ID>
```
## Running tests
To run unit tests:
```bash
make test
```

To run integration tests:
```bash
make integration-test
```
