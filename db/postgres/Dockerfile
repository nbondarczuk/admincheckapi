FROM postgres

RUN mkdir -p /docker-entrypoint-initdb.d

COPY init_db_test.sh /docker-entrypoint-initdb.d

EXPOSE 5432

