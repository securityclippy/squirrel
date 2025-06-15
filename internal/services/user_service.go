package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"reminder-service/internal/models"
)

type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{pool: pool}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, auth0_id, email, name, picture, email_verified, 
		       created_at, updated_at, last_login_at
		FROM users WHERE id = $1
	`
	
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Auth0ID, &user.Email, &user.Name, &user.Picture,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}

func (s *UserService) GetUserByAuth0ID(ctx context.Context, auth0ID string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, auth0_id, email, name, picture, email_verified, 
		       created_at, updated_at, last_login_at
		FROM users WHERE auth0_id = $1
	`
	
	err := s.pool.QueryRow(ctx, query, auth0ID).Scan(
		&user.ID, &user.Auth0ID, &user.Email, &user.Name, &user.Picture,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	var user models.User
	query := `
		INSERT INTO users (auth0_id, email, name, picture, email_verified)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, auth0_id, email, name, picture, email_verified, created_at, updated_at, last_login_at
	`
	
	err := s.pool.QueryRow(ctx, query, 
		req.Auth0ID, req.Email, req.Name, req.Picture, req.EmailVerified,
	).Scan(
		&user.ID, &user.Auth0ID, &user.Email, &user.Name, &user.Picture,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	
	if req.Picture != nil {
		updates = append(updates, fmt.Sprintf("picture = $%d", argIndex))
		args = append(args, *req.Picture)
		argIndex++
	}
	
	if req.EmailVerified != nil {
		updates = append(updates, fmt.Sprintf("email_verified = $%d", argIndex))
		args = append(args, *req.EmailVerified)
		argIndex++
	}
	
	if req.LastLoginAt != nil {
		updates = append(updates, fmt.Sprintf("last_login_at = $%d", argIndex))
		args = append(args, *req.LastLoginAt)
		argIndex++
	}

	if len(updates) == 0 {
		return s.GetUserByID(ctx, id)
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)
	
	query := fmt.Sprintf(`
		UPDATE users SET %s 
		WHERE id = $%d
		RETURNING id, auth0_id, email, name, picture, email_verified, created_at, updated_at, last_login_at
	`, 
		string(updates[0]), argIndex)

	for i := 1; i < len(updates); i++ {
		query = fmt.Sprintf("%s, %s", query[:len(query)-len(fmt.Sprintf("WHERE id = $%d", argIndex))], updates[i]) + fmt.Sprintf(" WHERE id = $%d", argIndex)
	}

	var user models.User
	err := s.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Auth0ID, &user.Email, &user.Name, &user.Picture,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (s *UserService) GetOrCreateUserFromAuth0(ctx context.Context, auth0ID string) (*models.User, error) {
	user, err := s.GetUserByAuth0ID(ctx, auth0ID)
	if err == nil {
		// Update last login
		now := time.Now()
		_, err = s.UpdateUser(ctx, user.ID, &models.UpdateUserRequest{
			LastLoginAt: &now,
		})
		return user, err
	}

	// TODO: Fetch user info from Auth0 API to create user
	// For now, create a minimal user record
	req := &models.CreateUserRequest{
		Auth0ID:       auth0ID,
		Email:         fmt.Sprintf("%s@example.com", auth0ID), // Placeholder
		Name:          "Unknown User",                          // Placeholder
		EmailVerified: false,
	}

	return s.CreateUser(ctx, req)
}

func (s *UserService) CreateAPIKey(ctx context.Context, userID int, req *models.CreateAPIKeyRequest) (*models.CreateAPIKeyResponse, error) {
	// Generate a random API key
	key, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	// Hash the key for storage
	keyHash, err := hashAPIKey(key)
	if err != nil {
		return nil, err
	}

	// Extract prefix for identification
	keyPrefix := key[:8]

	var apiKey models.APIKey
	query := `
		INSERT INTO api_keys (user_id, name, key_hash, key_prefix, permissions, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, key_hash, key_prefix, permissions, expires_at, last_used_at, is_active, created_at, updated_at
	`

	err = s.pool.QueryRow(ctx, query,
		userID, req.Name, keyHash, keyPrefix, req.Permissions, req.ExpiresAt,
	).Scan(
		&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.KeyHash, &apiKey.KeyPrefix,
		&apiKey.Permissions, &apiKey.ExpiresAt, &apiKey.LastUsedAt, &apiKey.IsActive,
		&apiKey.CreatedAt, &apiKey.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.CreateAPIKeyResponse{
		APIKey: &apiKey,
		Key:    key,
	}, nil
}

func (s *UserService) ValidateAPIKey(ctx context.Context, key string) (*models.User, *models.APIKey, error) {
	// We need to get the API key by prefix first, then verify the hash
	if len(key) < 8 {
		return nil, nil, fmt.Errorf("invalid API key format")
	}
	
	keyPrefix := key[:8]

	var user models.User
	var apiKey models.APIKey
	query := `
		SELECT u.id, u.auth0_id, u.email, u.name, u.picture, u.email_verified, 
		       u.created_at, u.updated_at, u.last_login_at,
		       ak.id, ak.user_id, ak.name, ak.key_hash, ak.key_prefix, ak.permissions,
		       ak.expires_at, ak.last_used_at, ak.is_active, ak.created_at, ak.updated_at
		FROM users u
		JOIN api_keys ak ON u.id = ak.user_id
		WHERE ak.key_prefix = $1 AND ak.is_active = true
		AND (ak.expires_at IS NULL OR ak.expires_at > NOW())
	`

	err := s.pool.QueryRow(ctx, query, keyPrefix).Scan(
		&user.ID, &user.Auth0ID, &user.Email, &user.Name, &user.Picture,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
		&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.KeyHash, &apiKey.KeyPrefix,
		&apiKey.Permissions, &apiKey.ExpiresAt, &apiKey.LastUsedAt, &apiKey.IsActive,
		&apiKey.CreatedAt, &apiKey.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, fmt.Errorf("invalid API key")
		}
		return nil, nil, err
	}

	// Verify the key hash
	err = bcrypt.CompareHashAndPassword([]byte(apiKey.KeyHash), []byte(key))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid API key")
	}

	// Update last used timestamp
	go func() {
		updateQuery := `UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`
		s.pool.Exec(context.Background(), updateQuery, apiKey.ID)
	}()

	return &user, &apiKey, nil
}

func (s *UserService) GetUserAPIKeys(ctx context.Context, userID int) ([]*models.APIKey, error) {
	query := `
		SELECT id, user_id, name, key_hash, key_prefix, permissions, expires_at, 
		       last_used_at, is_active, created_at, updated_at
		FROM api_keys 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []*models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		err := rows.Scan(
			&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.KeyHash, &apiKey.KeyPrefix,
			&apiKey.Permissions, &apiKey.ExpiresAt, &apiKey.LastUsedAt, &apiKey.IsActive,
			&apiKey.CreatedAt, &apiKey.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, &apiKey)
	}

	return apiKeys, nil
}

func (s *UserService) RevokeAPIKey(ctx context.Context, userID int, keyID int) error {
	query := `UPDATE api_keys SET is_active = false WHERE id = $1 AND user_id = $2`
	result, err := s.pool.Exec(ctx, query, keyID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return "sk_" + hex.EncodeToString(bytes), nil
}

func hashAPIKey(key string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}