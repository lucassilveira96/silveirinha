package utils

import (
	"fmt"
	"strconv"
	"teste/internal/infra/variables"
)

func ShowBanner() {
	cyan := "\033[0;36m"
	white := "\033[1;37m"
	port := strconv.Itoa(variables.ServerPort())

	fmt.Println(cyan)
	fmt.Println(white + " ____________________________________________________________________________________________________________________")
	fmt.Println(white + "|                                                                                                                    |")
	fmt.Println(cyan + "|  ######  #### ##       ##     ## ######## #### ########  #### ##    ## ##     ##    ###         #######   #######  |")
	fmt.Println(cyan + "| ##    ##  ##  ##       ##     ## ##        ##  ##     ##  ##  ###   ## ##     ##   ## ##       ##     ## ##     ## |")
	fmt.Println(cyan + "| ##        ##  ##       ##     ## ##        ##  ##     ##  ##  ####  ## ##     ##  ##   ##      ##        ##     ## |")
	fmt.Println(cyan + "|  ######   ##  ##       ##     ## ######    ##  ########   ##  ## ## ## ######### ##     ##     ##   #### ##     ## |")
	fmt.Println(cyan + "|       ##  ##  ##        ##   ##  ##        ##  ##   ##    ##  ##  #### ##     ## #########     ##    ##  ##     ## |")
	fmt.Println(cyan + "| ##    ##  ##  ##         ## ##   ##        ##  ##    ##   ##  ##   ### ##     ## ##     ## ### ##    ##  ##     ## |")
	fmt.Println(cyan + "|  ######  #### ########    ###    ######## #### ##     ## #### ##    ## ##     ## ##     ## ###  ######    #######  |")
	fmt.Println(white + "|____________________________________________________________________________________________________________________|")
	fmt.Println(white + "|                                                                                                                    |")
	fmt.Println(white + "|                This program automatically generates all the necessary files to create a CRUD                       |")
	fmt.Println(white + "|                    for your application in Go (Golang) using Fiber, GORM, Swagger                                  |")
	fmt.Println(white + "|                                        and several other dependencies.                                             |")
	fmt.Println(white + "|____________________________________________________________________________________________________________________|")
	fmt.Println(white + "|                                                                                                                    |")
	fmt.Println(white + "|                                                 GO CODE GENERATOR                                                  |")
	fmt.Println(white + "|                                                   version 1.0.0                                                    |")
	fmt.Println(white + "|____________________________________________________________________________________________________________________|")
	fmt.Println(white + "Project Name: " + variables.ServiceName())
	fmt.Println(white + "Version: " + variables.ServiceVersion())
	fmt.Println(white + "Project Started: " + variables.ServerHost() + ":" + port + "/" + variables.PrefixRoute())
	fmt.Println("")
	fmt.Println("")
}
