FROM alpine:latest

WORKDIR /app
COPY project-manager .
COPY dist ./dist
COPY config.yaml .

RUN chmod +x ./project-manager

ENTRYPOINT ["./project-manager"]

