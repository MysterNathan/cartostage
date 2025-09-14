# nginx/Dockerfile
FROM nginx:alpine

# Installer envsubst si pas présent
RUN apk add --no-cache gettext

# Copier le template
COPY nginx.conf.template /etc/nginx/templates/default.conf.template

# Script de démarrage
COPY start-nginx.sh /start-nginx.sh
RUN chmod +x /start-nginx.sh

EXPOSE 80 443

CMD ["/start-nginx.sh"]
