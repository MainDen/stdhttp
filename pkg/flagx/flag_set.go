package flagx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrHelp = errors.New("flagx: help requested")

type errorHandling int

const (
	ContinueOnError errorHandling = iota
	ExitOnError
	PanicOnError
)

type CommandHandler interface {
	Parse(args ...string) (int, error)
	Usage(args ...string)
	ShortUsage() string
}

var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

var defaultOutput io.Writer

var defaultHandler Handler = Nop()

func DefaultOutput() io.Writer {
	if defaultOutput != nil {
		return defaultOutput
	}
	return os.Stderr
}

func SetDefaultOutput(output io.Writer) {
	defaultOutput = output
}

var storedCommands []string

var commandStore Parser = Func(func(args ...string) (int, error) {
	storedCommands = append(storedCommands, args...)
	return len(args), nil
})

func GetStoredCommands() []string {
	return storedCommands
}

func GetStoredCommand() string {
	return strings.Join(storedCommands, " ")
}

type FlagSet struct {
	name              string
	prefix            string
	description       string
	shortUsage        string
	errorHandling     errorHandling
	printErrorOnError bool
	printUsageOnError bool
	output            io.Writer
	usage             func(args ...string)
	defaultHandler    Handler
	commandStore      Parser

	cmdHandlers        map[string]CommandHandler
	cmdHandlersIndex   []string
	cmdHandlersSorted  bool
	optFlags           map[string]*OptFlag
	optFlagsIndex      []string
	optFlagsSorted     bool
	envFlags           map[string]*EnvFlag
	envFlagsIndex      []string
	envFlagsSorted     bool
	paramFormats       map[string]string
	paramFormatsIndex  []string
	paramFormatsSorted bool
}

func NewFlagSet(name string, errorHandling errorHandling) *FlagSet {
	return &FlagSet{
		name:              name,
		errorHandling:     errorHandling,
		printErrorOnError: true,
		printUsageOnError: true,
		defaultHandler:    defaultHandler,
		commandStore:      commandStore,
	}
}

func (fs *FlagSet) SetName(name string) {
	fs.name = name
}

func SetName(name string) {
	CommandLine.SetName(name)
}

func (fs *FlagSet) Name() string {
	return fs.name
}

func Name() string {
	return CommandLine.Name()
}

func (fs *FlagSet) SetPrefix(prefix string) {
	fs.prefix = prefix
}

func SetPrefix(prefix string) {
	CommandLine.SetPrefix(prefix)
}

func (fs *FlagSet) AddPrefix(prefix string) {
	if fs.prefix != "" {
		fs.prefix += "_"
	}
	fs.prefix += prefix
}

func AddPrefix(prefix string) {
	CommandLine.AddPrefix(prefix)
}

func (fs *FlagSet) Prefix() string {
	return fs.prefix
}

func Prefix() string {
	return CommandLine.Prefix()
}

func (fs *FlagSet) SetDescription(description string) {
	fs.description = description
}

func SetDescription(description string) {
	CommandLine.SetDescription(description)
}

func (fs *FlagSet) Description() string {
	return fs.description
}

func Description() string {
	return CommandLine.Description()
}

func (fs *FlagSet) SetShortUsage(shortUsage string) {
	fs.shortUsage = shortUsage
	if fs.description == "" {
		fs.description = shortUsage
	}
}

func SetShortUsage(shortUsage string) {
	CommandLine.SetShortUsage(shortUsage)
}

func (fs *FlagSet) ShortUsage() string {
	return fs.shortUsage
}

func ShortUsage() string {
	return CommandLine.ShortUsage()
}

func (fs *FlagSet) SetErrorHandling(errorHandling errorHandling) {
	fs.errorHandling = errorHandling
}

func SetErrorHandling(errorHandling errorHandling) {
	CommandLine.SetErrorHandling(errorHandling)
}

func (fs *FlagSet) ErrorHandling() errorHandling {
	return fs.errorHandling
}

func ErrorHandling() errorHandling {
	return CommandLine.ErrorHandling()
}

func (fs *FlagSet) SetPrintErrorOnError(print bool) {
	fs.printErrorOnError = print
}

func SetPrintErrorOnError(print bool) {
	CommandLine.SetPrintErrorOnError(print)
}

func (fs *FlagSet) PrintErrorOnError() bool {
	return fs.printErrorOnError
}

func PrintErrorOnError() bool {
	return CommandLine.PrintErrorOnError()
}

