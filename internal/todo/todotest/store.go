// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package todotest

import (
	"context"
	"github.com/dokyan1989/goit/internal/todo"
	"sync"
)

// Ensure, that Mock does implement TodoStore.
// If this is not the case, regenerate this file with moq.
var _ TodoStore = &Mock{}

// Mock is a mock implementation of TodoStore.
//
//	func TestSomethingThatUsesTodoStore(t *testing.T) {
//
//		// make and configure a mocked TodoStore
//		mockedTodoStore := &Mock{
//			CreateTodoFunc: func(ctx context.Context, params todo.CreateTodoParams) (int64, error) {
//				panic("mock out the CreateTodo method")
//			},
//			DeleteTodoFunc: func(ctx context.Context, id int64) error {
//				panic("mock out the DeleteTodo method")
//			},
//			FindTodosFunc: func(ctx context.Context, params todo.FindTodosParams) ([]todo.Todo, error) {
//				panic("mock out the FindTodos method")
//			},
//			UpdateTodoFunc: func(ctx context.Context, id int64, params todo.UpdateTodoParams) error {
//				panic("mock out the UpdateTodo method")
//			},
//		}
//
//		// use mockedTodoStore in code that requires TodoStore
//		// and then make assertions.
//
//	}
type Mock struct {
	// CreateTodoFunc mocks the CreateTodo method.
	CreateTodoFunc func(ctx context.Context, params todo.CreateTodoParams) (int64, error)

	// DeleteTodoFunc mocks the DeleteTodo method.
	DeleteTodoFunc func(ctx context.Context, id int64) error

	// FindTodosFunc mocks the FindTodos method.
	FindTodosFunc func(ctx context.Context, params todo.FindTodosParams) ([]todo.Todo, error)

	// UpdateTodoFunc mocks the UpdateTodo method.
	UpdateTodoFunc func(ctx context.Context, id int64, params todo.UpdateTodoParams) error

	// calls tracks calls to the methods.
	calls struct {
		// CreateTodo holds details about calls to the CreateTodo method.
		CreateTodo []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Params is the params argument value.
			Params todo.CreateTodoParams
		}
		// DeleteTodo holds details about calls to the DeleteTodo method.
		DeleteTodo []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int64
		}
		// FindTodos holds details about calls to the FindTodos method.
		FindTodos []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Params is the params argument value.
			Params todo.FindTodosParams
		}
		// UpdateTodo holds details about calls to the UpdateTodo method.
		UpdateTodo []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int64
			// Params is the params argument value.
			Params todo.UpdateTodoParams
		}
	}
	lockCreateTodo sync.RWMutex
	lockDeleteTodo sync.RWMutex
	lockFindTodos  sync.RWMutex
	lockUpdateTodo sync.RWMutex
}

