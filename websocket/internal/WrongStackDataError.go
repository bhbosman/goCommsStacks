package internal

import (
	"github.com/bhbosman/goCommsStacks/websocket/common"
	"github.com/bhbosman/gocommon/model"
	common2 "github.com/bhbosman/gocomms/common"
	"reflect"
)

func WrongStackDataError(connectionType model.ConnectionType, stackData interface{}) error {
	return common2.NewWrongStackDataType(
		common.StackName,
		connectionType,
		reflect.TypeOf((*StackData)(nil)),
		reflect.TypeOf(stackData))
}
