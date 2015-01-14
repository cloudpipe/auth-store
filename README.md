# auth-store

*locally administered account management for cloudpipe*

[![Build Status](https://travis-ci.org/cloudpipe/auth-store.svg?branch=master)](https://travis-ci.org/cloudpipe/auth-store)

This is the default implementation of the [cloudpipe authentication backend protocol](https://github.com/cloudpipe/cloudpipe/wiki/Authentication) that stores account data in MongoDB, most likely the same instance that's used to host other internal Cloudpipe data. It exposes additional API endpoints to permit account creation and API key generation.

## Getting Started

 1. Install [Docker](https://docs.docker.com/installation/mac/) on your system.
 2. Install [fig](http://www.fig.sh/install.html).
 3. Use `script/genkeys` to generate a self-signed TLS certificate in `certificates/`.
 4. Run `fig build && fig up -d` to build and launch everything locally.

To run the tests, use `script/test`. You can also use `script/mongo` to connect to your local MongoDB database.

### Using the API

Once it's up and running, you can use `curl` to interact the auth API. Here are a few examples:

```bash
# If you're on a Mac and using boot2docker. Otherwise, you can use "localhost".
DOCKER=$(boot2docker ip 2>/dev/null)

# Create a new account.
curl -k -i -X POST https://${DOCKER}:8001/v1/accounts -d '{"name":"me@gmail.com","password":"shhh"}'

# Generate a new API key.
curl -k -i -X POST https://${DOCKER}:8001/v1/keys -u me@gmail.com:shhh

# Validate an existing key.
curl -k -i "https://${DOCKER}:8001/v1/validate?accountName=me%40gmail.com&apiKey=${KEY}"
```

## API

#### GET /v1/style

Return a descriptive string that [cloudpipe](https://github.com/cloudpipe/cloudpipe) can report to consumers of its API.

*Response*

The string "authstore" as plaintext.

#### GET /v1/validate?accountName={account}&apiKey={key}

Validate an API key against an account.

*Response*

 * **204 No Content:** when the account name and API key are valid.
 * **404 Not Found:** when the API key is not valid or the account does not exist.

#### POST /v1/accounts

Create a new account.

*Request*

```javascript
{
  "name": "", // Requested account name
  "password": "" // Password to use
}
```

*Response*

 * **201 Created:** Account created successfully.
 * **400 Bad Request:** Malformed JSON or incomplete document.
 * **409 Conflict:** Account name already taken.

#### POST /v1/keys

Generate a new API key and associate it with your account.

*Request*

Include a valid account name and password as HTTP basic auth.

*Response*

 * **200 OK:** Key generated successfully. Response body contains the generated API key as plaintext.
 * **401 Unauthorized:** Unable to authenticate with the provided credentials.

#### DELETE /v1/keys?accountName={name}&apiKey={key}

Revoke an API key from your account.

*Response*

 * **204 No Content:** The API key has been successfully revoked.
 * **400 Bad Request:** Request parameters are missing.
 * **401 Unauthorized:** Unrecognized account or API key.
