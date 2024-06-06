package meta

const (
	GeneralErrorCode = "G00001"
)

const (
	StatusSuccess = "Success"
	StatusFailure = "Failure"
)

type Status struct {
	Status    string       `json:"status,omitempty"`
	Data      interface{}  `json:"data,omitempty"`
	Reason    StatusReason `json:"reason,omitempty"`
	Code      int32        `json:"code,omitempty"`
	Message   string       `json:"message,omitempty"`
	ErrorCode string       `json:"error_code,omitempty"`
}

type StatusReason string

const (
	StatusReasonUnknown       StatusReason = ""
	StatusReasonUnauthorized  StatusReason = "Unauthorized"
	StatusReasonForbidden     StatusReason = "Forbidden"
	StatusReasonNotFound      StatusReason = "NotFound"
	StatusReasonAlreadyExists StatusReason = "AlreadyExists"
	StatusReasonBadRequest    StatusReason = "BadRequest"
	StatusReasonInternalError StatusReason = "InternalError"
	StatusReasonInvalid       StatusReason = "Invalid"
)

type Status2 struct {
	Status    int         `json:"status"`
	ErrorCode string      `json:"error_code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}
