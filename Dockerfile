FROM alpine
MAINTAINER jszls65@qq.com
CMD ["/bin/sh"]
WORKDIR /app
COPY ../devman /app
COPY ../www /app/www
COPY script /app/script
COPY ../config /app/config
RUN chmod +x /app/devman
EXPOSE 8559
# ENTRYPOINT ["/app/devman", "&"]
CMD ["/app/devman", "&"]