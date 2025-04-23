# Builder stage
FROM node:22 AS builder

WORKDIR /usr/src/sf-fe

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Production stage
FROM nginx:stable-alpine
COPY --from=builder /usr/src/sf-fe/build /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
