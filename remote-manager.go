/*
	Copyright © 2016 Jan Markup <mhmcze@gmail.com>
	This work is free. You can redistribute it and/or modify it under the
	terms of the Do What The Fuck You Want To Public License, Version 2,
	as published by Sam Hocevar. See the COPYING file for more details.
*/

package main

import (
	"remote-manager/config"
	"remote-manager/frontend/ncurses"
)

func main() {
	c := config.Config()
	ncurses.Run(c)

	return
}
