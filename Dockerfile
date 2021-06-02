# pull official base image
FROM node:14.4.0-alpine

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install app dependencies
COPY package.json ./
COPY package-lock.json ./
RUN npm install --silent && npm install pm2 -g

# add app
COPY . ./
RUN npm run build --silent

# start app
EXPOSE 8000
CMD ["pm2-runtime", "process.yml"]
