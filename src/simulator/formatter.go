package simulator

type Formatter interface {
	Format(map[string]Log)
	Save() error
}
