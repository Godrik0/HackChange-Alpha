package mocks

import "github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"

type MockLogger struct {
	DebugFunc func(msg string, args ...interface{})
	InfoFunc  func(msg string, args ...interface{})
	WarnFunc  func(msg string, args ...interface{})
	ErrorFunc func(msg string, args ...interface{})
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	if m.DebugFunc != nil {
		m.DebugFunc(msg, args...)
	}
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg, args...)
	}
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	if m.WarnFunc != nil {
		m.WarnFunc(msg, args...)
	}
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	if m.ErrorFunc != nil {
		m.ErrorFunc(msg, args...)
	}
}

func (m *MockLogger) With(args ...interface{}) interfaces.Logger {
	return m
}
