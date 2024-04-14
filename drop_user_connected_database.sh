#!/bin/bash

# Variables
USER_TO_DROP="newuser"
DB_NAME="connection_pool_test"

PGPASSWORD=postgres_password psql -U postgres -h localhost -d postgres -c "
-- Revoke all privileges in the specified database
REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM $USER_TO_DROP;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM $USER_TO_DROP;
REVOKE ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public FROM $USER_TO_DROP;
REVOKE USAGE ON SCHEMA public FROM $USER_TO_DROP;
REVOKE CONNECT ON DATABASE $DB_NAME FROM $USER_TO_DROP;

-- Disconnect any existing sessions
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = '$DB_NAME'
  AND pg_stat_activity.usename = '$USER_TO_DROP';

-- Drop the user
DROP ROLE $USER_TO_DROP;
"

echo "User $USER_TO_DROP has been dropped."

