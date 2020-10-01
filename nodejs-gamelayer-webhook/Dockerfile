FROM node:12-alpine 

ENV PORT ${PORT}
WORKDIR /app
COPY . /app
RUN npm ci --only=production
EXPOSE ${PORT}
RUN chown -R node:node /app
USER node
ENTRYPOINT [ "npm", "start", "--" ]
