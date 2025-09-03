# Multi-stage build pour optimiser la taille de l'image
FROM golang:1.25-alpine AS builder

# Installer les dépendances système nécessaires
RUN apk add --no-cache git ca-certificates tzdata

# Créer un utilisateur non-root
RUN adduser -D -g '' appuser

# Définir le répertoire de travail
WORKDIR /build

# Copier les fichiers go mod et go sum
COPY go.mod ./
COPY go.sum ./

# Télécharger les dépendances
RUN go mod download
RUN go mod verify

# Copier le code source
COPY . .

# Build de l'application avec optimisations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app ./cmd/server/

# Image de production minimale
FROM alpine:latest

# Installer ca-certificates pour HTTPS et wget pour health check
RUN apk --no-cache add ca-certificates wget
RUN adduser -D -g '' appuser

WORKDIR /root/

# Copier l'exécutable depuis le builder
COPY --from=builder /build/app .

# Utiliser l'utilisateur non-root
USER appuser

# Exposer le port 8080
EXPOSE 8080

# Variables d'environnement par défaut
ENV PORT=8080
ENV GIN_MODE=release

# Point d'entrée de l'application
CMD ["./app"]
