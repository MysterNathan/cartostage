package services_test

import (
	"auth/internal/repositories"
	"auth/internal/services"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"shared/models"
	sharedServices "shared/services"
)

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("Impossible de hasher le mot de passe : %v", err)
	}
	return string(hash)
}

func buildAuthService(t *testing.T, fakeRepo *repositories.MockAuthRepository) *services.AuthService {
	t.Helper()
	jwtService := sharedServices.NewJWTService("secret-de-test", "default")
	return services.NewAuthService(fakeRepo, jwtService)
}

func TestLogin_Success(t *testing.T) {
	fakeRepo := repositories.NewMockAuthRepository()
	fakeRepo.Users["alice"] = &models.User{
		ID:           1,
		Username:     "alice",
		PasswordHash: hashPassword(t, "monMotDePasse"),
		Role:         "user",
	}

	authService := buildAuthService(t, fakeRepo)

	req := &services.LoginRequest{
		Username: "alice",
		Password: "monMotDePasse",
	}

	// ACT - Appeler la fonction à tester
	response, err := authService.Login(req)

	// ASSERT - Vérifier les résultats
	if err != nil {
		t.Fatalf("Erreur inattendue : %v", err)
	}
	if response == nil {
		t.Fatal("La réponse ne devrait pas être nil")
	}
	if response.Token == "" {
		t.Error("Le token ne devrait pas être vide")
	}
	if response.User.PasswordHash != "" {
		t.Error("Le hash du mot de passe devrait être effacé dans la réponse")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	fakeRepo := repositories.NewMockAuthRepository()
	fakeRepo.Users["alice"] = &models.User{
		Username:     "alice",
		PasswordHash: hashPassword(t, "bonMotDePasse"),
	}

	authService := buildAuthService(t, fakeRepo)

	_, err := authService.Login(&services.LoginRequest{
		Username: "alice",
		Password: "mauvaisMotDePasse",
	})

	if err == nil {
		t.Fatal("Une erreur était attendue avec un mauvais mot de passe")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	fakeRepo := repositories.NewMockAuthRepository() // Vide, aucun user
	authService := buildAuthService(t, fakeRepo)

	_, err := authService.Login(&services.LoginRequest{
		Username: "inexistant",
		Password: "nimportequoi",
	})

	if err == nil {
		t.Fatal("Une erreur était attendue pour un user inexistant")
	}
}

func TestLogin_EmptyCredentials(t *testing.T) {
	// Table driven test - tester plusieurs cas similaires proprement
	tests := []struct {
		name     string
		username string
		password string
	}{
		{"username vide", "", "monMotDePasse"},
		{"password vide", "alice", ""},
		{"les deux vides", "", ""},
	}

	fakeRepo := repositories.NewMockAuthRepository()
	authService := buildAuthService(t, fakeRepo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := authService.Login(&services.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			})

			if err == nil {
				t.Errorf("Cas '%s' : une erreur était attendue", tt.name)
			}
		})
	}
}
