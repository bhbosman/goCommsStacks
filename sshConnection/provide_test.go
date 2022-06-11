package sshConnection

import (
	"fmt"
	"github.com/bhbosman/goCommsStacks/sshConnection/common"
	"github.com/bhbosman/gocommon/Services/IConnectionManager"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/fx/ApplicationContext"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	common2 "github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netListener"
	"github.com/cskr/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestProvideSshStack(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	connectionReactor := intf.NewMockIConnectionReactor(controller)
	connectionReactor.EXPECT().OnInitReturnDefault(
		func(b bool, i interface{}) {
			fmt.Print(reflect.TypeOf(t).String())
		}, nil).AnyTimes()
	connectionReactor.EXPECT().OnOpenReturn(nil).AnyTimes()
	connectionReactor.EXPECT().OnCloseReturn(nil).AnyTimes()

	crf := intf.NewMockIConnectionReactorFactory(controller)
	crf.EXPECT().OnCreateReturnDefault(connectionReactor).AnyTimes()

	connectionManager := IConnectionManager.NewMockIService(controller)
	connectionManager.EXPECT().OnStateReturn(IFxService.Started).AnyTimes()
	connectionManager.EXPECT().OnServiceNameReturn("ABC").AnyTimes()
	connectionManager.EXPECT().OnRegisterConnectionReturnDefault(nil).AnyTimes()
	connectionManager.EXPECT().OnDeregisterConnectionReturnDefault(nil).AnyTimes()
	connectionManager.EXPECT().OnConnectionInformationReceivedReturnDefault(nil).AnyTimes()

	onConnection := common2.NewMockIOnCreateConnection(controller)
	onConnection.EXPECT().OnOnCreateConnectionDoDefault(
		func(arg0 string, arg1 error, arg2 context.Context, arg3 context.CancelFunc) {
			fmt.Println("Dddd")
		}).AnyTimes()

	uniqueReferenceService := interfaces.NewMockIUniqueReferenceService(controller)
	uniqueReferenceService.EXPECT().OnNextReturnDefault("SomeRef").AnyTimes()

	var AppCallback messages.CreateAppCallback
	fxApp := fxtest.New(t,
		fx.Populate(&AppCallback),
		fx.Provide(func() interfaces.IUniqueReferenceService { return uniqueReferenceService }),
		fx.Provide(func() IConnectionManager.IService { return connectionManager }),
		ApplicationContext.ProvideApplicationContext(),
		fx.Provide(func() *zap.Logger { return zap.NewNop() }),
		PubSub.ProvidePubSub("Application"),
		fx.Provide(
			fx.Annotated{
				Target: func(params struct {
					fx.In
					PubSub             *pubsub.PubSub `name:"Application"`
					NetAppFuncInParams common2.NetAppFuncInParams
				}) messages.CreateAppCallback {
					f := netListener.NewNetListenApp(
						"",
						0,
						0,
						"Prefix",
						"tcp4://127.0.0.1:8888",
						common2.TransportFactoryOnlySSHStack,
						func() (intf.IConnectionReactorFactory, error) {
							return crf, nil
						},
						common2.MaxConnectionsSetting(1),
						common2.NewOverrideOnCreateConnectionFactory(
							func() (common2.IOnCreateConnection, error) {
								return onConnection, nil
							}),
						netListener.NewOverrideListenerAcceptFactory(
							func(listener net.Listener) (netListener.IListenerAccept, error) {
								return &ListenAcceptorForTesting{
									listener: listener,
									onAccept: nil,
								}, nil
							}),
						common2.NewConnectionInstanceOptions(
							fx.Provide(
								fx.Annotated{
									Group:  "TransportFactory",
									Target: common2.NewTransportFactory(common2.TransportFactoryOnlySSHStack, common.StackName)}),
						),
						common2.NewConnectionInstanceOptions(
							ProvideSshCommunicationStack(),
						),
					)
					return f(params.NetAppFuncInParams)
				},
			},
		))
	assert.NoError(t, fxApp.Err())
	assert.NoError(t, fxApp.Start(context.Background()))

	sut, err := AppCallback.Callback()
	assert.NoError(t, err)
	assert.NotNil(t, sut)

	if !assert.NoError(t, sut.Start(context.Background())) {
		return
	}

	time.Sleep(time.Hour)
}
