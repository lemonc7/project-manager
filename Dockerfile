# 前端
FROM node:24-slim AS frontend-builder
WORKDIR /app
RUN corepack enable pnpm
# 安装依赖，缓存优化
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY ./frontend .
RUN pnpm run build

# 后端
FROM golang:latest AS backend-builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o project-manager .

# 运行采用空镜像
FROM scratch
WORKDIR /app
COPY --from=backend-builder /app/project-manager /app/project-manager
COPY --from=frontend-builder /app/dist /app/dist

ENV TZ=Asia/Shanghai
EXPOSE 8080

ENTRYPOINT ["/app/project-manager"]

