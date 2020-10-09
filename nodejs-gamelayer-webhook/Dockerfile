FROM node:10-alpine AS build-env

ENV PORT ${PORT}
WORKDIR /app
COPY . /app
RUN npm ci --only=production
EXPOSE ${PORT}
RUN chown -R node:node /app
USER node

FROM gcr.io/distroless/nodejs:10
COPY --from=build-env /app /app
WORKDIR /app
CMD ["src/index.js"]
