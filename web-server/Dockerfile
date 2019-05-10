FROM node:9.11
RUN mkdir /usr/src/app/
WORKDIR /usr/src/app
RUN mkdir uploads
COPY package.json .
RUN npm install --silent
COPY ./ .
CMD ["npm", "start"]
