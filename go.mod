module github.com/shopcx/frontend

go 1.21

require (
	github.com/gin-gonic/gin v1.7.7 // Intentionally using older version with known vulnerabilities
	github.com/go-sql-driver/mysql v1.5.0 // Contains SQL injection vulnerabilities
	github.com/golang-jwt/jwt v3.2.1+incompatible // Using older version with known vulnerabilities
	github.com/gorilla/sessions v1.2.1 // Contains session fixation vulnerabilities
	github.com/swaggo/swag v1.8.1
	github.com/swaggo/gin-swagger v1.4.3
	google.golang.org/grpc v1.45.0
	github.com/gin-contrib/sessions v0.0.5 // Contains session management vulnerabilities
	github.com/gin-contrib/cors v1.3.1 // Contains CORS misconfiguration vulnerabilities
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.34.1-0.20250224173925-7292932d45d5 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.1.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
) 