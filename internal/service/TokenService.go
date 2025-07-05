package service

import (
	"test_go_project/internal/auth"
	"test_go_project/internal/common"
	"test_go_project/internal/dto"
	"test_go_project/internal/entity"
	"test_go_project/internal/repository"
	"test_go_project/internal/customErrors"

	"strings"
	"time"
	"errors"

	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

const base64Length = 32

type TokenService struct {
	repo     *repository.TokenRepo
	jwtUtils *auth.JWTUtils
	notifier common.Notifier
}

func NewTokenService(repo *repository.TokenRepo, jwtUtils *auth.JWTUtils, notifier common.Notifier) *TokenService {
	return &TokenService{repo, jwtUtils, notifier}
}

func (s *TokenService) GetNewTokenPair(id uuid.UUID, userAgent string, ip string) (*dto.TokenPair, error) {
	accessJTI := uuid.NewString()
	tokenPair, err := s.generateTokenPair(id, accessJTI)

	if err != nil {
		return nil, customErrors.ErrUnexpectedState
	}

	refreshHash, err := bcrypt.GenerateFromPassword([]byte(tokenPair.RefreshToken), bcrypt.DefaultCost)

	if err != nil {
		return nil, customErrors.ErrUnexpectedState
	}

	err = s.repo.SaveRefreshPair(&entity.RefreshRecord{UserID: id,
		RefreshHash: refreshHash, UserAgent: userAgent, ClientIp: ip, AccessJTI: accessJTI})

	if err != nil {
		return nil, customErrors.ErrDB
	}

	return tokenPair, nil
}

func (s *TokenService) UpdateTokenPair(id uuid.UUID, userAgent string,
	ip string, refreshToken string, accessJTI string) (*dto.TokenPair, error) {
	refreshRecord, err := s.repo.GetById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customErrors.ErrRecordNotFound
		}
		return nil, customErrors.ErrDB
	}

	if refreshRecord.ClientIp != ip {
		err = s.notifier.Notify(refreshRecord.ClientIp, ip, userAgent, id.String())
		if err != nil {
			return nil, customErrors.ErrWebhook
		}
	}

	if refreshRecord.UserAgent != userAgent {
		err = s.repo.DeleteById(refreshRecord.UserID)
		if err != nil {
			return nil, customErrors.ErrUnexpectedState
		}
		return nil, customErrors.ErrUserAgentMismatch
	}

	if refreshRecord.UserID != id {
		return nil, customErrors.ErrUserIDMismatch
	}

	if refreshRecord.AccessJTI != accessJTI {
		return nil, customErrors.ErrTokenPairing
	}

	err = bcrypt.CompareHashAndPassword(refreshRecord.RefreshHash, []byte(refreshToken))

	if err != nil {
		return nil, customErrors.ErrRefreshMismatch
	}

	return s.GetNewTokenPair(id, userAgent, ip)
}

func (s *TokenService) ExistsById(id uuid.UUID) error {
	_, err := s.repo.ExistsById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customErrors.ErrRecordNotFound
		}
		return customErrors.ErrDB
	}
	return err
}

func (s *TokenService) DeleteById(id uuid.UUID) error {
	err := s.repo.DeleteById(id)
	if err != nil {
		return customErrors.ErrUnexpectedState
	}
	return err
}

func (s *TokenService) generateTokenPair(id uuid.UUID, accessJTI string) (*dto.TokenPair, error) {
	genTime := time.Now()
	tokenPair := dto.TokenPair{}
	var err error

	tokenPair.AccessToken, err = s.jwtUtils.GenerateJWT(id, genTime.Add(time.Hour).Unix(), genTime.Unix(), accessJTI)

	if err != nil {
		return nil, customErrors.ErrUnexpectedState
	}

	bytes := make([]byte, base64Length)
	_, err = rand.Read(bytes)

	if err != nil {
		return nil, customErrors.ErrUnexpectedState
	}

	tokenPair.RefreshToken = strings.TrimRight(base64.URLEncoding.EncodeToString(bytes), "=")

	return &tokenPair, nil
}
