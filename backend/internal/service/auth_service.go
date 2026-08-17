package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/techblog/community/internal/cache"
	"github.com/techblog/community/internal/config"
	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

const (
	verifCodePrefix = "verif:code:"
	verifCodeTTL    = 5 * time.Minute
)

// AuthService 认证服务：注册 / 登录 / 验证码 / GitHub 登录
type AuthService struct {
	userRepo  *repository.UserRepository
	loginRepo *repository.LoginRepository
	userSvc   *UserService
	auditSvc  *AuditService
	redis     *cache.Redis
	cfg       *config.Config
}

// NewAuthService 构造认证服务
func NewAuthService(userRepo *repository.UserRepository, loginRepo *repository.LoginRepository, userSvc *UserService, auditSvc *AuditService, redis *cache.Redis, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, loginRepo: loginRepo, userSvc: userSvc, auditSvc: auditSvc, redis: redis, cfg: cfg}
}

// SendCode 发送邮箱验证码（演示模式返回 debug_code）
func (s *AuthService) SendCode(ctx context.Context, email string) (*dto.SendCodeResponse, error) {
	key := verifCodePrefix + email
	code, err := s.redis.GenerateCode(ctx, key, verifCodeTTL)
	if err != nil {
		util.LogError(util.GetRequestID(ctx), "send email code failed", "email", email, "error", err.Error())
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("验证码发送失败: email=%s", email), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserCodeSent, email, int(verifCodeTTL.Seconds())))
	resp := &dto.SendCodeResponse{Email: email, ExpiresIn: int(verifCodeTTL.Seconds())}
	if s.cfg.DebugCodeMode {
		resp.DebugCode = code
	}
	return resp, nil
}

// Register 邮箱注册（验证码 + 密码）
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest, ip string) (*dto.AuthResponse, error) {
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserRegisterAttempt, req.Email))
	if err := s.verifyCode(ctx, req.Email, req.Code); err != nil {
		return nil, err
	}
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, util.NewAppError(constants.CodeEmailExists, fmt.Sprintf(constants.MsgErrEmailExists, req.Email))
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询用户失败: email=%s", req.Email), err)
	}
	if err := s.userSvc.CheckNicknameUnique(ctx, req.Nickname, 0); err != nil {
		return nil, err
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "密码加密失败: field=password", err)
	}
	tagsJSON, _ := json.Marshal(req.TechTags)
	user := &model.User{
		Email:        req.Email,
		PasswordHash: hash,
		Nickname:     req.Nickname,
		Avatar:       req.Avatar,
		Bio:          req.Bio,
		TechTags:     string(tagsJSON),
		Role:         constants.RoleUser.Int(),
		Status:       constants.UserStatusActive.Int(),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicateEntry) {
			return nil, util.NewAppError(constants.CodeEmailExists, fmt.Sprintf(constants.MsgErrEmailExists, req.Email))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("创建用户失败: email=%s", req.Email), err)
	}
	_ = s.redis.Del(ctx, verifCodePrefix+req.Email)
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserRegistered, user.ID, user.Email, user.Nickname))
	s.auditSvc.Record(ctx, user.ID, user.Nickname, user.Role, "REGISTER", "auth", "用户注册", ip)
	s.recordLogin(ctx, user, ip)
	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}
	profile, err := s.userSvc.BuildProfile(ctx, user, user.ID)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: token, User: *profile}, nil
}

// Login 邮箱密码登录
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest, ip string) (*dto.AuthResponse, error) {
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserLoginAttempt, req.Email))
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserLoginFail, req.Email, "user not found"))
			return nil, util.NewAppError(constants.CodeInvalidCredentials, fmt.Sprintf("邮箱或密码错误: email=%s role=%d", req.Email, constants.RoleUser.Int()))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询用户失败: email=%s", req.Email), err)
	}
	if user.Status != constants.UserStatusActive.Int() {
		return nil, util.NewAppError(constants.CodeUserDisabled, fmt.Sprintf(constants.MsgErrUserDisabled, user.ID, user.Role))
	}
	if !util.VerifyPassword(user.PasswordHash, req.Password) {
		util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserLoginFail, req.Email, "wrong password"))
		return nil, util.NewAppError(constants.CodeInvalidCredentials, fmt.Sprintf("邮箱或密码错误: email=%s role=%d", req.Email, user.Role))
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserLoginSuccess, user.ID, user.Email))
	s.auditSvc.Record(ctx, user.ID, user.Nickname, user.Role, "LOGIN", "auth", "密码登录", ip)
	s.recordLogin(ctx, user, ip)
	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}
	profile, err := s.userSvc.BuildProfile(ctx, user, user.ID)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: token, User: *profile}, nil
}

