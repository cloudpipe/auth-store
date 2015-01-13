# auth-store

*locally administered account management for cloudpipe*

[![Build Status](https://travis-ci.org/cloudpipe/auth-store.svg?branch=master)](https://travis-ci.org/cloudpipe/auth-store)

This is the default implementation of the [cloudpipe authentication backend protocol](https://github.com/cloudpipe/cloudpipe/wiki/Authentication) that stores account data in MongoDB, most likely the same instance that's used to host other internal Cloudpipe data. It exposes additional API endpoints to permit account creation and API key generation.
