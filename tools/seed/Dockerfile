FROM node:14.15.4-alpine

WORKDIR /app

COPY package.json yarn.lock /app/

RUN yarn --prod

COPY . .

CMD ["node", "index.js"]




