package bytecode

import (
	"dap/internal/common"
	"fmt"
	"math"
)

// VM is the stack-based Bytecode Virtual Machine execution engine.
type VM struct {
	chunk   *Chunk
	ip      int
	stack   []common.Value
	globals map[string]common.Value
}

// NewVM creates a new Virtual Machine instance.
func NewVM() *VM {
	return &VM{
		stack:   make([]common.Value, 0, 256),
		globals: make(map[string]common.Value),
	}
}

// Push pushes a value onto the VM operand stack.
func (vm *VM) Push(val common.Value) {
	vm.stack = append(vm.stack, val)
}

// Pop pops and returns the top value from the VM operand stack.
func (vm *VM) Pop() common.Value {
	if len(vm.stack) == 0 {
		return common.Null{}
	}
	val := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return val
}

// Peek returns the value at distance from top of stack without popping.
func (vm *VM) Peek(distance int) common.Value {
	idx := len(vm.stack) - 1 - distance
	if idx < 0 || idx >= len(vm.stack) {
		return common.Null{}
	}
	return vm.stack[idx]
}

// Run executes a compiled bytecode Chunk on the Virtual Machine.
func (vm *VM) Run(c *Chunk) error {
	vm.chunk = c
	vm.ip = 0

	for vm.ip < len(c.Instructions) {
		inst := c.Instructions[vm.ip]
		vm.ip++

		switch inst.Op {
		case OP_NOP:
			// No operation

		case OP_CONST:
			val := c.Constants[inst.Arg]
			vm.Push(val)

		case OP_NIL:
			vm.Push(common.Null{})

		case OP_TRUE:
			vm.Push(common.Boolean{Value: true})

		case OP_FALSE:
			vm.Push(common.Boolean{Value: false})

		case OP_POP:
			vm.Pop()

		case OP_DUP:
			top := vm.Peek(0)
			vm.Push(top)

		case OP_LOAD:
			name := c.Names[inst.Arg]
			val, exists := vm.globals[name]
			if !exists {
				return fmt.Errorf("line %d: undefined variable '%s'", inst.Line, name)
			}
			vm.Push(val)

		case OP_STORE:
			name := c.Names[inst.Arg]
			val := vm.Pop()
			vm.globals[name] = val

		case OP_ADD:
			b := vm.Pop()
			a := vm.Pop()
			if numA, okA := a.(common.Number); okA {
				res, err := numA.Add(b)
				if err != nil {
					return fmt.Errorf("line %d: %s", inst.Line, err.Details)
				}
				vm.Push(res)
			} else if strA, okA := a.(common.String); okA {
				res, err := strA.Add(b)
				if err != nil {
					return fmt.Errorf("line %d: %s", inst.Line, err.Details)
				}
				vm.Push(res)
			} else {
				return fmt.Errorf("line %d: unsupported operand for '+'", inst.Line)
			}

		case OP_SUB:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '-' expects number operands", inst.Line)
			}
			res, err := numA.Subtract(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_MUL:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '*' expects number operands", inst.Line)
			}
			res, err := numA.Multiply(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_DIV:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '/' expects number operands", inst.Line)
			}
			res, err := numA.Divide(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_IDIV:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: 'div' expects integer operands", inst.Line)
			}
			res, err := numA.Div(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_MOD:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '%%' expects number operands", inst.Line)
			}
			res, err := numA.Mod(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_POW:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '**' expects number operands", inst.Line)
			}
			res, err := numA.Power(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_NEG:
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: unary minus expects number operand", inst.Line)
			}
			vm.Push(common.Number{Value: -numA.Value, IsReal: numA.IsReal})

		case OP_BIT_AND:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '&' expects integer operands", inst.Line)
			}
			res, err := numA.BitwiseAnd(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_BIT_OR:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '|' expects integer operands", inst.Line)
			}
			res, err := numA.BitwiseOr(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_BIT_XOR:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '^' expects integer operands", inst.Line)
			}
			res, err := numA.BitwiseXor(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_BIT_NOT:
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '~' expects integer operand", inst.Line)
			}
			res, err := numA.BitwiseNot()
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_SHL:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '<<' expects integer operands", inst.Line)
			}
			res, err := numA.ShiftLeft(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_SHR:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			if !okA {
				return fmt.Errorf("line %d: '>>' expects integer operands", inst.Line)
			}
			res, err := numA.ShiftRight(b)
			if err != nil {
				return fmt.Errorf("line %d: %s", inst.Line, err.Details)
			}
			vm.Push(res)

		case OP_EQ:
			b := vm.Pop()
			a := vm.Pop()
			if numA, okA := a.(common.Number); okA {
				if numB, okB := b.(common.Number); okB {
					vm.Push(common.Boolean{Value: numA.Value == numB.Value})
				} else {
					vm.Push(common.Boolean{Value: false})
				}
			} else if strA, okA := a.(common.String); okA {
				if strB, okB := b.(common.String); okB {
					vm.Push(common.Boolean{Value: strA.Value == strB.Value})
				} else {
					vm.Push(common.Boolean{Value: false})
				}
			} else if boolA, okA := a.(common.Boolean); okA {
				if boolB, okB := b.(common.Boolean); okB {
					vm.Push(common.Boolean{Value: boolA.Value == boolB.Value})
				} else {
					vm.Push(common.Boolean{Value: false})
				}
			} else {
				vm.Push(common.Boolean{Value: false})
			}

		case OP_NEQ:
			b := vm.Pop()
			a := vm.Pop()
			if numA, okA := a.(common.Number); okA {
				if numB, okB := b.(common.Number); okB {
					vm.Push(common.Boolean{Value: numA.Value != numB.Value})
				} else {
					vm.Push(common.Boolean{Value: true})
				}
			} else if strA, okA := a.(common.String); okA {
				if strB, okB := b.(common.String); okB {
					vm.Push(common.Boolean{Value: strA.Value != strB.Value})
				} else {
					vm.Push(common.Boolean{Value: true})
				}
			} else {
				vm.Push(common.Boolean{Value: true})
			}

		case OP_LT:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			numB, okB := b.(common.Number)
			if !okA || !okB {
				return fmt.Errorf("line %d: '<' expects number operands", inst.Line)
			}
			vm.Push(common.Boolean{Value: numA.Value < numB.Value})

		case OP_LTE:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			numB, okB := b.(common.Number)
			if !okA || !okB {
				return fmt.Errorf("line %d: '<=' expects number operands", inst.Line)
			}
			vm.Push(common.Boolean{Value: numA.Value <= numB.Value})

		case OP_GT:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			numB, okB := b.(common.Number)
			if !okA || !okB {
				return fmt.Errorf("line %d: '>' expects number operands", inst.Line)
			}
			vm.Push(common.Boolean{Value: numA.Value > numB.Value})

		case OP_GTE:
			b := vm.Pop()
			a := vm.Pop()
			numA, okA := a.(common.Number)
			numB, okB := b.(common.Number)
			if !okA || !okB {
				return fmt.Errorf("line %d: '>=' expects number operands", inst.Line)
			}
			vm.Push(common.Boolean{Value: numA.Value >= numB.Value})

		case OP_NOT:
			a := vm.Pop()
			vm.Push(common.Boolean{Value: !a.IsTrue()})

		case OP_JUMP:
			vm.ip = inst.Arg

		case OP_JUMP_IF_FALSE:
			cond := vm.Pop()
			if !cond.IsTrue() {
				vm.ip = inst.Arg
			}

		case OP_JUMP_IF_TRUE:
			cond := vm.Pop()
			if cond.IsTrue() {
				vm.ip = inst.Arg
			}

		case OP_LOOP:
			vm.ip = inst.Arg

		case OP_PRINT:
			val := vm.Pop()
			if num, ok := val.(common.Number); ok {
				if num.IsReal {
					if num.Value == math.Trunc(num.Value) {
						fmt.Printf("%.1f\n", num.Value)
					} else {
						fmt.Printf("%v\n", num.Value)
					}
				} else {
					fmt.Printf("%d\n", int64(num.Value))
				}
			} else {
				fmt.Println(common.PrintValueInterpreter(val))
			}

		case OP_HALT:
			return nil

		default:
			return fmt.Errorf("unknown opcode %s", inst.Op.String())
		}
	}

	return nil
}
