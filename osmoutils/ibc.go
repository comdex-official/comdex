package osmoutils

import (
	"fmt"
	"encoding/json"
	"context"
	"strings"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	grpc1 "github.com/gogo/protobuf/grpc"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"text/template"
	"reflect"
	"time"
	"strconv"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

type FieldReadLocation = bool

const paginationType = "*query.PageRequest"

var lastQueryModuleName string

type CustomFieldParserFn = func(arg string, flags *pflag.FlagSet) (valueToSet any, usedArg FieldReadLocation, err error)

type FlagAdvice struct {
	HasPagination bool

	// Map of FieldName -> FlagName
	CustomFlagOverrides map[string]string
	CustomFieldParsers  map[string]CustomFieldParserFn

	// Tx sender value
	IsTx              bool
	TxSenderFieldName string
	FromValue         string
}

type FlagDesc struct {
	RequiredFlags []*pflag.FlagSet
	OptionalFlags []*pflag.FlagSet
}

type QueryDescriptor struct {
	Use   string
	Short string
	Long  string

	HasPagination bool

	QueryFnName string

	Flags FlagDesc
	// Map of FieldName -> FlagName
	CustomFlagOverrides map[string]string
	// Map of FieldName -> CustomParseFn
	CustomFieldParsers map[string]CustomFieldParserFn

	ParseQuery func(args []string, flags *pflag.FlagSet) (proto.Message, error)

	ModuleName string
	numArgs    int
}

type LongMetadata struct {
	BinaryName    string
	CommandPrefix string
	Short         string

	// Newline Example:
	ExampleHeader string
}

const IbcAcknowledgementErrorType = "ibc-acknowledgement-error"

// NewEmitErrorAcknowledgement creates a new error acknowledgement after having emitted an event with the
// details of the error.
func NewEmitErrorAcknowledgement(ctx sdk.Context, err error, errorContexts ...string) channeltypes.Acknowledgement {
	logger := ctx.Logger().With("module", IbcAcknowledgementErrorType)

	attributes := make([]sdk.Attribute, len(errorContexts)+1)
	attributes[0] = sdk.NewAttribute("error", err.Error())
	for i, s := range errorContexts {
		attributes[i+1] = sdk.NewAttribute("error-context", s)
		logger.Error(fmt.Sprintf("error-context: %v", s))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			IbcAcknowledgementErrorType,
			attributes...,
		),
	})

	return channeltypes.NewErrorAcknowledgement(err)
}

// IsAckError checks an IBC acknowledgement to see if it's an error.
// This is a replacement for ack.Success() which is currently not working on some circumstances
func IsAckError(acknowledgement []byte) bool {
	var ackErr channeltypes.Acknowledgement_Error
	if err := json.Unmarshal(acknowledgement, &ackErr); err == nil && len(ackErr.Error) > 0 {
		return true
	}
	return false
}


// Index command, but short is not set. That is left to caller.
func IndexCmd(moduleName string) *cobra.Command {
	return &cobra.Command{
		Use:                        moduleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", moduleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       indexRunCmd,
	}
}

