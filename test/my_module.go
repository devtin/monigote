package mymodule

//go:generate go run ../generator/generator.go -a my_module.go -b MyModule

type MyModule interface {
	Get(name string) (string, error)
	Set(name string, value string) error
}