func (fs *FlagSet) SetPrintUsageOnError(print bool) {
	fs.printUsageOnError = print
}

func SetPrintUsageOnError(print bool) {
	CommandLine.SetPrintUsageOnError(print)
}

func (fs *FlagSet) PrintUsageOnError() bool {
	return fs.printUsageOnError
}

func PrintUsageOnError() bool {
	return CommandLine.PrintUsageOnError()
}

func (fs *FlagSet) SetOutput(output io.Writer) {
	fs.output = output
}

func SetOutput(output io.Writer) {
	CommandLine.SetOutput(output)
}

func (fs *FlagSet) Output() io.Writer {
	if fs.output != nil {
		return fs.output
	}
	return DefaultOutput()
}

func Output() io.Writer {
	return CommandLine.Output()
}

func (fs *FlagSet) SetUsage(usage func(args ...string)) {
	fs.usage = usage
}

func SetUsage(usage func(args ...string)) {
	CommandLine.SetUsage(usage)
}

func (fs *FlagSet) Usage(args ...string) {
	if fs.usage != nil {
		fs.usage(args...)
		return
	}
	if len(args) != 0 && fs.cmdHandlers[args[0]] != nil {
		fs.cmdHandlers[args[0]].Usage(args[1:]...)
		return
	}
	_, _ = io.WriteString(fs.Output(), formatUsage(fs))
}

func Usage(args ...string) {
	CommandLine.Usage(args...)
}

func (fs *FlagSet) SetDefaultHandler(handler Handler) {
	fs.defaultHandler = handler
}

func SetDefaultHandler(handler Handler) {
	CommandLine.SetDefaultHandler(handler)
}

func (fs *FlagSet) DefaultHandler() Handler {
	return fs.defaultHandler
}

func DefaultHandler() Handler {
	return CommandLine.DefaultHandler()
}

func (fs *FlagSet) SetCommandStore(parser Parser) {
	fs.commandStore = parser
}

func SetCommandStore(parser Parser) {
	CommandLine.SetCommandStore(parser)
}

func (fs *FlagSet) CommandStore() Parser {
	return fs.commandStore
}

func CommandStore() Parser {
	return CommandLine.CommandStore()
}

func (fs *FlagSet) SetSort(cmd bool, opt bool, env bool, param bool) {
	fs.cmdHandlersSorted = cmd
	fs.optFlagsSorted = opt
	fs.envFlagsSorted = env
	fs.paramFormatsSorted = param
}

func SetSort(cmd bool, opt bool, env bool, param bool) {
	CommandLine.SetSort(cmd, opt, env, param)
}

func (fs *FlagSet) RegisterCmd(name string, handler CommandHandler) {
	if _, ok := fs.cmdHandlers[name]; ok {
		panic(fmt.Sprintf("'%v': command already registered", name))
	}
	if fs.cmdHandlers == nil {
		fs.cmdHandlers = make(map[string]CommandHandler)
	}
	fs.cmdHandlers[name] = handler
	fs.cmdHandlersIndex = append(fs.cmdHandlersIndex, name)
	if fs.defaultHandler == defaultHandler {
		fs.defaultHandler = nil
	}
}

func RegisterCmd(name string, handler CommandHandler) {
	CommandLine.RegisterCmd(name, handler)
}

func (fs *FlagSet) RegisterOpt(flag *OptFlag) {
	if _, ok := fs.optFlags[flag.Name]; ok {
		panic(fmt.Sprintf("'%v': option already registered", flag.Name))
	}
	if flag.Alias != 0 {
		if _, ok := fs.optFlags[string(flag.Alias)]; ok {
			panic(fmt.Sprintf("'%v': option already registered", string(flag.Alias)))
		}
	}
	if fs.optFlags == nil {
		fs.optFlags = make(map[string]*OptFlag)
	}
	fs.optFlags[flag.Name] = flag
	fs.optFlagsIndex = append(fs.optFlagsIndex, flag.Name)
	if flag.Alias != 0 {
		fs.optFlags[string(flag.Alias)] = flag
	}
}

func RegisterOpt(flag *OptFlag) {
	CommandLine.RegisterOpt(flag)
}

func (fs *FlagSet) RegisterEnv(flag *EnvFlag) {
	if _, ok := fs.envFlags[flag.Name]; ok {
		panic(fmt.Sprintf("'%v': env already registered", flag.Name))
	}
	if fs.envFlags == nil {
		fs.envFlags = make(map[string]*EnvFlag)
	}
	fs.envFlags[flag.Name] = flag
	fs.envFlagsIndex = append(fs.envFlagsIndex, flag.Name)
}

