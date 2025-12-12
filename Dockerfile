FROM alpine:latest

WORKDIR /app
COPY project-manager .
COPY dist ./dist

RUN chmod +x ./project-manager

ENTRYPOINT ["./project-manager"]

