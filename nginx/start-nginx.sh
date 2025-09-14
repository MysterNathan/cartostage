#!/bin/sh
# nginx/start-nginx.sh

# Remplacer les variables d'environnement dans le template
envsubst '${DOMAIN_NAME} ${SSL_CERT_PATH} ${SSL_KEY_PATH}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf

# Démarrer nginx
nginx -g 'daemon off;'
