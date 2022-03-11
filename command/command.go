package command

import "errors"

var FailExecError = errors.New("failed to execute command, please troubleshoot logs")
