# Shopping cart mockup

### Introduction
This application consists about the shopping cart exercise. 

As requested, the application implements the endpoint /items with GET and POST methods.

- Get method gets the items of the cart
- Post method adds items to the cart

For simplicity, there is only 1 cart, so the get and post methods always access the same cart. The endpoints are unauthenticated

More details about the endpoints syntax can be found in api/shopping-cart-api.yaml

## Usage
- Checkout git project
- Have installed docker and docker-compose (I have tested it in linux)
- Open the project in vs code. Install dev container extension.
- Reopen in container. It will create all the necessay dockers and open the dev environment
- Open bash
- Execute 
```
go build -o shopping-cart cmd/main.go && ./shopping-cart

```

The exported port is 8080. Postman can be used to access and exercise the endpoints.

## Database initialization

The application stores the current cart status in a "cartItem" table in a mysql database.

The first time that the app is run, the application user and the table where to store the data need to be created. There is an script called in scripts/sql/databaseInitialization.sql that performs the required operations.

```bash
mysql -v -h 127.0.0.1 -P 45478 -u root -proot < < scripts/sql/databaseInitialization.sql 
```

For production environment, further db migration tools may be considered, like flyway, that can help with the database structure migration (https://github.com/flyway/flyway). These tools keep the database structure up to date and can be tracked in github via update sscripts.

## Architecture and directory structure

The solution follows hexagonal architeture and go standard directory patterns. This helps with the code mantainability and code decoupling.

All the requests are handled by handlers (in some literature are called controllers). These handlers only accept requests check basic parameters and pass the request to the service, the core application. If any other protocol should be supported in the future (I.E flatbuffers, raw sockets), only the handler should be changed.

The service layer performs the business logic. 

The repo layer stores the data in persistent storage: a database, a bucket, cache, hard drive, etc.

#### Directory structure
- .devcontainer: files for the dev container environment. Docker compose and Docker file for the dev environment.
- .vscode: common vscode settings. It helps with the dev environment harmonization.
- api: openAPI files (former swagger) with the interfacce definition
- cmd: main app
- internal: go packages internal to the app. Not to be exported.
  - adapters: adapters that the core use to perform its operation
    - handlers: http handlers. May other handlers be added, here is the place
    - repositories: persistency modules.
  - core: where the core application lives. 
    - services: business logic
    - mocks: mocks of the interfaces
    - ports: ports or interfaces of the different modules
    - model: data model structure to be shared across the solution
- scripts: useful utility scripts. Here have only included the database initialization script
- root: Docker compose and main Docker file


## Deployability

The solution is built using docker and docker-compose.

The main solution is built using a simple Dockerfile. The dockerfile uses a go docker for build an alpine linux for the run time. This reduces the image footprint

The dev environment mimics this structure in order to get a dev environment as similar as possible to the production envirnoment. The dev environment also uses the docker-compose.yaml file in the root folder. This file sets up a local environment that mimics the production environment. It launches, databases, caches, S3, etc. This enables the local development and decouples it from the cloud.


## Used libraries

I've used some libraries to help with the development
- gin: http handler library. Handles http requests and allows middlewares (although they are not implemented in this solution)
- resty: http client. Very compact code. It helps with creating the requests, the transmission and even the marshal/unmarshal of the data
- gomock: for unit tests. Great mocking utility
- sqlbuilder: eases the sql creation. It has simple ORM capabilities. A more powerful ORM could be used, if needed.


## Other considerations

All mocks are automatically generated from the ports using go generate. In order to refresh them simply type. 
```bash
go generate ./...
```


Errors are only written in the handler with errors wrapped. The intent is to write only a single line of error per query. This avoids a "log mess" with tons of unrelated lines for a single error.

We are Not responding very verbose erros to the customer for security reasons. The details of the errors are dumped in the logs.


## Tests

Every module has its unit tests. The only exception are the main and the initialization modules. 

The tests follow the go table driven unit tests. It makes easier to add a new test case and concentrate every type of call into a single function

All the tests follow the Gherkin notation: GivenXXXX_WhenYYYY_ThenZZZZ

## Potential improvements

Things that have not been developed for the sake of simplicity

#### Authentication
Requests are unauthenticated. In a production environment they should be

#### Request ID
A middleware that injects a unique request id could be implemented. This helps if some queries, for whatever reason, issue several lines of log. It allows to bring together all the logs that belong to the same query. This request id could be returned in the error response for support purposes.


#### Modification of log levels
The default log level could be changed in run time for debugging purposes. An ad-hoc endpoint cloud be developed for this purpose.

#### Secrets manager
Add secrets manager for storing secret values (users, passwords, etc)

#### Env variables depending on stack
Add env variables depending on stack (dev, stg, pro) in order to tweak values that change between staks, like db host, reserve host, etc.

#### Kubernetes files
Some files for tweaking kubernetes resouces (depending on the stack) could be added.