func indexRunCmd(cmd *cobra.Command, args []string) error {
	usageTemplate := `Usage:{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}
  
{{if .HasAvailableSubCommands}}Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
	cmd.SetUsageTemplate(usageTemplate)
	return cmd.Help()
}

func QueryIndexCmd(moduleName string) *cobra.Command {
	cmd := IndexCmd(moduleName)
	cmd.Short = fmt.Sprintf("Querying commands for the %s module", moduleName)
	lastQueryModuleName = moduleName
	return cmd
}

func NewLongMetadata(moduleName string) *LongMetadata {
	commandPrefix := fmt.Sprintf("$ %s q %s", version.AppName, moduleName)
	return &LongMetadata{
		BinaryName:    version.AppName,
		CommandPrefix: commandPrefix,
	}
}

func ParseUint(arg string, fieldName string) (uint64, error) {
	v, err := strconv.ParseUint(arg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse %s as uint for field %s: %w", arg, fieldName, err)
	}
	return v, nil
}

func ParseFloat(arg string, fieldName string) (float64, error) {
	v, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse %s as float for field %s: %w", arg, fieldName, err)
	}
	return v, nil
}

func ParseInt(arg string, fieldName string) (int64, error) {
	v, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse %s as int for field %s: %w", arg, fieldName, err)
	}
	return v, nil
}

func ParseDenom(arg string, fieldName string) (string, error) {
	return strings.TrimSpace(arg), nil
}

func ParseFieldFromArg(fVal reflect.Value, fType reflect.StructField, arg string) error {
	// We cant pass in a negative number due to the way pflags works...
	// This is an (extraordinarily ridiculous) workaround that checks if a negative int is encapsulated in square brackets,
	// and if so, trims the square brackets
	if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") && arg[1] == '-' {
		arg = strings.TrimPrefix(arg, "[")
		arg = strings.TrimSuffix(arg, "]")
	}

	switch fType.Type.Kind() {
	// SetUint allows anyof type u8, u16, u32, u64, and uint
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		u, err := ParseUint(arg, fType.Name)
		if err != nil {
			return err
		}
		fVal.SetUint(u)
		return nil
	// SetInt allows anyof type i8,i16,i32,i64 and int
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		typeStr := fType.Type.String()
		var i int64
		var err error
		if typeStr == "time.Duration" {
			dur, err2 := time.ParseDuration(arg)
			i, err = int64(dur), err2
		} else {
			i, err = ParseInt(arg, fType.Name)
		}
		if err != nil {
			return err
		}
		fVal.SetInt(i)
		return nil
	case reflect.Float32, reflect.Float64:
		typeStr := fType.Type.String()
		f, err := ParseFloat(arg, typeStr)
		if err != nil {
			return err
		}
		fVal.SetFloat(f)
		return nil
	case reflect.String:
		s, err := ParseDenom(arg, fType.Name)
		if err != nil {
			return err
		}
		fVal.SetString(s)
		return nil
	case reflect.Ptr:
	case reflect.Slice:
		typeStr := fType.Type.String()
		if typeStr == "[]uint64" {
			// Parse comma-separated uint64 values into []uint64 slice
			strValues := strings.Split(arg, ",")
			values := make([]uint64, len(strValues))
			for i, strValue := range strValues {
				u, err := strconv.ParseUint(strValue, 10, 64)
				if err != nil {
					return err
				}
				values[i] = u
			}
			fVal.Set(reflect.ValueOf(values))
			return nil
		}
		if typeStr == "types.Coins" {
			coins, err := ParseCoins(arg, fType.Name)
			if err != nil {
				return err
			}
			fVal.Set(reflect.ValueOf(coins))
			return nil
		}
	case reflect.Struct:
		typeStr := fType.Type.String()
		var v any
		var err error
		if typeStr == "types.Coin" {
			v, err = ParseCoin(arg, fType.Name)
		} else if typeStr == "types.Int" {
			v, err = ParseSdkInt(arg, fType.Name)
		} else if typeStr == "time.Time" {
			v, err = ParseUnixTime(arg, fType.Name)
		} else if typeStr == "types.Dec" {
			v, err = ParseSdkDec(arg, fType.Name)
		} else {
			return fmt.Errorf("struct field type not recognized. Got type %v", fType)
		}

		if err != nil {
			return err
		}
		fVal.Set(reflect.ValueOf(v))
		return nil
	}
	fmt.Println(fType.Type.Kind().String())
	return fmt.Errorf("field type not recognized. Got type %v", fType)
}

func ParseSdkDec(arg, fieldName string) (sdk.Dec, error) {
	i, err := sdk.NewDecFromStr(arg)
	if err != nil {
		return sdk.Dec{}, fmt.Errorf("could not parse %s as sdk.Dec for field %s: %w", arg, fieldName, err)
	}
	return i, nil
}

func ParseUnixTime(arg string, fieldName string) (time.Time, error) {
	timeUnix, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		parsedTime, err := time.Parse(sdk.SortableTimeFormat, arg)
		if err != nil {
			return time.Time{}, fmt.Errorf("could not parse %s as time for field %s: %w", arg, fieldName, err)
		}

		return parsedTime, nil
	}
	startTime := time.Unix(timeUnix, 0)
	return startTime, nil
}

func ParseSdkInt(arg string, fieldName string) (sdk.Int, error) {
	i, ok := sdk.NewIntFromString(arg)
	if !ok {
		return sdk.Int{}, fmt.Errorf("could not parse %s as sdk.Int for field %s", arg, fieldName)
	}
	return i, nil
}

func ParseCoin(arg string, fieldName string) (sdk.Coin, error) {
	coin, err := sdk.ParseCoinNormalized(arg)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("could not parse %s as sdk.Coin for field %s: %w", arg, fieldName, err)
	}
	return coin, nil
}

func ParseCoins(arg string, fieldName string) (sdk.Coins, error) {
	coins, err := sdk.ParseCoinsNormalized(arg)
	if err != nil {
		return sdk.Coins{}, fmt.Errorf("could not parse %s as sdk.Coins for field %s: %w", arg, fieldName, err)
	}
	return coins, nil
}

func parseFieldFromDirectlySetFlag(fVal reflect.Value, fType reflect.StructField, flagAdvice FlagAdvice, flagName string, flags *pflag.FlagSet) error {
	// get string. If its a string great, run through arg parser. Otherwise try setting directly
	s, err := flags.GetString(flagName)
	if err != nil {
		flag := flags.Lookup(flagName)
		if flag == nil {
			return fmt.Errorf("Programmer set the flag name wrong. Flag %s does not exist", flagName)
		}
		t := flag.Value.Type()
		if t == "uint64" {
			u, err := flags.GetUint64(flagName)
			if err != nil {
				return err
			}
			fVal.SetUint(u)
			return nil
		}
	}
	return ParseFieldFromArg(fVal, fType, s)
}

// ParseFieldFromFlag attempts to parses the value of a field in a struct from a flag.
// The field is identified by the provided `reflect.StructField`.
// The flag advice and `pflag.FlagSet` are used to determine the flag to parse the field from.
// If the field corresponds to a value from a flag, true is returned.
// Otherwise, `false` is returned.
// In the true case, the parsed value is set on the provided `reflect.Value`.
// An error is returned if there is an issue parsing the field from the flag.
func ParseFieldFromFlag(fVal reflect.Value, fType reflect.StructField, flagAdvice FlagAdvice, flags *pflag.FlagSet) (bool, error) {
	lowercaseFieldNameStr := strings.ToLower(fType.Name)
	if flagName, ok := flagAdvice.CustomFlagOverrides[lowercaseFieldNameStr]; ok {
		return true, parseFieldFromDirectlySetFlag(fVal, fType, flagAdvice, flagName, flags)
	}

	kind := fType.Type.Kind()
	switch kind {
	case reflect.String:
		if flagAdvice.IsTx {
			// matchesFieldName is true if lowercaseFieldNameStr is the same as TxSenderFieldName,
			// or if TxSenderFieldName is left blank, then matches fields named "sender" or "owner"
			matchesFieldName := (flagAdvice.TxSenderFieldName == lowercaseFieldNameStr) ||
				(flagAdvice.TxSenderFieldName == "" && (lowercaseFieldNameStr == "sender" || lowercaseFieldNameStr == "owner"))
			if matchesFieldName {
				fVal.SetString(flagAdvice.FromValue)
				return true, nil
			}
		}
	case reflect.Ptr:
		if flagAdvice.HasPagination {
			typeStr := fType.Type.String()
			if typeStr == paginationType {
				pageReq, err := client.ReadPageRequest(flags)
				if err != nil {
					return true, err
				}
				fVal.Set(reflect.ValueOf(pageReq))
				return true, nil
			}
		}
	}
	return false, nil
}

// ParseField parses field #fieldIndex from either an arg or a flag.
// Returns true if it was parsed from an argument.
// Returns error if there was an issue in parsing this field.
func ParseField(v reflect.Value, t reflect.Type, fieldIndex int, arg string, flagAdvice FlagAdvice, flags *pflag.FlagSet) (bool, error) {
	fVal := v.Field(fieldIndex)
	fType := t.Field(fieldIndex)
	// fmt.Printf("Field %d: %s %s %s\n", fieldIndex, fType.Name, fType.Type, fType.Type.Kind())

	lowercaseFieldNameStr := strings.ToLower(fType.Name)
	if parseFn, ok := flagAdvice.CustomFieldParsers[lowercaseFieldNameStr]; ok {
		v, usedArg, err := parseFn(arg, flags)
		if err == nil {
			fVal.Set(reflect.ValueOf(v))
		}
		return usedArg, err
	}

	parsedFromFlag, err := ParseFieldFromFlag(fVal, fType, flagAdvice, flags)
	if err != nil {
		return false, err
	}
	if parsedFromFlag {
		return false, nil
	}
	return true, ParseFieldFromArg(fVal, fType, arg)
}

func MakeNew[T any]() T {
	var v T
	if typ := reflect.TypeOf(v); typ.Kind() == reflect.Ptr {
		elem := typ.Elem()
		//nolint:forcetypeassert
		return reflect.New(elem).Interface().(T) // must use reflect
	} else {
		return *new(T) // v is not ptr, alloc with new
	}
}

// makes an exception, where it allows Pagination to come from flags.
func ParseFieldsFromFlagsAndArgs[reqP any](flagAdvice FlagAdvice, flags *pflag.FlagSet, args []string) (reqP, error) {
	req := MakeNew[reqP]()
	v := reflect.ValueOf(req).Elem()
	t := v.Type()

	argIndexOffset := 0
	// Iterate over the fields in the struct
	for i := 0; i < t.NumField(); i++ {
		arg := ""
		if len(args) > i+argIndexOffset {
			arg = args[i+argIndexOffset]
		}
		usedArg, err := ParseField(v, t, i, arg, flagAdvice, flags)
		if err != nil {
			return req, err
		}
		if !usedArg {
			argIndexOffset -= 1
		}
	}
	return req, nil
}

func BuildQueryCli[reqP proto.Message, querier any](desc *QueryDescriptor, newQueryClientFn func(grpc1.ClientConn) querier) *cobra.Command {
	prepareDescriptor[reqP](desc)
	if desc.ParseQuery == nil {
		desc.ParseQuery = func(args []string, fs *pflag.FlagSet) (proto.Message, error) {
			flagAdvice := FlagAdvice{
				HasPagination:       desc.HasPagination,
				CustomFlagOverrides: desc.CustomFlagOverrides,
				CustomFieldParsers:  desc.CustomFieldParsers,
			}.Sanitize()
			return ParseFieldsFromFlagsAndArgs[reqP](flagAdvice, fs, args)
		}
	}

	cmd := &cobra.Command{
		Use:   desc.Use,
		Short: desc.Short,
		Long:  desc.Long,
		Args:  cobra.ExactArgs(desc.numArgs),
		RunE:  queryLogic(desc, newQueryClientFn),
	}
	flags.AddQueryFlagsToCmd(cmd)
	AddFlags(cmd, desc.Flags)
	if desc.HasPagination {
		cmdName := strings.Split(desc.Use, " ")[0]
		flags.AddPaginationFlagsToCmd(cmd, cmdName)
	}

	return cmd
}

func queryLogic[querier any](desc *QueryDescriptor,
	newQueryClientFn func(grpc1.ClientConn) querier,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientQueryContext(cmd)
		if err != nil {
			return err
		}
		queryClient := newQueryClientFn(clientCtx)

		req, err := desc.ParseQuery(args, cmd.Flags())
		if err != nil {
			return err
		}

		res, err := callQueryClientFn(cmd.Context(), desc.QueryFnName, req, queryClient)
		if err != nil {
			return err
		}

		return clientCtx.PrintProto(res)
	}
}

func callQueryClientFn(ctx context.Context, fnName string, req proto.Message, q any) (res proto.Message, err error) {
	qVal := reflect.ValueOf(q)
	method := qVal.MethodByName(fnName)
	if (method == reflect.Value{}) {
		return nil, fmt.Errorf("Method %s does not exist on the querier."+
			" You likely need to override QueryFnName in your Query descriptor", fnName)
	}
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(req),
	}
	results := method.Call(args)
	if len(results) != 2 {
		panic("We got something wrong")
	}
	if !results[1].IsNil() {
		//nolint:forcetypeassert
		err = results[1].Interface().(error)
		return res, err
	}
	//nolint:forcetypeassert
	res = results[0].Interface().(proto.Message)
	return res, nil
}

func (f FlagAdvice) Sanitize() FlagAdvice {
	// map CustomFlagOverrides & CustomFieldParser keys to lower-case
	// initialize if uninitialized
	newFlagOverrides := make(map[string]string, len(f.CustomFlagOverrides))
	for k, v := range f.CustomFlagOverrides {
		newFlagOverrides[strings.ToLower(k)] = v
	}
	f.CustomFlagOverrides = newFlagOverrides
	newFlagParsers := make(map[string]CustomFieldParserFn, len(f.CustomFieldParsers))
	for k, v := range f.CustomFieldParsers {
		newFlagParsers[strings.ToLower(k)] = v
	}
	f.CustomFieldParsers = newFlagParsers
	return f

}

// Required flags are marked as required.
func AddFlags(cmd *cobra.Command, desc FlagDesc) {
	for i := 0; i < len(desc.OptionalFlags); i++ {
		cmd.Flags().AddFlagSet(desc.OptionalFlags[i])
	}
	for i := 0; i < len(desc.RequiredFlags); i++ {
		fs := desc.RequiredFlags[i]
		cmd.Flags().AddFlagSet(fs)

		// mark all these flags as required.
		fs.VisitAll(func(flag *pflag.Flag) {
			err := cmd.MarkFlagRequired(flag.Name)
			if err != nil {
				panic(err)
			}
		})
	}
}

func ParseHasPagination[reqP any]() bool {
	req := MakeNew[reqP]()
	t := reflect.ValueOf(req).Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		fType := t.Field(i)
		if fType.Type.String() == paginationType {
			return true
		}
	}
	return false
}

func prepareDescriptor[reqP proto.Message](desc *QueryDescriptor) {
	if !desc.HasPagination {
		desc.HasPagination = ParseHasPagination[reqP]()
	}
	if desc.QueryFnName == "" {
		desc.QueryFnName = ParseExpectedQueryFnName[reqP]()
	}
	if strings.Contains(desc.Long, "{") {
		if desc.ModuleName == "" {
			desc.ModuleName = lastQueryModuleName
		}
		desc.FormatLong(desc.ModuleName)
	}

	desc.numArgs = ParseNumFields[reqP]() - len(desc.CustomFlagOverrides)
	if desc.HasPagination {
		desc.numArgs = desc.numArgs - 1
	}
}

func ParseNumFields[reqP any]() int {
	req := MakeNew[reqP]()
	v := reflect.ValueOf(req).Elem()
	t := v.Type()
	return t.NumField()
}

func (desc *QueryDescriptor) FormatLong(moduleName string) {
	desc.Long = FormatLongDesc(desc.Long, NewLongMetadata(moduleName).WithShort(desc.Short))
}

func FormatLongDesc(longString string, meta *LongMetadata) string {
	template, err := template.New("long_description").Parse(longString)
	if err != nil {
		panic("incorrectly configured long message")
	}
	bld := strings.Builder{}
	meta.ExampleHeader = "\n\nExample:"
	err = template.Execute(&bld, meta)
	if err != nil {
		panic("incorrectly configured long message")
	}
	return strings.TrimSpace(bld.String())
}

func (m *LongMetadata) WithShort(short string) *LongMetadata {
	m.Short = short
	return m
}

func GetParams[reqP proto.Message, querier any](moduleName string,
	newQueryClientFn func(grpc1.ClientConn) querier,
) *cobra.Command {
	return BuildQueryCli[reqP](&QueryDescriptor{
		Use:         "params [flags]",
		Short:       fmt.Sprintf("Get the params for the x/%s module", moduleName),
		QueryFnName: "Params",
	}, newQueryClientFn)
}

func ParseExpectedQueryFnName[reqP any]() string {
	req := MakeNew[reqP]()
	v := reflect.ValueOf(req).Elem()
	s := v.Type().String()
	// handle some non-std queries
	var prefixTrimmed string
	if strings.Contains(s, "Query") {
		prefixTrimmed = strings.Split(s, "Query")[1]
	} else {
		prefixTrimmed = strings.Split(s, ".")[1]
	}
	suffixTrimmed := strings.TrimSuffix(prefixTrimmed, "Request")
	return suffixTrimmed
}
