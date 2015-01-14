# API Documentation

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

```
accountName={account}&password={password}
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
