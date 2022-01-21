# 選擇node:alpine
FROM node:15-alpine

RUN apk update

# 指定預設/工作資料夾
WORKDIR /var/www

RUN npm install

EXPOSE 3000

# 指定啟動container後執行命令
CMD ["npm", "start"]
