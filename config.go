package shync

type Config struct {
	Store, Username, Password string
	In, Out                   string
	AllTemplates              bool
	Templates                 []string
}
