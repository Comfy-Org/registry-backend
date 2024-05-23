package gateways

import (
	"github.com/stretchr/testify/mock"
)

type MockSlackService struct {
	mock.Mock
}

func (m *MockSlackService) SendRegistryMessageToSlack(msg string) error {
	args := m.Called(msg)
	return args.Error(0)
}
