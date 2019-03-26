FROM postgres:11
COPY create-dragonfly-testing-db.sql /docker-entrypoint-initdb.d/