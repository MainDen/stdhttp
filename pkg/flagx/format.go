package flagx

import (
	"fmt"
	"sort"
	"strings"
)

const (
	usageFormat        = "Usage:\n"
	progParamsFormat   = "      %v\n"
	descriptionFormat  = "%v\n"
	cmdsFormat         = "Commands:\n"
	cmdFormat          = "      %v"
	optsFormat         = "Options:\n"
	optFormat          = "      --%v"
	optWithAliasFormat = "  -%v  --%v"
	envsFormat         = "Env:\n"
	envFormat          = "  %v"
	paramsFormat       = "Parameters:\n"
	paramFormat        = "  %v"
	shortUsagePadding  = 40
)

func getIndex(index []string, sorted bool) []string {
	if !sorted {
		return index
	}
	sortedIndex := make([]string, len(index))
	copy(sortedIndex, index)
	sort.Strings(sortedIndex)
	return sortedIndex
}

func formatQuote(text string) string {
	if text == "" || strings.ContainsAny(text, " \t,.") {
		return fmt.Sprintf("\"%v\"", text)
	}
	return text
}

func formatQuoteList(list []string) []string {
	quotedList := make([]string, len(list))
	for i, text := range list {
		quotedList[i] = formatQuote(text)
	}
	return quotedList
}

func formatNotEmpty(prefix string, text string) string {
	if text == "" {
		return ""
	}
	return prefix + text
}

func formatError(err error) string {
	return fmt.Sprintf("Syntax error: %v.\n", err)
}

func formatShortUsage(params string, usage string) string {
	if usage == "" {
		return params + "\n"
	}
	newLine := "\n" + strings.Repeat(" ", shortUsagePadding)
	if len(params) < shortUsagePadding {
		return params + strings.Repeat(" ", shortUsagePadding-len(params)) + strings.ReplaceAll(usage, "\n", newLine) + "\n"
	}
	return params + newLine + strings.ReplaceAll(usage, "\n", newLine) + "\n"
}

func writeProgParams(fs *FlagSet, builder *strings.Builder) {
	params := fs.name
	if len(fs.optFlags) > 0 {
		params += " [Option ...]"
	}
	if fs.defaultHandler != nil && fs.defaultHandler.Params() == "" {
		_, _ = builder.WriteString(fmt.Sprintf(progParamsFormat, params))
	}
	if len(fs.cmdHandlers) > 0 {
		_, _ = builder.WriteString(fmt.Sprintf(progParamsFormat, params+" Command [Arg ...]"))
	}
	if fs.defaultHandler != nil && fs.defaultHandler.Params() != "" {
		if len(fs.optFlags) > 0 || len(fs.cmdHandlers) > 0 {
			_, _ = builder.WriteString(fmt.Sprintf(progParamsFormat, params+" [--] "+fs.defaultHandler.Params()))
		} else {
			_, _ = builder.WriteString(fmt.Sprintf(progParamsFormat, params+" "+fs.defaultHandler.Params()))
		}
	}
	_, _ = builder.WriteString("\n")
}

func writeDescription(fs *FlagSet, builder *strings.Builder) {
	if fs.description != "" {
		_, _ = builder.WriteString(fmt.Sprintf(descriptionFormat, fs.description))
		_, _ = builder.WriteString("\n")
	}
}

func formatCmd(name string, cmd CommandHandler) string {
	return formatShortUsage(fmt.Sprintf(cmdFormat, name), cmd.ShortUsage())
}

func writeCmds(fs *FlagSet, builder *strings.Builder) {
	if len(fs.cmdHandlers) > 0 {
		_, _ = builder.WriteString(cmdsFormat)
		for _, name := range getIndex(fs.cmdHandlersIndex, fs.cmdHandlersSorted) {
			_, _ = builder.WriteString(formatCmd(name, fs.cmdHandlers[name]))
		}
		_, _ = builder.WriteString("\n")
	}
}

func formatOptParams(f *OptFlag) string {
	prefix := " "
	if IsInlined(f.Value) {
		prefix = "="
	}
	if f.Alias != 0 {
		return fmt.Sprintf(optWithAliasFormat, string(f.Alias), f.Name) + formatNotEmpty(prefix, f.Params)
	}
	return fmt.Sprintf(optFormat, f.Name) + formatNotEmpty(prefix, f.Params)
}

func formatOptUsage(f *OptFlag) string {
	if len(f.Defaults) > 1 {
		return fmt.Sprintf("%v  (defaults: %v)", f.Usage, strings.Join(formatQuoteList(f.Defaults), " "))
	}
	if len(f.Defaults) == 1 {
		return fmt.Sprintf("%v  (default: %v)", f.Usage, formatQuote(f.Defaults[0]))
	}
	return f.Usage
}

func formatOpt(f *OptFlag) string {
	return formatShortUsage(formatOptParams(f), formatOptUsage(f))
}

func writeOpts(fs *FlagSet, builder *strings.Builder) {
	if len(fs.optFlags) > 0 {
		_, _ = builder.WriteString(optsFormat)
		for _, name := range getIndex(fs.optFlagsIndex, fs.optFlagsSorted) {
			_, _ = builder.WriteString(formatOpt(fs.optFlags[name]))
		}
		_, _ = builder.WriteString("\n")
	}
}

func formatEnvParams(f *EnvFlag) string {
	return fmt.Sprintf(envFormat, f.Name+formatNotEmpty("=", f.Params))
}

func formatEnvUsage(f *EnvFlag) string {
	if f.Default != "" {
		return fmt.Sprintf("%v  (default: %v)", f.Usage, formatQuote(f.Default))
	}
	return f.Usage
}

func formatEnv(f *EnvFlag) string {
	return formatShortUsage(formatEnvParams(f), formatEnvUsage(f))
}

func writeEnvs(fs *FlagSet, builder *strings.Builder) {
	if len(fs.envFlags) > 0 {
		_, _ = builder.WriteString(envsFormat)
		for _, name := range getIndex(fs.envFlagsIndex, fs.envFlagsSorted) {
			_, _ = builder.WriteString(formatEnv(fs.envFlags[name]))
		}
		_, _ = builder.WriteString("\n")
	}
}

func formatParam(name string) string {
	return fmt.Sprintf(paramFormat, name)
}

func writeParams(fs *FlagSet, builder *strings.Builder) {
	if len(fs.paramFormats) > 0 {
		_, _ = builder.WriteString(paramsFormat)
		for _, name := range getIndex(fs.paramFormatsIndex, fs.paramFormatsSorted) {
			_, _ = builder.WriteString(formatShortUsage(formatParam(name), fs.paramFormats[name]))
		}
		_, _ = builder.WriteString("\n")
	}
}

func formatUsage(fs *FlagSet) string {
	builder := new(strings.Builder)
	_, _ = builder.WriteString(usageFormat)
	writeProgParams(fs, builder)
	writeDescription(fs, builder)
	writeCmds(fs, builder)
	writeOpts(fs, builder)
	writeEnvs(fs, builder)
	writeParams(fs, builder)
	return builder.String()
}
