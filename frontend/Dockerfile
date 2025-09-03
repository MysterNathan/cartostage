# Multi-stage build pour optimiser la taille de l'image
FROM node:18-alpine AS base

# Installer les dépendances uniquement quand nécessaire
FROM base AS deps
# Vérifier https://github.com/nodejs/docker-node/tree/b4117f9333da4138b03a546ec926ef50a31506c3#nodealpine pour comprendre pourquoi libc6-compat pourrait être nécessaire.
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Installer les dépendances basées sur le gestionnaire de paquets préféré
COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* ./
RUN \
  if [ -f yarn.lock ]; then yarn --frozen-lockfile; \
  elif [ -f package-lock.json ]; then npm ci; \
  elif [ -f pnpm-lock.yaml ]; then corepack enable pnpm && pnpm i --frozen-lockfile; \
  else echo "Lockfile not found." && exit 1; \
  fi

# Rebuild du code source uniquement quand nécessaire
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Variables d'environnement pour le build
# Next.js collecte des données d'utilisation anonymes sur l'utilisation générale.
# Apprenez-en plus ici: https://nextjs.org/telemetry
# Décommenter la ligne suivante si vous voulez désactiver la télémétrie pendant le build.
ENV NEXT_TELEMETRY_DISABLED 1

# Build de l'application
RUN \
  if [ -f yarn.lock ]; then yarn run build; \
  elif [ -f package-lock.json ]; then npm run build; \
  elif [ -f pnpm-lock.yaml ]; then corepack enable pnpm && pnpm run build; \
  else echo "Lockfile not found." && exit 1; \
  fi

# Image de production
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
# Décommenter la ligne suivante si vous voulez désactiver la télémétrie pendant l'exécution.
ENV NEXT_TELEMETRY_DISABLED 1

# Créer un utilisateur non-root pour exécuter l'application
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copier les fichiers nécessaires depuis le builder
COPY --from=builder /app/public ./public

# Set the correct permission for prerender cache
RUN mkdir .next
RUN chown nextjs:nodejs .next

# Copier automatiquement les fichiers de sortie en utilisant les permissions appropriées
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

# Exposer le port 80
EXPOSE 80

ENV PORT 80
ENV HOSTNAME "0.0.0.0"

# Démarrer l'application
CMD ["node", "server.js"]
