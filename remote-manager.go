package main

import (
	"fmt"
	"remote-manager/config"
)

func main() {
	c := config.Config()
	c.SaveConfig()
	fmt.Println(c)
}
