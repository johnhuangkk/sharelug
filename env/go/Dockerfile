FROM golang:latest

ENV GOLANG_VERSION=1.14.4

RUN apt-get update && apt-get install -y init && apt-get clean all && apt install --reinstall awscli -y && apt-get install -y wkhtmltopdf

RUN mkdir -p /opt/jdk && \
    curl https://download.java.net/java/GA/jdk11/13/GPL/openjdk-11.0.1_linux-x64_bin.tar.gz -o /opt/jdk-11.0.1.tgz && \
    tar -xvf /opt/jdk-11.0.1.tgz --strip-components=1 -C /opt/jdk

ENV JAVA_HOME /opt/jdk

ENV PATH "$PATH:$JAVA_HOME/bin"

ENV APP_NAME sharelug

RUN mkdir -p /go/src/sharelug

WORKDIR /go/src/sharelug

VOLUME /go/src/sharelug

EXPOSE 8001

CMD ["bash"]