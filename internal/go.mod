module github.com/Tamplier2911/gorest/internal

go 1.16

replace github.com/Tamplier2911/gorest/pkg => ../pkg

require (
	github.com/Tamplier2911/gorest/pkg v0.0.0-20210815190001-755e69443498
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.1.2
	github.com/labstack/echo/v4 v4.5.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/echo-swagger v1.1.2
	github.com/swaggo/swag v1.7.1
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/oauth2 v0.0.0-20210810183815-faf39c7919d5
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5 // indirect
	gorm.io/gorm v1.21.12
)
