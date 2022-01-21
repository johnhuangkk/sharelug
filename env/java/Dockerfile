FROM openjdk:15-jdk-slim-buster

#Install Open JDK 6
RUN apt-get update && apt-get install -y init && apt-get install -y procps && apt-get clean all && apt install --reinstall awscli -y

WORKDIR /usr/src

COPY ECPayment.jar /usr/src/ECPayment.jar

CMD ["java","-jar","/usr/src/ECPayment.jar"]