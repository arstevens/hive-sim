package simulator

type Analzyer interface {
	Process(map[string]interface{}) error
	Summary() map[string]interface{}
	Save() error
}
