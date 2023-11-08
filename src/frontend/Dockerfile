FROM node:18.18.0-alpine3.18 AS build

WORKDIR /app
COPY ["package.json", "package-lock.json", "."]
RUN ["npm", "ci"]
COPY . .
RUN ["npx", "ng", "build"]

FROM node:18.18.0-alpine3.18
WORKDIR /app
RUN npm install -g serve@14.2.1
COPY --from=build /app/dist/frontend/browser .
CMD ["serve", "-n", "--no-port-switching", "-p", "4000"]
