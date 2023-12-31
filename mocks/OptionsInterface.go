// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	model "quiz/model"

	mock "github.com/stretchr/testify/mock"
)

// OptionsInterface is an autogenerated mock type for the OptionsInterface type
type OptionsInterface struct {
	mock.Mock
}

// DeleteOption provides a mock function with given fields: questionsId, userId
func (_m *OptionsInterface) DeleteOption(questionsId uint, userId uint) int {
	ret := _m.Called(questionsId, userId)

	var r0 int
	if rf, ok := ret.Get(0).(func(uint, uint) int); ok {
		r0 = rf(questionsId, userId)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetAllOptionsFromQuiz provides a mock function with given fields: questionsId
func (_m *OptionsInterface) GetAllOptionsFromQuiz(questionsId uint) ([]model.Options, int) {
	ret := _m.Called(questionsId)

	var r0 []model.Options
	var r1 int
	if rf, ok := ret.Get(0).(func(uint) ([]model.Options, int)); ok {
		return rf(questionsId)
	}
	if rf, ok := ret.Get(0).(func(uint) []model.Options); ok {
		r0 = rf(questionsId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Options)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) int); ok {
		r1 = rf(questionsId)
	} else {
		r1 = ret.Get(1).(int)
	}

	return r0, r1
}

// GetOptionByID provides a mock function with given fields: id
func (_m *OptionsInterface) GetOptionByID(id uint) (*model.Options, int) {
	ret := _m.Called(id)

	var r0 *model.Options
	var r1 int
	if rf, ok := ret.Get(0).(func(uint) (*model.Options, int)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *model.Options); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Options)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) int); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Get(1).(int)
	}

	return r0, r1
}

// InsertOption provides a mock function with given fields: newOptions
func (_m *OptionsInterface) InsertOption(newOptions model.Options) (*model.Options, int) {
	ret := _m.Called(newOptions)

	var r0 *model.Options
	var r1 int
	if rf, ok := ret.Get(0).(func(model.Options) (*model.Options, int)); ok {
		return rf(newOptions)
	}
	if rf, ok := ret.Get(0).(func(model.Options) *model.Options); ok {
		r0 = rf(newOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Options)
		}
	}

	if rf, ok := ret.Get(1).(func(model.Options) int); ok {
		r1 = rf(newOptions)
	} else {
		r1 = ret.Get(1).(int)
	}

	return r0, r1
}

// UpdateOption provides a mock function with given fields: updateOptions, userId
func (_m *OptionsInterface) UpdateOption(updateOptions model.Options, userId uint) (*model.Options, int) {
	ret := _m.Called(updateOptions, userId)

	var r0 *model.Options
	var r1 int
	if rf, ok := ret.Get(0).(func(model.Options, uint) (*model.Options, int)); ok {
		return rf(updateOptions, userId)
	}
	if rf, ok := ret.Get(0).(func(model.Options, uint) *model.Options); ok {
		r0 = rf(updateOptions, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Options)
		}
	}

	if rf, ok := ret.Get(1).(func(model.Options, uint) int); ok {
		r1 = rf(updateOptions, userId)
	} else {
		r1 = ret.Get(1).(int)
	}

	return r0, r1
}

// NewOptionsInterface creates a new instance of OptionsInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOptionsInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *OptionsInterface {
	mock := &OptionsInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
