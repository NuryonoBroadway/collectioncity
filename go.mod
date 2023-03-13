module collection-squad/collection/collection-city

go 1.16

require (
	cloud.google.com/go/storage v1.27.0
	github.com/Shopify/sarama v1.38.1
	github.com/aws/aws-sdk-go v1.43.26
	github.com/brainlabs/snowflake v0.3.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pressly/goose/v3 v3.5.3
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.8.1
	github.com/xdg/scram v1.0.5
	github.com/ziutek/mymysql v1.5.4
	gitlab.privy.id/privypass/privypass-package-core v1.7.6
	google.golang.org/api v0.103.0
	google.golang.org/grpc v1.51.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.37.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/pierrec/lz4 v2.6.1+incompatible // indirect
