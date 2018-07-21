FROM alpine:3.7  

RUN apk --no-cache add ca-certificates \
    && addgroup -S camel \
    && adduser -S -g app camel

WORKDIR /home/camel
COPY ./camel-go /home/camel/camel-go

RUN chown -R camel:camel ./camel-go

ENTRYPOINT [ "./camel-go" ]
