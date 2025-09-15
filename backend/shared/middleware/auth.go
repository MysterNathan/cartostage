package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"shared/services"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

type AuthMiddleware struct {
	jwtService *services.JWTService
}

func NewAuthMiddleware(jwtService *services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Middleware de base - vérifie juste le token
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraire le token
		token := m.extractToken(r)
		if token == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Valider le token
		claims, err := m.jwtService.ValidateAndParseToken(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Ajouter les claims au contexte
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware avec vérification de scope
func (m *AuthMiddleware) RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())

			if !claims.HasScope(scope) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

// Middleware avec vérification de rôle
func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())

			if !claims.IsRole(role) {
				http.Error(w, "Insufficient role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

// Middleware avec vérifications multiples
func (m *AuthMiddleware) RequirePermissions(requiredRole string, requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())

			// Vérifier le rôle
			if requiredRole != "" && !claims.IsRole(requiredRole) {
				http.Error(w, "Insufficient role", http.StatusForbidden)
				return
			}

			// Vérifier le scope
			if requiredScope != "" && !claims.HasScope(requiredScope) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

func (m *AuthMiddleware) extractToken(r *http.Request) string {
	// Vérifier le header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Fallback sur query param (pour websockets par exemple)
	return r.URL.Query().Get("token")
}

// Helper pour récupérer les claims du contexte
func GetClaimsFromContext(ctx context.Context) *services.CustomClaims {
	if claims, ok := ctx.Value(ClaimsContextKey).(*services.CustomClaims); ok {
		return claims
	}
	return nil
}

// Helper pour récupérer l'ID de l'entreprise depuis le contexte
func GetEnterpriseIDFromContext(ctx context.Context) (int, bool) {
	claims := GetClaimsFromContext(ctx)
	if claims == nil {
		log.Println("DEBUG: No claims found in context")
		return 0, false
	}

	log.Printf("DEBUG: User role: %s, EntityID: %v", claims.Role, claims.EntityID)

	// Vérifier que c'est bien un utilisateur entreprise
	if claims.Role != "enterprise" {
		log.Printf("DEBUG: User is not enterprise role, got: %s", claims.Role)
		return 0, false
	}

	// EntityID pour un user entreprise = ID de l'entreprise
	if claims.EntityID == nil {
		log.Println("DEBUG: EntityID is nil")
		return 0, false
	}

	log.Printf("DEBUG: Returning enterprise ID: %d", *claims.EntityID)
	return *claims.EntityID, true
}

// Helper pour récupérer l'ID utilisateur depuis le contexte
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	claims := GetClaimsFromContext(ctx)
	if claims == nil {
		return 0, false
	}
	return claims.UserID, true
}

// Helper pour récupérer le rôle depuis le contexte
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	claims := GetClaimsFromContext(ctx)
	if claims == nil {
		return "", false
	}
	return claims.Role, true
}
