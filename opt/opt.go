package opt

// Boolean flag types with constants
type AppendFlag bool
const (
	Append    AppendFlag = true
	Overwrite AppendFlag = false
)

type IgnoreInterruptFlag bool
const (
	IgnoreInterrupt   IgnoreInterruptFlag = true
	NoIgnoreInterrupt IgnoreInterruptFlag = false
)

// Flags represents the configuration options for the tee command
type Flags struct {
	Append          AppendFlag          // Append to files instead of overwriting
	IgnoreInterrupt IgnoreInterruptFlag // Ignore interrupt signals
}

// Configure methods for the opt system
func (f AppendFlag) Configure(flags *Flags) { flags.Append = f }
func (f IgnoreInterruptFlag) Configure(flags *Flags) { flags.IgnoreInterrupt = f }
