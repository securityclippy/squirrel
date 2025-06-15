package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"reminder-service/internal/models"
	"reminder-service/internal/services"
)

type AuthMiddleware struct {
	userService *services.UserService
	domain      string
	audience    string
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Scope string `json:"scope"`
}

type contextKey string

const (
	UserContextKey   contextKey = "user"
	APIKeyContextKey contextKey = "apikey"
)

func NewAuthMiddleware(userService *services.UserService) *AuthMiddleware {
	domain := os.Getenv("AUTH0_DOMAIN")
	audience := os.Getenv("AUTH0_AUDIENCE")
	
	if domain == "" || audience == "" {
		// Log warning but don't panic in development
		fmt.Printf("Warning: AUTH0_DOMAIN and AUTH0_AUDIENCE environment variables not set\n")
		domain = "dev.auth0.com"
		audience = "dev-api"
	}

	return &AuthMiddleware{
		userService: userService,
		domain:      domain,
		audience:    audience,
	}
}

func (m *AuthMiddleware) AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		token, err := m.validateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Get or create user based on Auth0 ID
		user, err := m.userService.GetOrCreateUserFromAuth0(r.Context(), claims.Subject)
		if err != nil {
			http.Error(w, "Failed to authenticate user", http.StatusInternalServerError)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) AuthenticateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := extractAPIKey(r)
		if apiKey == "" {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		user, keyInfo, err := m.userService.ValidateAPIKey(r.Context(), apiKey)
		if err != nil {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		// Add user and API key info to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, APIKeyContextKey, keyInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) AuthenticateAny(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AuthenticateAny: Processing request to %s", r.URL.Path)
		
		// Try JWT first
		tokenString := extractToken(r)
		if tokenString != "" {
			log.Printf("AuthenticateAny: Found JWT token, attempting validation")
			token, err := m.validateJWT(tokenString)
			if err == nil {
				if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
					log.Printf("AuthenticateAny: JWT validation successful for subject: %s", claims.Subject)
					user, err := m.userService.GetOrCreateUserFromAuth0(r.Context(), claims.Subject)
					if err == nil {
						log.Printf("AuthenticateAny: User authenticated successfully via JWT")
						ctx := context.WithValue(r.Context(), UserContextKey, user)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					} else {
						log.Printf("AuthenticateAny: Failed to get/create user from Auth0: %v", err)
					}
				} else {
					log.Printf("AuthenticateAny: Invalid JWT claims or token")
				}
			} else {
				log.Printf("AuthenticateAny: JWT validation failed: %v", err)
			}
		} else {
			log.Printf("AuthenticateAny: No JWT token found")
		}

		// Try API key
		apiKey := extractAPIKey(r)
		if apiKey != "" {
			log.Printf("AuthenticateAny: Found API key, attempting validation")
			user, keyInfo, err := m.userService.ValidateAPIKey(r.Context(), apiKey)
			if err == nil {
				log.Printf("AuthenticateAny: API key validation successful")
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				ctx = context.WithValue(ctx, APIKeyContextKey, keyInfo)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				log.Printf("AuthenticateAny: API key validation failed: %v", err)
			}
		} else {
			log.Printf("AuthenticateAny: No API key found")
		}

		log.Printf("AuthenticateAny: Authentication failed - no valid JWT or API key")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
	})
}

func (m *AuthMiddleware) validateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the key ID from the token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in token header")
		}

		// Get the JWKS URL
		jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", m.domain)
		
		// Fetch and get the key
		keySet, err := jwk.Fetch(context.Background(), jwksURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
		}

		key, found := keySet.LookupKeyID(kid)
		if !found {
			return nil, fmt.Errorf("key not found in JWKS")
		}

		var publicKey interface{}
		if err := key.Raw(&publicKey); err != nil {
			return nil, fmt.Errorf("failed to get raw key: %w", err)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Validate audience
	found := false
	for _, aud := range claims.Audience {
		if aud == m.audience {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("invalid audience")
	}

	// Validate issuer
	expectedIssuer := fmt.Sprintf("https://%s/", m.domain)
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	return token, nil
}


func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func extractAPIKey(r *http.Request) string {
	// Check Authorization header with API key format
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "apikey" {
			return parts[1]
		}
	}

	// Check X-API-Key header
	return r.Header.Get("X-API-Key")
}

func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	return user, ok
}

func GetAPIKeyFromContext(ctx context.Context) (*models.APIKey, bool) {
	apiKey, ok := ctx.Value(APIKeyContextKey).(*models.APIKey)
	return apiKey, ok
}