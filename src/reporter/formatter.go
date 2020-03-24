package simulator

type Formatter interface {
	Format(map[string]interface{}) error
	Save() error
}
