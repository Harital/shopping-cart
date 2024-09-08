# Shopping cart mockup

# Usage
- Checkout git project
- Have installed docker and docker-compose
- Open the project in vs code. Install dev container extensions
- Reopen in container. It will create all the necessay dockers

### Database initialization

The application user and the table where to store the data needs to be created. There is an script called in scripts/sql/databaseInitialization.sql that performs the required operations

```bash
mysql -v -h 127.0.0.1 -P 45478 -u root -proot
```

For production environment, further db migration tools may be considered, like flyway, that can help with the database structure migration (https://github.com/flyway/flyway)

# Directory structure

# Used libraries
gin
resty
gomock 


# Other considerations
Go generate for mocks
Errors only written in the handler with errors wrapped, so we are avoiding n lines of error for a single failure
Not responding very verbose erros to the customer for security reasons. Details are dumped in the logs.
Mocks are generated with 
```bash
go generate ./...
```

## Tests
Test structure
Gherkin notation

## potential imprivements
Authentication
Request id to track queries. This request id could be returned in the error response and in the logs in order to match the erroneous queries.
Modification of log levels
Add secrets manager for storing secret values (users, passwords, etc)
Add env variables depending on stack (dev, stg, pro) in order to tweak values as db host, reserve host....

