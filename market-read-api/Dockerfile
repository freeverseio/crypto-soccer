FROM node:10-alpine AS build-env

WORKDIR /app
COPY . /app
RUN npm ci --only=production
EXPOSE 5000
RUN chown -R node:node /app
USER node

FROM gcr.io/distroless/nodejs:10
COPY --from=build-env /app /app
WORKDIR /app
CMD ["src/index.js"]
