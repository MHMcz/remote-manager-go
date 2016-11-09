/*
	Copyright © 2016 Jan Markup <mhmcze@gmail.com>
	This work is free. You can redistribute it and/or modify it under the
	terms of the Do What The Fuck You Want To Public License, Version 2,
	as published by Sam Hocevar. See the COPYING file for more details.
*/

package ncurses

import (
	gc "github.com/rthornton128/goncurses"
	"remote-manager/config"
	"strings"
)

func Run(c *config.Configuration) {
	stdscr, _ := gc.Init()
	defer gc.End()

	gc.StartColor()
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)

	gc.InitPair(1, gc.C_WHITE, gc.C_BLUE)

	reprintGroups(c, stdscr)

	for {
		gc.Update()
		ch := stdscr.GetChar()
		switch ch {
		case 27:
			return
		case gc.KEY_F2:
			createNewGroup(c, stdscr)
			reprintGroups(c, stdscr)
		}
	}
}

func reprintGroups(c *config.Configuration, window *gc.Window) {
	window.Clear()

	for _, group := range c.Groups {
		window.Println(group.Name)
		for _, remote := range group.Remotes {
			window.Println("--- " + remote.Name)
		}
	}

	window.Println("")
	window.Println("New group with F2.")
	window.Println("New remote with F3.")
	window.Println("Exit with ESC.")
}

func createNewGroup(c *config.Configuration, stdscr *gc.Window) {
	gc.Cursor(1)

	rows, cols := stdscr.MaxYX()
	width, height := 60, 15
	y, x := (rows-height)/2, (cols-width)/2

	var win *gc.Window
	win = stdscr.Sub(height, width, y, x)
	win.Keypad(true)
	var winForm *gc.Window
	winForm = win.Derived(10, width, 0, 0)
	winForm.Keypad(true)
	var winMenu *gc.Window
	winMenu = win.Derived(5, width, 10, 0)
	winMenu.Keypad(true)

	fields := make([]*gc.Field, 5)
	for i := 0; i < 3; i++ {
		fields[i], _ = gc.NewField(1, 25, int32(4+(i*2)), 25, 0, 0)
		defer fields[i].Free()
		fields[i].SetBackground(gc.A_UNDERLINE)
		fields[i].SetOptionsOff(gc.FO_AUTOSKIP)
		fields[i].SetOptionsOff(gc.FO_STATIC)
	}

	form, _ := gc.NewForm(fields)
	form.SetSub(winForm)
	form.Post()
	defer form.UnPost()
	defer form.Free()

	winForm.MovePrint(4, 10, "Name:")
	winForm.MovePrint(6, 10, "Alias:")
	winForm.MovePrint(8, 10, "Alias for mc:")
	win.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	win.Refresh()
	winForm.Refresh()

	menu_items := []string{"CANCEL", "OK"}
	items := make([]*gc.MenuItem, len(menu_items))
	for i, val := range menu_items {
		items[i], _ = gc.NewItem(val, "")
		defer items[i].Free()
	}

	menu, _ := gc.NewMenu(items)
	menu.SetWindow(winMenu)
	menu.Format(1, 2)
	defer menu.Free()
	menu.Post()
	winMenu.Refresh()

	form.Driver(gc.REQ_FIRST_FIELD)

	ch := winForm.GetChar()
formloop:
	for ch != 27 { // ESC
		switch ch {
		case gc.KEY_ENTER, 10: // enter
			winMenu.Refresh()
		menuloop:
			for {
				ch := winForm.GetChar()
				switch ch {
				case gc.KEY_TAB:
					winForm.Refresh()
					break menuloop
				case gc.KEY_ENTER, 10: // enter
					activeIndex := menu.Current(nil).Index()
					if activeIndex == 0 {
						break formloop
					} else if activeIndex == 1 {
						err := form.Driver(gc.REQ_VALIDATION)
						if err == nil {
							group := config.GroupConfig{}
							group.Name = strings.TrimSpace(fields[0].Buffer())
							group.Alias = strings.TrimSpace(fields[1].Buffer())
							group.AliasMc = strings.TrimSpace(fields[2].Buffer())
							group.Remotes = make([]config.RemoteConfig, 0)
							c.Groups = append(c.Groups, group)
							c.SaveConfig()
							break formloop
						}
					}
				default:
					menu.Driver(gc.DriverActions[ch])
				}
				winMenu.Refresh()
			}
		case gc.KEY_DOWN, gc.KEY_TAB:
			form.Driver(gc.REQ_NEXT_FIELD)
			form.Driver(gc.REQ_END_LINE)
		case gc.KEY_UP:
			form.Driver(gc.REQ_PREV_FIELD)
			form.Driver(gc.REQ_END_LINE)
		case gc.KEY_LEFT:
			form.Driver(gc.REQ_PREV_CHAR)
		case gc.KEY_RIGHT:
			form.Driver(gc.REQ_NEXT_CHAR)
		case gc.KEY_BACKSPACE:
			form.Driver(gc.REQ_DEL_PREV)
		case gc.KEY_DC:
			form.Driver(gc.REQ_DEL_CHAR)
		default:
			form.Driver(ch)
		}
		ch = winForm.GetChar()
		stdscr.Refresh()
	}
	gc.Cursor(0)
}
