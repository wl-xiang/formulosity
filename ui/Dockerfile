FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --no-audit --no-fund

COPY . .

# Requests are proxied by nginx to the api service, so the UI is same-origin.
ARG VITE_API_ADDR=/api
ENV VITE_API_ADDR=${VITE_API_ADDR}

RUN npm run build && npm cache clean --force && rm -rf /root/.npm

FROM nginx:1.27-alpine

RUN rm -rf /usr/share/nginx/html/* /var/cache/apk/*

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
