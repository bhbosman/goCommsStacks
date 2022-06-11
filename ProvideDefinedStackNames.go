package goCommsStacks

import (
	bottomStack "github.com/bhbosman/goCommsStacks/bottom/common"
	bvisMessageBreakerStack "github.com/bhbosman/goCommsStacks/bvisMessageBreaker/common"
	messageCompressorStack "github.com/bhbosman/goCommsStacks/messageCompressor/common"
	messageNumberStack "github.com/bhbosman/goCommsStacks/messageNumber/common"
	pingPongStack "github.com/bhbosman/goCommsStacks/pingPong/common"
	tlsConnectionStack "github.com/bhbosman/goCommsStacks/tlsConnection/common"
	topStack "github.com/bhbosman/goCommsStacks/top/common"
	webSocketStack "github.com/bhbosman/goCommsStacks/websocket/common"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
)

func ProvideDefinedStackNames() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotated{
			Group: "TransportFactory",
			Target: common.NewTransportFactory(
				common.TransportFactoryCompressedTlsName,
				topStack.StackName,
				pingPongStack.StackName,
				messageCompressorStack.StackName,
				messageNumberStack.StackName,
				bvisMessageBreakerStack.StackName,
				tlsConnectionStack.StackName,
				bottomStack.StackName)}),
		fx.Provide(fx.Annotated{
			Group: "TransportFactory",
			Target: common.NewTransportFactory(
				common.TransportFactoryUnCompressedTlsName,
				topStack.StackName,
				pingPongStack.StackName,
				messageNumberStack.StackName,
				bvisMessageBreakerStack.StackName,
				tlsConnectionStack.StackName,
				bottomStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryCompressedName,
					topStack.StackName,
					pingPongStack.StackName,
					messageCompressorStack.StackName,
					messageNumberStack.StackName,
					bvisMessageBreakerStack.StackName,
					bottomStack.StackName,
				)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryUnCompressedName,
					topStack.StackName,
					pingPongStack.StackName,
					messageNumberStack.StackName,
					bvisMessageBreakerStack.StackName,
					bottomStack.StackName,
				)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryEmptyName,
					topStack.StackName,
					bottomStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.WebSocketName,
					topStack.StackName,
					webSocketStack.StackName,
					bottomStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.EmptyStackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryOnlyBottomStack,
					bottomStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryOnlyBvisMessageBreakerStack,
					bvisMessageBreakerStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryOnlyTlsConnectionStack,
					tlsConnectionStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryOnlyMessageNumberStack,
					messageNumberStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryOnlyMessageCompressorStack,
					messageCompressorStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryOnlyPingPongStack,
					pingPongStack.StackName)}),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: common.NewTransportFactory(
					common.TransportFactoryOnlyTopStack,
					topStack.StackName)}),
	)
}
