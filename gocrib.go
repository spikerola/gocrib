package main

import (
    "fmt"
    "log"
    "time"
    "encoding/hex"
    "github.com/jroimartin/gocui"
)

func xorTest(a, b byte) bool {
    c := a^b
    return c >= 20 && c <= 126
}

func xor(a, b []byte) []byte {
    if len(a) != len(b) {
        panic("different lenght")
    }
    c := make([]byte, len(a))
    for i := 0; i < len(a); i++ {
        c[i] = a[i] ^ b[i]
    }
    return c
}

func crib(a [][]byte, key []byte, alphabet string) string {
    key = key[:len(key)-1]
    var res string = ""
    if len(a[0]) <= len(key) {
        // if we already know the key
        res = "DONE!\n"
        key = key[:len(a[0])]
        for _, b := range a {
            r := xor(b[:len(key)], key)
            res += fmt.Sprintf("%s\n", r)
        }
        return res
    }

    // try each character in the alphabet
    for _, k := range alphabet {
        var z bool = false
        // check if it's meaningful for each ciphertexts
        for _, b := range a {
            z = z || !xorTest(b[len(key)], byte(k))
        }
        // otherwise skip the character
        if z {
            continue
        }

        // append the results
        res += fmt.Sprintf("\n= %c =\n", k)
        for _, b := range a {
            r := xor(b[:len(key)+1], append(key, byte(k)))
            res += fmt.Sprintf("%s\n", r)
        }
    }
    return res
}

func main() {
    fmt.Print("Enter key alphabet: ")
    var keyAlphabet string
    fmt.Scanf("%s", &keyAlphabet)

    cipherTexts := make([][]byte, 0)

    // take ciphertexts, exit on empty string
    for {
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

    // text update function
    go func() {
        var tmp string
        for {
            resView, err := g.View("results")
            if err != nil {
                continue
            }
            inputView, err := g.View("fast")
            if err != nil {
                continue
            }
            in := inputView.Buffer()
            if in != tmp { // only on text change
                resView.Clear()
                time.Sleep(100 * time.Millisecond)
                fmt.Fprintf(resView, "%s", crib(cipherTexts, []byte(in), keyAlphabet))
                tmp = in
            }
            time.Sleep(100 * time.Millisecond)
        }
    }()

// Gui setup
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

// Layout

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
