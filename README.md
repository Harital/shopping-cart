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

# Other considerations
Go generate for mocks

## Tests
Test structure
Gherkin notation