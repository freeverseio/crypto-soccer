FROM node:10

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /appa
COPY . /app

RUN cd truffle-core && npm install && ./node_modules/.bin/truffle compile
RUN cd nodejs-horizon && npm install
# If you are building your code for production
# RUN npm ci --only=production

EXPOSE 4000

CMD ["node", "nodejs-horizon/src/index.js"]