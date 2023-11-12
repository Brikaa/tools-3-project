FROM node:18.18.0-alpine3.18 AS build

WORKDIR /app
COPY ["package.json", "package-lock.json", "./"]
RUN ["npm", "ci"]
COPY . .
RUN ["npx", "ng", "build"]

FROM nginx:1.25-alpine3.18
COPY ./nginx/default.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /app/dist/frontend/browser /var/www/html/
