FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD root
ENV MYSQL_USER user
ENV MYSQL_PASSWORD password

COPY ./scripts/sql/init.sql /docker-entrypoint-initdb.d/

EXPOSE 3306
