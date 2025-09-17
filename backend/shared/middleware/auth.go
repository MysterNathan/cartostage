package middleware

import (
	"net/http"
	"strings"

	sharedContext "shared/context"
	"shared/models"
	"shared/services"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Middleware de base - vérifie juste le token
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Ajouter les claims au contexte
		ctx := sharedContext.SetClaimsInContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware avec vérification de rôle (version simplifiée)
func (m *AuthMiddleware) RequireRole(role models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := sharedContext.GetClaimsFromContext(r.Context())

			if claims.Role != role {
				http.Error(w, "Insufficient role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

// Middleware pour vérifier les permissions sur une ressource
func (m *AuthMiddleware) RequirePermission(resourceType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sharedContext.HasPermissionFor(r.Context(), resourceType, nil) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

// Middleware pour admins seulement
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireRole(models.RoleAdmin)(next)
}

// Middleware pour enseignants et admins
func (m *AuthMiddleware) RequireTeacherOrAdmin(next http.Handler) http.Handler {
	return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := sharedContext.GetClaimsFromContext(r.Context())

		if claims.Role != models.RoleTeacher && claims.Role != models.RoleAdmin {
			http.Error(w, "Teacher or Admin role required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}))
}

func (m *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	return r.URL.Query().Get("token")
}
