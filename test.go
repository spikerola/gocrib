package main

import (
    "fmt"
    "encoding/hex"
)

func main() {
    k := []byte("FLAG{THIS_15_TH3_K3Y___YEP_WE_DO_THIS}")
    fmt.Println(len(k))
    /*
    a := []byte("Succulents cred skateboard tousled pok pok pabst hell of meh shoreditch offal gluten-free viral. Palo santo bicycle rights kinfolk, af pok pok jean shorts iPhone hoodie normcore poke 90's gastropub. Slow-carb you probably haven't heard of them unicorn distillery meh readymade bespoke vinyl. Gentrify quinoa tacos semiotics bitters skateboard craft beer, master cleanse umami keffiyeh. Microdosing small batch cold-pressed kinfolk, humblebrag bushwick direct trade organic chillwave asymmetrical mumblecore chicharrones. Vape lomo portland, drinking")
    for i := 0; i < len(a); i++ {
        if i % len(k) == 0 {
            fmt.Printf("\n")
        }
        fmt.Printf("%c", a[i])
    }
    fmt.Println("--")
    */
    a := []byte("Succulents cred skateboard tousled pok")
    b := []byte(" pok pabst hell of meh shoreditch offa")
    c := []byte("l gluten-free viral. Palo santo bicycl")
    d := []byte("e rights kinfolk, af pok pok jean shor")
    e := []byte("ts iPhone hoodie normcore poke 90's ga")
    f := []byte("stropub. Slow-carb you probably haven'")

    t := [][]byte{a, b, c, d, e, f}
    for _, x := range t {
        fmt.Printf("%s\n", hex.EncodeToString(xor(x, k)))
    }
}

func xor(a, b []byte) []byte {
    if len(a) != len(b) {
        panic("different length")
    }
    c := make([]byte, len(a))
    for i := 0; i < len(a); i++ {
        c[i] = a[i] ^ b[i]
    }
    return c
}

