FROM node:18.18.0-alpine3.18 AS build

WORKDIR /app
COPY ["package.json", "package-lock.json", "./"]
RUN ["npm", "ci"]
COPY . .
RUN ["npx", "ng", "build"]

FROM nginxinc/nginx-unprivileged:alpine3.18-perl
COPY --chown=nginx:nginx ./nginx/default.conf.template /etc/nginx/templates/default.conf.template
COPY --chown=nginx:nginx --from=build /app/dist/frontend/browser /var/www/html/
