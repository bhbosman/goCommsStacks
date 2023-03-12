package protoBuf

import (
	"fmt"
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/zap"
	"reflect"
	"strconv"
)

// TODO: Create Constructors
type stackHandler struct {
	errorState error
	logger     *zap.Logger
	counterMap map[reflect.Type]int
	prefix     string
}

func (self *stackHandler) addCounter(of reflect.Type) {
	counter, ok := self.counterMap[of]
	if ok {
		self.counterMap[of] = counter + 1
	} else {
		self.counterMap[of] = 1
	}
}

func (self *stackHandler) ReadMessage(i interface{}) error {
	if self.errorState != nil {
		return self.errorState
	}
	switch v := i.(type) {
	case *model.ClearCounters:
		self.counterMap = make(map[reflect.Type]int)
	case *model.PublishRxHandlerCounters:
		for r, i := range self.counterMap {
			v.AddMapData(fmt.Sprintf("ProtoBuf %v %v", self.prefix, r.String()), strconv.Itoa(i))
		}
		return nil
	}
	return nil
}
