package eltype

type When struct {
	Message       Message     `yaml:"message"`
	OperationList []Operation `yaml:"operation"`
}