func RegisterEnv(flag *EnvFlag) {
	CommandLine.RegisterEnv(flag)
}

func (fs *FlagSet) RegisterParam(name string, format string) {
	if _, ok := fs.paramFormats[name]; ok {
		panic(fmt.Sprintf("'%v': parameter already registered", name))
	}
	if fs.paramFormats == nil {
		fs.paramFormats = make(map[string]string)
	}
	fs.paramFormats[name] = format
	fs.paramFormatsIndex = append(fs.paramFormatsIndex, name)
}

func RegisterParam(name string, format string) {
	CommandLine.RegisterParam(name, format)
}

func (fs *FlagSet) AddOpt(name string, alias rune, params string, usage string, value Value, wrappers ...Wrapper) *OptFlag {
	value = Wrap(value, wrappers...)
	flag := &OptFlag{
		Name:     name,
		Alias:    alias,
		Params:   params,
		Usage:    usage,
		Value:    value,
		Defaults: FormatSlice(value),
	}
	fs.RegisterOpt(flag)
	return flag
}

func (fs *FlagSet) AddCmd(name string) *FlagSet {
	cmd := NewFlagSet(fs.name+" "+name, fs.errorHandling)
	cmd.SetPrefix(fs.prefix)
	cmd.SetOutput(fs.output)
	cmd.SetSort(fs.optFlagsSorted, fs.envFlagsSorted, fs.cmdHandlersSorted, fs.paramFormatsSorted)
	cmd.SetPrintErrorOnError(fs.printErrorOnError)
	cmd.SetPrintUsageOnError(fs.printUsageOnError)
	cmd.SetCommandStore(fs.commandStore)
	fs.RegisterCmd(name, cmd)
	return cmd
}

func AddCmd(name string) *FlagSet {
	return CommandLine.AddCmd(name)
}

func AddOpt(name string, alias rune, params string, usage string, value Value, wrappers ...Wrapper) *OptFlag {
	return CommandLine.AddOpt(name, alias, params, usage, value, wrappers...)
}

func (fs *FlagSet) AddEnv(name string, params string, usage string, value Value, wrappers ...Wrapper) *EnvFlag {
	value = Wrap(value, wrappers...)
	if fs.prefix != "" {
		name = fs.prefix + "_" + name
	}
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ToUpper(name)
	flag := &EnvFlag{
		Name:    name,
		Params:  params,
		Usage:   usage,
		Value:   value,
		Default: Format(value),
	}
	fs.RegisterEnv(flag)
	return flag
}

func AddEnv(name string, params string, usage string, value Value, wrappers ...Wrapper) *EnvFlag {
	return CommandLine.AddEnv(name, params, usage, value)
}

func (fs *FlagSet) AddOptEnv(name string, alias rune, params string, usage string, value Value, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	optFlag := fs.AddOpt(name, alias, params, usage, value, wrappers...)
	envFlag := fs.AddEnv(name, params, usage, value, wrappers...)
	return optFlag, envFlag
}

func AddOptEnv(name string, alias rune, params string, usage string, value Value, wrappers ...Wrapper) (*OptFlag, *EnvFlag) {
	return CommandLine.AddOptEnv(name, alias, params, usage, value, wrappers...)
}

func (fs *FlagSet) AddParam(name string, format string) {
	fs.RegisterParam(name, format)
}

func AddParam(name string, format string) {
	CommandLine.AddParam(name, format)
}

func (fs *FlagSet) setDefaults() error {
	for _, flag := range fs.envFlags {
		if flag.Default != "" {
			if _, err := flag.Value.Parse(flag.Default); err != nil {
				return fmt.Errorf("setting default: env '%v': %w", flag.Name, err)
			}
		}
	}
	for _, flag := range fs.optFlags {
		if flag.Defaults != nil {
			if _, err := flag.Value.Parse(flag.Defaults...); err != nil {
				return fmt.Errorf("setting default: option '%v': %w", flag.Name, err)
			}
		}
	}
	return nil
}

func (fs *FlagSet) parseOptName(name string, inlined bool, args ...string) (int, error) {
	if flag, ok := fs.optFlags[name]; ok && flag.Name == name {
		if IsInlined(flag.Value) && !inlined {
			args = nil
		}
		n, err := flag.Value.Parse(args...)
		if err != nil {
			return n, fmt.Errorf("option '%v': %w", name, err)
		}
		if inlined && n < len(args) {
			return n, fmt.Errorf("extra value %v", formatQuote(args[n]))
		}
		return n, nil
	}
	return 0, fmt.Errorf("unknown option '%v'", name)
}

