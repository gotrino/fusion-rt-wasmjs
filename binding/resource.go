package binding

func ResourceFrom[T any](repo Repository[T], id string) ResourceRepositoryAdapter[T] {
	return ResourceRepositoryAdapter[T]{repo, id}
}

type ResourceRepositoryAdapter[T any] struct {
	repo Repository[T]
	id   string
}

func (r ResourceRepositoryAdapter[T]) Load() (T, error) {
	return r.repo.Load(r.id)
}

func (r ResourceRepositoryAdapter[T]) Save(t T) error {
	return r.repo.Save(t)
}

func (r ResourceRepositoryAdapter[T]) Delete() error {
	return r.repo.Delete(r.id)
}
