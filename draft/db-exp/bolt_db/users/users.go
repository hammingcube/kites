package users

type User struct {
	Id       string
	Username string
	Email    string
}

type Store interface {
	Open(dbName, bucketName string) error
	Close()
	Get(key string) (*User, error)
	GetAll() ([]User, error)
	Post(u *User) error
	Delete(key string) error
}
