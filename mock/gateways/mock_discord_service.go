package gateways

import (
	"github.com/stretchr/testify/mock"
)

type MockDiscordService struct {
	mock.Mock
}

func (m *MockDiscordService) SendSecurityCouncilMessage(msg string) error {
	args := m.Called(msg)
	return args.Error(0)
}
