module github.com/Muhammadpurwanto/ecommerce-baju/services/gateway

go 1.25.0

require (
	github.com/Muhammadpurwanto/ecommerce-baju/services/common v0.0.0-00010101000000-000000000000
	github.com/gofiber/fiber/v2 v2.52.13
	github.com/gofiber/storage/redis/v3 v3.4.8
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/redis/go-redis/v9 v9.20.0
	github.com/spf13/viper v1.21.0
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.81.1
)

require (
	github.com/klauspost/cpuid/v2 v2.2.11 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/bytedance/gopkg v0.1.4
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.23 // indirect
	github.com/pelletier/go-toml/v2 v2.3.1 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tinylib/msgp v1.6.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/Muhammadpurwanto/ecommerce-baju/services/auth => ../auth

replace github.com/Muhammadpurwanto/ecommerce-baju/services/user => ../user

replace github.com/Muhammadpurwanto/ecommerce-baju/services/product => ../product

replace github.com/Muhammadpurwanto/ecommerce-baju/services/cart => ../cart

replace github.com/Muhammadpurwanto/ecommerce-baju/services/order => ../order

replace github.com/Muhammadpurwanto/ecommerce-baju/services/payment => ../payment

replace github.com/Muhammadpurwanto/ecommerce-baju/services/common => ../common
