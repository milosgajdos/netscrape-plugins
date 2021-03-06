package dgraph

// Op is dgraph operation.
type Op int

const (
	// AddOp is add operation
	AddOp Op = iota
	// DelOp is delete operation
	DelOp
	// GetOp is get operation
	GetOp
	// LinkOp is link operation
	LinkOp
	// UnlinkOp is unlink operation
	UnlinkOp
	// UnknownOp is unknown operation
	UnknownOp
)

// String implements fmt.Stringer.
func (op Op) String() string {
	switch op {
	case AddOp:
		return "AddOp"
	case DelOp:
		return "DelOp"
	case GetOp:
		return "GetOp"
	case LinkOp:
		return "LinkOp"
	case UnlinkOp:
		return "UnlinkOp"
	default:
		return "UnknownOp"
	}
}
