package mention

type Storage interface {
	Get(key string) (string, error)
	List(key string) ([]string, error)
	Add(key string, value string) error
	Remove(key string, value string) error
	RemoveAll(key string) error
	Rotate(key string) error
}
