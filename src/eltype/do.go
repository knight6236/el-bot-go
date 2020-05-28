package eltype

type Do struct {
	IsCount       bool        `yaml:"isCount"`
	Message       Message     `yaml:"message"`
	OperationList []Operation `yaml:"operation"`
}
