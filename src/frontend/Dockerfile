FROM node:18.18.0-alpine3.18 AS build

WORKDIR /app
COPY ["package.json", "package-lock.json", "./"]
RUN ["npm", "ci"]
COPY . .
RUN ["npx", "ng", "build"]

FROM nginx:1.25-alpine3.18
RUN touch /var/run/nginx.pid && \
  chown -R nginx:nginx /var/cache/nginx /var/run/nginx.pid &&\
  chown -R nginx:nginx /etc/nginx/conf.d
USER nginx
COPY --chown=nginx:nginx ./nginx/default.conf.template /etc/nginx/templates/default.conf.template
COPY --chown=nginx:nginx --from=build /app/dist/frontend/browser /var/www/html/
