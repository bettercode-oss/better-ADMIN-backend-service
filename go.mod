module better-admin-backend-service

go 1.18

require (
	github.com/bettercode-oss/rest v0.0.4
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.8.1
	github.com/go-ldap/ldap/v3 v3.3.0
	github.com/go-testfixtures/testfixtures/v3 v3.5.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/configor v1.2.1
	github.com/keepeye/logrus-filename v0.0.0-20190711075016-ce01a4391dd1
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.8.1
	github.com/wesovilabs/koazee v0.0.5
	golang.org/x/crypto v0.4.0
	gorm.io/driver/mysql v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.9
	gorm.io/plugin/dbresolver v1.1.0
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/avast/retry-go v3.0.0+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.9.0 // indirect
	github.com/ernesto-jimenez/httplogger v0.0.0-20150224132909-86cc44f6150a // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/jackc/pgx/v4 v4.10.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.1+incompatible
