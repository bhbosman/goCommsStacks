package internal

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"reflect"
)

func WrongStackDataError(connectionType model.ConnectionType, stackData interface{}) error {
	return common2.NewWrongStackDataType(
		goCommsDefinitions.WebSocketStackName,
		connectionType,
		reflect.TypeOf((*Data)(nil)),
		reflect.TypeOf(stackData))
}
