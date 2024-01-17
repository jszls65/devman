#!/bin/bash
FROM alpine
MAINTAINER jszls65@icloud.com
CMD ["/bin/sh"]
RUN mkdir -p /app/alertman
WORKDIR /app/alertman
COPY ../devman /app/alertman/
COPY ../www /app/alertman/www
COPY script /app/alertman/script
COPY ../config /app/alertman/config
WORKDIR /app/alertman
EXPOSE 8559
ENTRYPOINT ["./devman", "&"]
