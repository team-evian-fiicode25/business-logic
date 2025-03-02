package services

import (
	"context"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

type ConsumerService struct {
	client graphql.Client
}

func NewConsumerService(endpoint string) *ConsumerService {
	return &ConsumerService{
		client: graphql.NewClient(endpoint, http.DefaultClient),
	}
}

func (s *ConsumerService) CreateConsumer(ctx context.Context, username, email, password string) (*CreateLoginResponse, error) {
	resp, err := CreateLogin(ctx, s.client, username, email, "", password)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ConsumerService) LogInWithPassword(ctx context.Context, username, password string) (*LogInWithPasswordResponse, error) {
	resp, err := LogInWithPassword(ctx, s.client, username, password)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
