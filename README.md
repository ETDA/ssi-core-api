
<h1 align="center">
    Core API 
</h1>

<p align="center">
  <a href="#about">About</a> •
  <a href="#development">Development</a> 
</p>

## About

The SSI Core API are the service that interact with blockchain for verifying documents in the system. 

This repository contain 4 services
- abci
- did-api
- vc-status-api
- vc-verify

these 4 services always communicate with each other and easier to develop the services.

However, we run the services separately, by define the environment when start the service.

## Development

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)

#### Start Service

- Copy file `.env.sample` to `.env`
- run `docker-compose up -d`
- you can access the service via `http:localhost:{port}`
- port will be specified in `docker-compose.yml`
- services defualt port:
    - did-api : `8080`
    - vc-status-api : `8082`
    - vc-verify : `8082`

#### After service `tendermint_init` was run
go to .storage/tendermint/config to edit config.toml

- edit `create_empty_blocks` to `false` (default=`true`)
- edit `proxy_app` to `"tcp://abci:26658"` (default=`tcp://127.0.0.1:26658`)


#don't need for version 0.34.11
- edit `[rpc] laddr = "tcp://0.0.0.0:26657"` (default=`tcp://127.0.0.1:26657`)

Then, restart all services again by `docker-compose restart`


