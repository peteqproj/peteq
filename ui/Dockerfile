FROM node:14.14.0-alpine3.12 AS dev

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

RUN npm install -g create-react-app

COPY ./package.json yarn.lock /app/

RUN yarn --silent

COPY . /app

ENV REACT_APP_API=http://a58f5419b766e43269d5eddad79cf2f7-1019754084.us-east-1.elb.amazonaws.com:8080

RUN yarn build

FROM nginx:1.17.8-alpine

COPY --from=dev /app/build /usr/share/nginx/html

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx/nginx.conf /etc/nginx/conf.d

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]