// CreateTodo calls CreateTodoFunc.
func (mock *Mock) CreateTodo(ctx context.Context, params todo.CreateTodoParams) (int64, error) {
	if mock.CreateTodoFunc == nil {
		panic("Mock.CreateTodoFunc: method is nil but TodoStore.CreateTodo was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Params todo.CreateTodoParams
	}{
		Ctx:    ctx,
		Params: params,
	}
	mock.lockCreateTodo.Lock()
	mock.calls.CreateTodo = append(mock.calls.CreateTodo, callInfo)
	mock.lockCreateTodo.Unlock()
	return mock.CreateTodoFunc(ctx, params)
}

// CreateTodoCalls gets all the calls that were made to CreateTodo.
// Check the length with:
//
//	len(mockedTodoStore.CreateTodoCalls())
func (mock *Mock) CreateTodoCalls() []struct {
	Ctx    context.Context
	Params todo.CreateTodoParams
} {
	var calls []struct {
		Ctx    context.Context
		Params todo.CreateTodoParams
	}
	mock.lockCreateTodo.RLock()
	calls = mock.calls.CreateTodo
	mock.lockCreateTodo.RUnlock()
	return calls
}

// ResetCreateTodoCalls reset all the calls that were made to CreateTodo.
func (mock *Mock) ResetCreateTodoCalls() {
	mock.lockCreateTodo.Lock()
	mock.calls.CreateTodo = nil
	mock.lockCreateTodo.Unlock()
}

// DeleteTodo calls DeleteTodoFunc.
func (mock *Mock) DeleteTodo(ctx context.Context, id int64) error {
	if mock.DeleteTodoFunc == nil {
		panic("Mock.DeleteTodoFunc: method is nil but TodoStore.DeleteTodo was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  int64
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockDeleteTodo.Lock()
	mock.calls.DeleteTodo = append(mock.calls.DeleteTodo, callInfo)
	mock.lockDeleteTodo.Unlock()
	return mock.DeleteTodoFunc(ctx, id)
}

// DeleteTodoCalls gets all the calls that were made to DeleteTodo.
// Check the length with:
//
//	len(mockedTodoStore.DeleteTodoCalls())
func (mock *Mock) DeleteTodoCalls() []struct {
	Ctx context.Context
	ID  int64
} {
	var calls []struct {
		Ctx context.Context
		ID  int64
	}
	mock.lockDeleteTodo.RLock()
	calls = mock.calls.DeleteTodo
	mock.lockDeleteTodo.RUnlock()
	return calls
}

// ResetDeleteTodoCalls reset all the calls that were made to DeleteTodo.
func (mock *Mock) ResetDeleteTodoCalls() {
	mock.lockDeleteTodo.Lock()
	mock.calls.DeleteTodo = nil
	mock.lockDeleteTodo.Unlock()
}

// FindTodos calls FindTodosFunc.
func (mock *Mock) FindTodos(ctx context.Context, params todo.FindTodosParams) ([]todo.Todo, error) {
	if mock.FindTodosFunc == nil {
		panic("Mock.FindTodosFunc: method is nil but TodoStore.FindTodos was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Params todo.FindTodosParams
	}{
		Ctx:    ctx,
		Params: params,
	}
	mock.lockFindTodos.Lock()
	mock.calls.FindTodos = append(mock.calls.FindTodos, callInfo)
	mock.lockFindTodos.Unlock()
	return mock.FindTodosFunc(ctx, params)
}

// FindTodosCalls gets all the calls that were made to FindTodos.
// Check the length with:
//
//	len(mockedTodoStore.FindTodosCalls())
func (mock *Mock) FindTodosCalls() []struct {
	Ctx    context.Context
	Params todo.FindTodosParams
} {
	var calls []struct {
		Ctx    context.Context
		Params todo.FindTodosParams
	}
	mock.lockFindTodos.RLock()
	calls = mock.calls.FindTodos
	mock.lockFindTodos.RUnlock()
	return calls
}

// ResetFindTodosCalls reset all the calls that were made to FindTodos.
func (mock *Mock) ResetFindTodosCalls() {
	mock.lockFindTodos.Lock()
	mock.calls.FindTodos = nil
	mock.lockFindTodos.Unlock()
}

// UpdateTodo calls UpdateTodoFunc.
func (mock *Mock) UpdateTodo(ctx context.Context, id int64, params todo.UpdateTodoParams) error {
	if mock.UpdateTodoFunc == nil {
		panic("Mock.UpdateTodoFunc: method is nil but TodoStore.UpdateTodo was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		ID     int64
		Params todo.UpdateTodoParams
	}{
		Ctx:    ctx,
		ID:     id,
		Params: params,
	}
	mock.lockUpdateTodo.Lock()
	mock.calls.UpdateTodo = append(mock.calls.UpdateTodo, callInfo)
	mock.lockUpdateTodo.Unlock()
	return mock.UpdateTodoFunc(ctx, id, params)
}

// UpdateTodoCalls gets all the calls that were made to UpdateTodo.
// Check the length with:
//
//	len(mockedTodoStore.UpdateTodoCalls())
func (mock *Mock) UpdateTodoCalls() []struct {
	Ctx    context.Context
	ID     int64
	Params todo.UpdateTodoParams
} {
	var calls []struct {
		Ctx    context.Context
		ID     int64
		Params todo.UpdateTodoParams
	}
	mock.lockUpdateTodo.RLock()
	calls = mock.calls.UpdateTodo
	mock.lockUpdateTodo.RUnlock()
	return calls
}

// ResetUpdateTodoCalls reset all the calls that were made to UpdateTodo.
func (mock *Mock) ResetUpdateTodoCalls() {
	mock.lockUpdateTodo.Lock()
	mock.calls.UpdateTodo = nil
	mock.lockUpdateTodo.Unlock()
}

// ResetCalls reset all the calls that were made to all mocked methods.
func (mock *Mock) ResetCalls() {
	mock.lockCreateTodo.Lock()
	mock.calls.CreateTodo = nil
	mock.lockCreateTodo.Unlock()

	mock.lockDeleteTodo.Lock()
	mock.calls.DeleteTodo = nil
	mock.lockDeleteTodo.Unlock()

	mock.lockFindTodos.Lock()
	mock.calls.FindTodos = nil
	mock.lockFindTodos.Unlock()

	mock.lockUpdateTodo.Lock()
	mock.calls.UpdateTodo = nil
	mock.lockUpdateTodo.Unlock()
}