FROM scratch

WORKDIR /app
COPY project-manager .
COPY dist ./dist

ENTRYPOINT ["./project-manager"]