// LoginWithCode 邮箱验证码登录
func (s *AuthService) LoginWithCode(ctx context.Context, req *dto.LoginCodeRequest, ip string) (*dto.AuthResponse, error) {
	if err := s.verifyCode(ctx, req.Email, req.Code); err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// 验证码登录自动注册
			hash, _ := util.HashPassword(fmt.Sprintf("code-%s-%d", req.Email, time.Now().UnixNano()))
			user = &model.User{
				Email:        req.Email,
				PasswordHash: hash,
				Nickname:     "用户" + req.Email[:6],
				Role:         constants.RoleUser.Int(),
				Status:       constants.UserStatusActive.Int(),
				TechTags:     "[]",
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("创建用户失败: email=%s", req.Email), err)
			}
		} else {
			return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询用户失败: email=%s", req.Email), err)
		}
	}
	if user.Status != constants.UserStatusActive.Int() {
		return nil, util.NewAppError(constants.CodeUserDisabled, fmt.Sprintf(constants.MsgErrUserDisabled, user.ID, user.Role))
	}
	_ = s.redis.Del(ctx, verifCodePrefix+req.Email)
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserCodeLoginSuccess, user.ID, user.Email))
	s.recordLogin(ctx, user, ip)
	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}
	profile, err := s.userSvc.BuildProfile(ctx, user, user.ID)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: token, User: *profile}, nil
}

// GithubLogin GitHub 第三方登录（演示模式：使用 code 映射 github_id）
func (s *AuthService) GithubLogin(ctx context.Context, req *dto.GithubLoginRequest, ip string) (*dto.AuthResponse, error) {
	githubID := "gh_" + req.Code
	user, err := s.userRepo.FindByGithubID(ctx, githubID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			email := req.Email
			if email == "" {
				email = fmt.Sprintf("gh_%s@github.local", req.Code)
			}
			nickname := req.Nickname
			if nickname == "" {
				nickname = "GitHub用户"
			}
			hash, _ := util.HashPassword("github-oauth-" + githubID)
			user = &model.User{
				Email:        email,
				PasswordHash: hash,
				Nickname:     nickname,
				Role:         constants.RoleUser.Int(),
				Status:       constants.UserStatusActive.Int(),
				GithubID:     githubID,
				TechTags:     "[]",
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				if errors.Is(err, repository.ErrDuplicateEntry) {
					return nil, util.NewAppError(constants.CodeEmailExists, fmt.Sprintf(constants.MsgErrEmailExists, email))
				}
				return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("创建 GitHub 用户失败: github_id=%s", githubID), err)
			}
		} else {
			return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询 GitHub 用户失败: github_id=%s", githubID), err)
		}
	}
	if user.Status != constants.UserStatusActive.Int() {
		return nil, util.NewAppError(constants.CodeUserDisabled, fmt.Sprintf(constants.MsgErrUserDisabled, user.ID, user.Role))
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogGithubLogin, user.ID, user.GithubID))
	s.recordLogin(ctx, user, ip)
	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}
	profile, err := s.userSvc.BuildProfile(ctx, user, user.ID)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: token, User: *profile}, nil
}


// Me 当前用户信息（复用 UserService.GetProfile）
func (s *AuthService) Me(ctx context.Context, userID uint) (*dto.UserProfileDTO, error) {
	return s.userSvc.GetProfile(ctx, userID, userID)
}

// verifyCode 校验验证码（内部方法，被 Register / LoginWithCode 复用）
func (s *AuthService) verifyCode(ctx context.Context, email, code string) error {
	key := verifCodePrefix + email
	got, err := s.redis.Get(ctx, key)
	if err != nil {
		if err.Error() == "redis: nil" {
			return util.NewAppError(constants.CodeCodeExpired, fmt.Sprintf(constants.MsgErrCodeExpired, email))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("读取验证码失败: email=%s", email), err)
	}
	if got != code {
		return util.NewAppError(constants.CodeInvalidCode, fmt.Sprintf(constants.MsgErrInvalidCode, email))
	}
	return nil
}

// issueToken 签发 JWT（内部方法）
func (s *AuthService) issueToken(user *model.User) (string, error) {
	token, err := util.GenerateToken(s.cfg.JWTSecret, user.ID, user.Role, user.Email, s.cfg.JWTTTL)
	if err != nil {
		return "", util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("签发令牌失败: user_id=%d role=%d", user.ID, user.Role), err)
	}
	return token, nil
}

// recordLogin 记录登录日志并更新最后登录时间（内部方法）
func (s *AuthService) recordLogin(ctx context.Context, user *model.User, ip string) {
	now := time.Now()
	_ = s.loginRepo.Create(ctx, &model.LoginLog{UserID: user.ID, LoginDate: now.Format("2006-01-02"), IP: ip})
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID, now)
}
