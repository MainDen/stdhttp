package flagx

type Handler interface {
	Parse(args ...string) (int, error)
	Params() string
}

type Parser interface {
	Parse(args ...string) (int, error)
}

type handlerParams struct {
	params string
	parser Parser
}

func Params(params string, parser Parser) *handlerParams {
	return &handlerParams{params: params, parser: parser}
}

func (handler *handlerParams) Parse(args ...string) (int, error) {
	return handler.parser.Parse(args...)
}

func (handler *handlerParams) Params() string {
	return handler.params
}

func (fs *FlagSet) SetDefaultHandlerParams(params string, parser Parser) {
	fs.defaultHandler = Params(params, parser)
}

func SetDefaultHandlerParams(params string, parser Parser) {
	CommandLine.SetDefaultHandlerParams(params, parser)
}

type nopHandler struct{}

func (*nopHandler) Parse(args ...string) (int, error) {
	return 0, nil
}

func (*nopHandler) Params() string {
	return ""
}

func Nop() *nopHandler {
	return &nopHandler{}
}