func (fs *FlagSet) parseOptAliases(aliases string, inlined bool, args ...string) (int, error) {
	n, err := 0, error(nil)
	for i, alias := range aliases {
		args := args
		if i != len(aliases)-1 {
			args = nil
		}
		if flag, ok := fs.optFlags[string(alias)]; ok && flag.Alias == alias {
			if IsInlined(flag.Value) && !inlined {
				args = nil
			}
			n, err = flag.Value.Parse(args...)
			if err != nil {
				return n, fmt.Errorf("option '%v': %w", string(alias), err)
			}
			if inlined && n < len(args) {
				return n, fmt.Errorf("extra value: %v", formatQuote(args[n]))
			}
			continue
		}
		return i, fmt.Errorf("unknown option '%v'", string(alias))
	}
	return n, nil
}

func (fs *FlagSet) parseOpts(args ...string) (int, error) {
	if len(fs.optFlags) == 0 {
		return 0, nil
	}
	for i := 0; i < len(args); i++ {
		if args[i] == "--" {
			return i, nil
		}
		if args[i] == "-" {
			return i, errors.New("unexpected argument '-'")
		}
		arg, args := args[i], args[i+1:]
		arg, value, inlined := strings.Cut(arg, "=")
		if inlined {
			args = []string{value}
		}
		n, err := 0, error(nil)
		switch {
		case strings.HasPrefix(arg, "--"):
			n, err = fs.parseOptName(arg[2:], inlined, args...)
		case strings.HasPrefix(arg, "-"):
			n, err = fs.parseOptAliases(arg[1:], inlined, args...)
		default:
			return i, nil
		}
		if !inlined {
			i += n
		}
		if err != nil {
			return i, err
		}
	}
	return len(args), nil
}

func (fs *FlagSet) parseCmd(args ...string) (int, error) {
	if len(args) > 0 && args[0] == "--" && fs.defaultHandler != nil && fs.defaultHandler.Params() != "" && (len(fs.cmdHandlers) > 0 || len(fs.optFlags) > 0) {
		n, err := fs.defaultHandler.Parse(args[1:]...)
		return n + 1, err
	}
	if len(args) > 0 {
		if cmd, ok := fs.cmdHandlers[args[0]]; ok {
			n, err := cmd.Parse(args[1:]...)
			if err != nil {
				return n + 1, fmt.Errorf("command '%v': %w", args[0], err)
			}
			if fs.commandStore != nil {
				m, err := fs.commandStore.Parse(args[0])
				if err != nil {
					return n + 1, fmt.Errorf("store command '%v': %w", args[0], err)
				}
				if m != 1 {
					return n + 1, fmt.Errorf("store command '%v': unexpected result", args[0])
				}
			}
			return n + 1, nil
		}
	}
	if fs.defaultHandler != nil {
		n, err := fs.defaultHandler.Parse(args...)
		return n, err
	}
	return 0, fmt.Errorf("missing command")
}

func (fs *FlagSet) parse(args ...string) (int, error) {
	err := fs.setDefaults()
	if err != nil {
		return 0, err
	}
	n, err := fs.parseOpts(args...)
	if err != nil {
		return n, err
	}
	m, err := fs.parseCmd(args[n:]...)
	n += m
	if n < len(args) && err == nil {
		if len(args)-n == 1 {
			return n, fmt.Errorf("extra argument: %v", formatQuote(args[n]))
		}
		return n, fmt.Errorf("extra arguments: %v", strings.Join(formatQuoteList(args[n:]), " "))
	}
	if err != nil {
		return n, err
	}
	return n, nil
}

func (fs *FlagSet) Parse(args ...string) (int, error) {
	n, err := fs.parse(args...)
	if err != nil {
		if fs.printErrorOnError && !errors.Is(err, ErrHelp) {
			_, _ = io.WriteString(fs.Output(), formatError(err))
		}
		if fs.printUsageOnError && !errors.Is(err, ErrHelp) {
			fs.Usage()
		}
		switch fs.errorHandling {
		case ContinueOnError:
			return n, err
		case ExitOnError:
			if errors.Is(err, ErrHelp) {
				os.Exit(0)
			}
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return n, nil
}

func Parse(args ...string) (int, error) {
	return CommandLine.Parse(args...)
}
