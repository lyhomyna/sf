FROM nginx:stable

RUN mkdir -p /etc/nginx/certs

COPY ./certs/* /etc/nginx/certs

RUN chown root:root /etc/nginx/certs/* && \
    chmod 600 /etc/nginx/certs/*-key.pem && \
    chmod 644 /etc/nginx/certs/*.pem

COPY nginx.conf /etc/nginx/nginx.conf
