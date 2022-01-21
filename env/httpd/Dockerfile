FROM nginx:alpine

RUN apt-get update && apt-get install -y init && apt-get install -y procps && apt-get clean all && apt install --reinstall awscli -y && apt install -y net-tools

WORKDIR /var/src

COPY ./nginx.conf /etc/nginx/nginx.conf

EXPOSE 8888

CMD ["/bin/bash", "nginx", "-g", "daemon off;"]

