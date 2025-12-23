FROM oven/bun:1-alpine AS development

WORKDIR /app

COPY package.json bun.lockb* ./

RUN bun install

COPY . .

EXPOSE 5173

CMD ["bun", "run", "dev", "--", "--host"]

FROM oven/bun:1-alpine AS builder

WORKDIR /app

COPY package.json bun.lockb* ./

RUN bun install

COPY . .

RUN bun run build

FROM nginx:alpine AS production

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]
