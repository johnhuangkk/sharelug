FROM node:14.15-alpine

RUN apk update && apk add bash

RUN npm install pm2 -g

WORKDIR /angular/sharelug

VOLUME /angular/sharelug

EXPOSE 4000

CMD ["/bin/bash"]
