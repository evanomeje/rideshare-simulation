FROM node:19-alpine

WORKDIR /app


RUN yarn add pg

COPY simulation ./simulation
COPY dbconfig.json .

CMD cd simulation && node main.js