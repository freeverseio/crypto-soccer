# Stage 1
FROM node:12.2.0-alpine AS react-build
ARG REACT_APP_API="https://dashboard.gorengine.com/api" 
WORKDIR /app
COPY . ./
RUN npm install --silent
RUN REACT_APP_API=${REACT_APP_API} npm run-script build

# Stage 2 - the production environment
FROM nginx:alpine
# COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=react-build /app/build /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]