package http2

// PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n
var CONNECTION_PREFACE = []byte{0x50, 0x52, 0x49, 0x20, 0x2a, 0x20, 0x48, 0x54, 0x54, 0x50, 0x2f, 0x32, 0x2e, 0x30, 0x0d, 0x0a, 0x0d, 0x0a, 0x53, 0x4d, 0x0d, 0x0a, 0x0d, 0x0a}

type FLAG byte

const (
	NO          FLAG = 0x00
	ACK              = 0x01
	END_STREAM       = 0x01
	END_HEADERS      = 0x04
	PADDED           = 0x08
	PRIORITY         = 0x20
)

var flagNames map[FLAG]string = map[FLAG]string{
	NO:          "NO",
	ACK:         "ACK or END",
	END_HEADERS: "END_HEADEDRS",
	PADDED:      "PADDED",
	PRIORITY:    "PRIORTY",
}

func (flag FLAG) String() string {
	return flagNames[flag]
}

type TYPE byte

const (
	DATA_FRAME TYPE = iota
	HEADERS_FRAME
	PRIORITY_FRAME
	RST_STREAM_FRAME
	SETTINGS_FRAME
	PING_FRAME
	GOAWAY_FRAME
	WINDOW_UPDATE_FRAME
	CONTINUATION_FRAME
)

var frameNames []string = []string{
	"DATA",
	"HEADERS",
	"PRIORITY",
	"RST_STREAM",
	"SETTINGS",
	"PING",
	"GOAWAY",
	"WINDOW_UPDATE",
	"CONTINUATION",
}

func (frameType TYPE) String() string {
	return frameNames[int(frameType)]
}

type STATE byte

const (
	IDLE = iota
	RESERVED_LOCAL
	RESERVED_REMOTE
	OPEN
	HALF_CLOSED_LOCAL
	HALF_CLOSED_REMOTE
	CLOSED
)

var streamState []string = []string{
	"IDLE",
	"RESERVED_LOCAL",
	"RESERVED_REMOTE",
	"OPEN",
	"HALF_CLOSED_LOCAL",
	"HALF_CLOSED_REMOTE",
	"CLOSED",
}

func (state STATE) String() string {
	return streamState[int(state)]
}

type SETTING uint16

const (
	NO_SETTING SETTING = iota
	HEADER_TABLE_SIZE
	ENABLE_PUSE
	MAX_CONCURRENT_STREAMS
	INITIAL_WINDOW_SIZE
	MAX_FRAME_SIZE
	MAX_HEADER_LIST_SIZE
)

var settingIDs []string = []string{
	"NO_SETTING",
	"HEADER_TABLE_SIZE",
	"ENABLE_PUSH",
	"MAX_CONCURRENT_STREAMS",
	"INITIAL_WINDOW_SIZE",
	"MAX_FRAME_SIZE",
	"MAX_HEADER_LIST_SIZE",
}

func (setting SETTING) String() string {
	return settingIDs[int(setting)]
}

type ERROR uint32

const (
	NO_ERROR ERROR = iota
	PROTOCOL_ERROR
	INTERNAL_ERROR
	FLOW_CONTROL_ERROR
	SETTINGS_TIMEOUT
	STREAM_CLOSED
	FRAME_SIZE_ERROR
	REFUSED_STREAM
	CANCEL
	COMPRESSION_ERROR
	CONNECT_ERROR
	ENHANCE_YOUR_CALM
	INADEQUATE_SECURITY
	HTTP_1_1_REQUIRED
)

var errorNames []string = []string{
	"NO_ERROR",
	"PROTOCOL_ERROR",
	"INTERNAL_ERROR",
	"FLOW_CONTROL_ERROR",
	"SETTINGS_TIMEOUT",
	"STREAM_CLOSED",
	"FRAME_SIZE_ERROR",
	"REFUSED_STREAM",
	"CANCEL",
	"COMPRESSION_ERROR",
	"CONNECT_ERROR",
	"ENHANCE_YOUR_CALM",
	"INADEQUATE_SECURITY",
	"HTTP_1_1_REQUIRED",
}

func (err ERROR) String() string {
	return errorNames[int(err)]
}
