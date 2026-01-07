package errno

import "fmt"

type Error struct {
	HTTPCode int
	Code     int
	Msg      string
	Err      error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

// 　常用错误
var (
	ErrInternal              = &Error{HTTPCode: 500, Code: 10000, Msg: "Internal server error"}
	ErrDB                    = &Error{HTTPCode: 500, Code: 10001, Msg: "Database connection error"}
	ErrUnauthorized          = &Error{HTTPCode: 401, Code: 10002, Msg: "Unauthorized"}
	ErrPostNotFound          = &Error{HTTPCode: 404, Code: 10003, Msg: "Post not found"}
	ErrCommentNotFound       = &Error{HTTPCode: 404, Code: 10004, Msg: "Comment not found"}
	ErrUserExists            = &Error{HTTPCode: 400, Code: 10005, Msg: "User already exists"}
	ErrPostExists            = &Error{HTTPCode: 400, Code: 10005, Msg: "User already exists"}
	ErrInvalidParameter      = &Error{HTTPCode: 405, Code: 10006, Msg: "Invalid request parameter"}
	ErrStatusForbidden       = &Error{HTTPCode: 406, Code: 10007, Msg: "You are not the author of the post"}
	ErrBadPassword           = &Error{HTTPCode: 407, Code: 10008, Msg: "Failed to hash password"}
	ErrInvalidNameOrPassword = &Error{HTTPCode: 408, Code: 10009, Msg: "Invalid username or password"}
	ErrGenerateToken         = &Error{HTTPCode: 409, Code: 10010, Msg: "Failed to generate token"}
)

// 包装数据库错误
func DB(err error) *Error {
	return &Error{HTTPCode: 500, Code: 10001, Msg: "Datsbase error", Err: err}
}

// 包装内部错误
func Internal(err error) *Error {
	return &Error{HTTPCode: 500, Code: 10000, Msg: "Internal server error", Err: err}
}

// 包装内部错误
func InvalidParameter(err error) *Error {
	return &Error{HTTPCode: 500, Code: 10006, Msg: "Invalid request parameter", Err: err}
}
