module github.com/bhbosman/goCommsStacks

go 1.24.0

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20250308000247-4306925b3dfd
	github.com/bhbosman/goMessages v0.0.0-20250308063841-56e83c4ed9eb
	github.com/bhbosman/gocommon v0.0.0-20250308052839-0ebeb121f996
	github.com/bhbosman/gocomms v0.0.0-20250308053431-cfd3e830d8a5
	github.com/bhbosman/goerrors v0.0.0-20250307194237-312d070c8e38
	github.com/bhbosman/gomessageblock v0.0.0-20250307141417-ab783e8e2eba
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/gobwas/ws v1.4.0
	github.com/reactivex/rxgo/v2 v2.5.0
	github.com/stretchr/testify v1.10.0
	go.uber.org/fx v1.23.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.37.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/bhbosman/goConnectionManager v0.0.0-20250308054322-587c786bedb0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/icza/gox v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	go.uber.org/dig v1.18.1 // indirect
	golang.org/x/sys v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20250308071159-4cf72f668c72

replace github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af
