FROM node:14.15.4-alpine

WORKDIR /app

COPY package.json yarn.lock /app/

RUN yarn --prod

CMD ["node", "index.js"]




