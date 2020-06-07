package eltype

type Stack struct {
	index int
	list  []string
}

func NewStack() (*Stack, error) {
	stack := new(Stack)
	stack.index = -1
	return stack, nil
}

func (stack *Stack) Push(item string) {
	stack.index++
	stack.list = append(stack.list, item)
}

func (stack *Stack) Top() string {
	if stack.index < 0 {
		return ""
	}
	return stack.list[stack.index]
}

func (stack *Stack) Pop() string {
	if stack.index < 0 {
		return ""
	}
	ret := stack.list[stack.index]
	stack.index--
	return ret
}
