package email

import (
	"context"
	authv1 "fbt/backend/gen/proto/go/auth/v1"
	"fbt/backend/gen/proto/go/auth/v1/authv1connect"
	"fbt/backend/internal/domain/auth/model"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/errors"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/util"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

type Server struct {
	service service.Service
	repo    Repo
}

func NewServiceHandler(service service.Service, repo Repo, opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewEmailServiceHandler(&Server{service, repo}, opts...)
}

func (s *Server) SendOTP(ctx context.Context, in *authv1.EmailServiceSendOTPRequest) (*authv1.EmailServiceSendOTPResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	otp, err := util.GenerateOTP(6)
	if err != nil {
		return nil, err
	}

	verificationId := util.GenerateBase32UUID()
	if err := s.service.SendVerificationMail(in.Email, otp); err != nil {
		return nil, err
	}

	err = s.repo.CreateEmailVerification(ctx, &model.EmailVerification{
		UserID:         auth.Session.UserId,
		VerificationID: verificationId,
		Otp:            otp,
		ExpiresAt:      pgtype.Timestamp{Time: time.Now().Add(time.Hour), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &authv1.EmailServiceSendOTPResponse{VerificationId: verificationId}, nil
}

func (s *Server) Verify(ctx context.Context, in *authv1.EmailServiceVerifyRequest) (*authv1.EmailServiceVerifyResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	auth, err := interceptor.FromAuthContext(ctx)
	if err != nil {
		return nil, err
	}

	emailVerification, err := s.repo.GetEmailVerification(ctx, auth.Session.UserId)
	if err != nil {
		return nil, err
	}

	if time.Now().After(emailVerification.ExpiresAt.Time) {
		if err := s.repo.DeleteEmailVerification(ctx, auth.Session.UserId); err != nil {
			return nil, errors.DBError
		}
		return nil, errors.SessionExpire
	}

	if emailVerification.VerificationID != in.VerificationId || emailVerification.Otp != in.Otp {
		return nil, errors.BadRequest
	}

	err = s.repo.VerifyEmail(ctx, auth.Session.UserId)
	if err != nil {
		return nil, err
	}

	session, err := s.service.CreateSession(ctx, auth.Session.UserId, false)
	if err != nil {
		return nil, err
	}
	return &authv1.EmailServiceVerifyResponse{Session: session.ToProto()}, nil
}
