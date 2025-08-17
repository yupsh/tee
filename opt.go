package command

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

type flags struct {
	Append          AppendFlag
	IgnoreInterrupt IgnoreInterruptFlag
}

func (f AppendFlag) Configure(flags *flags)          { flags.Append = f }
func (f IgnoreInterruptFlag) Configure(flags *flags) { flags.IgnoreInterrupt = f }
