package testutil

type MockRepo[T any] struct {
	Model  *T
	Models []*T
}

func (r *MockRepo[T]) Create(t *T) (*T, error) {
	return r.Model, nil
}

func (r *MockRepo[T]) Update(t *T) (*T, error) {
	return r.Model, nil
}

func (r *MockRepo[T]) Delete(t *T) error {
	return nil
}

func (r *MockRepo[T]) FindById(id any) (*T, error) {
	return r.Model, nil
}

func (r *MockRepo[T]) FindAllByIds(id []any) ([]*T, error) {
	return r.Models, nil
}
