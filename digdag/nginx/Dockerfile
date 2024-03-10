FROM nginx:stable-alpine

ARG APP_HOST=scheduler:65432

COPY nginx.conf /tmp/

RUN set -eux \
 && envsubst '$$APP_HOST' < /tmp/nginx.conf > /etc/nginx/nginx.conf \
 && :

CMD ["nginx", "-g", "daemon off;"]
