#!/bin/bash

DB_NAME="connection_pool_test"
NEW_USER="new_connection_pool_test"
NEW_PASSWORD="satapril2024$"


echo "Database $DB_NAME created successfully."
# Connect again to set up the user and privileges
PGPASSWORD=postgres_password psql -U postgres -h localhost -d $DB_NAME -c "
CREATE ROLE $NEW_USER WITH LOGIN PASSWORD '$NEW_PASSWORD';
GRANT CONNECT ON DATABASE $DB_NAME TO $NEW_USER;
GRANT USAGE ON SCHEMA public TO $NEW_USER;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO $NEW_USER;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO $NEW_USER;
"
echo "User $NEW_USER created and privileges granted on database $DB_NAME"


