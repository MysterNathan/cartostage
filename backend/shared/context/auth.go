package context

import (
	"context"
	"fmt"
	"shared/models"
)

type contextKey string

const (
	userClaimsKey contextKey = "userClaims"
)

// SetUserClaims ajoute les claims de l'utilisateur au contexte
func SetUserClaims(ctx context.Context, claims *models.CustomClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

// Alias pour compatibilité avec le middleware
func SetClaimsInContext(ctx context.Context, claims *models.CustomClaims) context.Context {
	return SetUserClaims(ctx, claims)
}

// GetUserClaims récupère les claims de l'utilisateur depuis le contexte
func GetUserClaims(ctx context.Context) *models.CustomClaims {
	claims, ok := ctx.Value(userClaimsKey).(*models.CustomClaims)
	if !ok {
		return nil
	}
	return claims
}

// Alias pour compatibilité avec le middleware
func GetClaimsFromContext(ctx context.Context) *models.CustomClaims {
	return GetUserClaims(ctx)
}

// GetUserRoleFromContext récupère le rôle de l'utilisateur depuis le contexte
func GetUserRoleFromContext(ctx context.Context) models.UserRole {
	claims := GetUserClaims(ctx)
	if claims == nil {
		return ""
	}
	return claims.Role
}

// HasPermissionFor vérifie les permissions sur une ressource
func HasPermissionFor(ctx context.Context, resourceType string, resource interface{}) bool {
	claims := GetUserClaims(ctx)
	if claims == nil {
		return false
	}

	// Les admins ont toutes les permissions
	if claims.IsAdmin() {
		return true
	}

	// Logique simple : chaque rôle peut accéder à ses propres ressources
	switch resourceType {
	case "users":
		// Les utilisateurs peuvent voir les utilisateurs de leur propre rôle
		return true
	case "stages":
		// Tous les rôles authentifiés peuvent voir les stages
		return true
	case "filieres":
		// Tous les rôles authentifiés peuvent voir les filières
		return true
	default:
		return false
	}
}

// RequireAuth vérifie que l'utilisateur est authentifié
func RequireAuth(ctx context.Context) (*models.CustomClaims, error) {
	claims := GetUserClaims(ctx)
	if claims == nil {
		return nil, fmt.Errorf("authentication required")
	}
	return claims, nil
}

// RequireRole vérifie que l'utilisateur a un rôle spécifique
func RequireRole(ctx context.Context, role models.UserRole) (*models.CustomClaims, error) {
	claims, err := RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	if claims.Role != role {
		return nil, fmt.Errorf("insufficient privileges: requires %s role", role)
	}

	return claims, nil
}

// RequireAdmin vérifie que l'utilisateur est admin
func RequireAdmin(ctx context.Context) (*models.CustomClaims, error) {
	return RequireRole(ctx, models.RoleAdmin)
}

// RequireTeacherOrAdmin vérifie que l'utilisateur est enseignant ou admin
func RequireTeacherOrAdmin(ctx context.Context) (*models.CustomClaims, error) {
	claims, err := RequireAuth(ctx)
	if err != nil {
		return nil, err
	}

	if !claims.IsAdmin() && !claims.IsTeacher() {
		return nil, fmt.Errorf("insufficient privileges: requires teacher or admin role")
	}

	return claims, nil
}
