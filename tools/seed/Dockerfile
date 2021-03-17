FROM node:14.16.0-alpine

WORKDIR /app

COPY package.json yarn.lock /app/

RUN yarn --prod

COPY . .

CMD ["node", "index.js"]




