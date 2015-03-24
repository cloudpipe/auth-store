# auth-store

*locally administered account management for cloudpipe*

[![Build Status](https://travis-ci.org/cloudpipe/auth-store.svg?branch=master)](https://travis-ci.org/cloudpipe/auth-store)

This is the default implementation of the [cloudpipe authentication backend protocol](https://github.com/cloudpipe/cloudpipe/wiki/Authentication) that stores account data in MongoDB, most likely the same instance that's used to host other internal Cloudpipe data. It exposes additional API endpoints to permit account creation and API key generation.

## Getting Started

 1. Install [Docker](https://docs.docker.com/installation/mac/) on your system.
 2. Install [Compose](https://docs.docker.com/compose/install/).
 3. Use `script/genkeys` to generate self-signed TLS keypairs in `certificates/`.
 4. Run `docker-compose build && docker-compose up -d` to build and launch everything locally.

To run the tests, use `script/test`. You can also use `script/mongo` to connect to your local MongoDB database.

### Using the API

Once it's up and running, you can use `curl` to interact the auth API. Here are a few examples:

```bash
# If you're on a Mac and using boot2docker. Otherwise, you can use "localhost".
DOCKER=$(boot2docker ip 2>/dev/null)

# Create a new account.
curl -k -i -X POST https://${DOCKER}:9000/v1/accounts -d 'accountName=me%40gmail.com&password=shhh'

# Generate a new API key.
curl -k -i -X POST https://${DOCKER}:9000/v1/keys -d 'accountName=me%40gmail.com&password=shhh'
```

### API Documentation

Current API documentation may be found [in the `docs/` directory](docs/api.md).
