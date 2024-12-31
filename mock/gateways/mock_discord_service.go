package gateways

import (
	"github.com/stretchr/testify/mock"
)

type MockDiscordService struct {
	mock.Mock
}

func (m *MockDiscordService) SendSecurityCouncilMessage(msg string, private bool	) error {
	args := m.Called(msg, private)
	return args.Error(0)
}

