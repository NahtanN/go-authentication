package main

import "github.com/nahtann/go-lab/cmd"

//	@title			Go-Lab API
//	@version		1.0
//	@description	Golang API service. It`s just for studies purposes.

//	@contact.name	NahtanN
//	@contact.url	https://www.linkedin.com/in/nahtann/
//	@contact.email	nahtann@outlook.com

//	@host		localhost:3333
//	@BasePath	/api

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Example: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
func main() {
	cmd.Execute()
}
