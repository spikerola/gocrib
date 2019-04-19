package main

import (
    "fmt"
    "log"
    "time"
    "encoding/hex"
    "github.com/jroimartin/gocui"
)

func main() {
    fmt.Print("Enter key alphabet: ")
    var keyAlphabet string
    fmt.Scanf("%s", &keyAlphabet)

    cipherTexts := make([][]byte, 0)
    for ;; {
        fmt.Print("Enter ciphertext: ")
        var t string
        fmt.Scanf("%s", &t)
        if t == "" {
            break
        }

        b, err := hex.DecodeString(t)
        if err != nil {
            fmt.Println("hex string not accepted")
            continue
        }

        cipherTexts = append(cipherTexts, b)
    }

    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

    go func() {
        var tmp string
        for ;; {
            resView, err := g.View("results")
            if err != nil {
                continue
            }
            inputView, err := g.View("fast")
            if err != nil {
                continue
            }
            in := inputView.Buffer()
            if in != tmp {
                resView.Clear()
                time.Sleep(100 * time.Millisecond)
                fmt.Fprintf(resView, "%s\n%s\n%s\n", in, in, in)
                tmp = in
            }
            time.Sleep(100 * time.Millisecond)
        }
    }()

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
        func(g *gocui.Gui, v *gocui.View) error {
            resView, err := g.View("results")
            if err != nil {
                return err
            }
            scrollView(resView, -1)
            return nil
        }); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
        func(g *gocui.Gui, v *gocui.View) error {
            resView, err := g.View("results")
            if err != nil {
                panic(err)
            }
            scrollView(resView, 1)
            return nil
        }); err != nil {
		log.Panicln(err)
	}

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
    if v, err := g.SetView("results", 0, 0, maxX-1, maxY-7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
        v.Wrap = true
    }

	if v, err := g.SetView("fast", 0, maxY-7, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
        v.Editable = true
        v.Wrap = true
        if _, err = setCurrentViewOnTop(g, "fast"); err != nil {
            return err
        }
        g.Cursor = true
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    // TODO
    resView, err := g.View("results")
    if err != nil {
        return err
    }
    inputView, err := g.View("fast")
    if err != nil {
        return err
    }
    log.Printf("%s", resView.Buffer())
    log.Printf("%s", inputView.Buffer())
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}


func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}