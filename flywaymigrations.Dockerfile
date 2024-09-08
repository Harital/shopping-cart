FROM flyway/flyway:10.11-alpine

USER 1002

COPY ./flyway/flyway_migration.sh /
COPY ./flyway/sql /flyway/sql
COPY ./flyway/flyway.toml /flyway/conf/

# Execute the script to init the container
ENTRYPOINT ["sh", "-c", "/flyway_migration.sh"]
