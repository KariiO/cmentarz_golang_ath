package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

// Magazyn
// 100 zniczy
// 50 wiazanek

// 4 babki
// 2 babki - znicze
// 2 babki - wiazanki
// babka - jedna sztuka

// znicze i wiazanki w koszach pobieranie przez poslancow
// kosz na znicze pojenosc 10, na wiazanki 10

// 5 poslancow
// maksymalnie 1 wiazanke, 2 znicze

var start_magazyn_znicze = 100
var start_magazyn_wiazanki = 50
var magazyn_znicze = start_magazyn_znicze
var magazyn_wiazanki = start_magazyn_wiazanki
var max_ilosc_babek_znicze = 2
var max_ilosc_babek_wiazanki = 2
var max_ilosc_poslancow = 5
var max_kosz_na_znicze = 10
var max_kosz_na_wiazanki = 10
var kosz_na_znicze = 0
var kosz_na_wiazanki = 0
var poslaniec_max_wiazanki = 1
var poslaniec_max_znicze = 2

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/magazyn", magazyn)
	mux.HandleFunc("/kosz", kosz)
	mux.HandleFunc("/babka/znicze/", babkaZnicze)
	mux.HandleFunc("/babka/wiazanki/", babkaWiazanki)
	mux.HandleFunc("/poslaniec/", poslaniec)
	log.Fatal(http.ListenAndServe(":3000", recoverMw(mux, true)))
}

func recoverMw(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				if !dev {
					http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, string(stack))
			}
		}()

		nw := &responseWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		nw.flush()
	}
}

// type ResponseWriter interface {
// 	Header() Header
// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

type responseWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter does not support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
}

func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}

func magazyn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,
		"Stan magazynu:\nZnicze: ",
		magazyn_znicze,
		"\nWiazanki: ",
		magazyn_wiazanki,
	)
}
func kosz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,
		"Stan kosza:\nZnicze: ",
		kosz_na_znicze,
		"\nWiazanki: ",
		kosz_na_wiazanki,
	)
}

func babkaZnicze(w http.ResponseWriter, r *http.Request) {
	babka_numer, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/babka/znicze/"))

	if err != nil {
		panic("Blad! Nie mozna pobrac numeru babki!")
	}
	if babka_numer > max_ilosc_babek_znicze {
		msg := fmt.Sprintf("Blad! Maksymalna ilosc babek od zniczy to: %d", max_ilosc_babek_znicze)
		panic(msg)
	}
	if magazyn_znicze == 0 {
		panic("Blad! Magazyn nie posiada wiecej zniczy!")
	}
	if kosz_na_znicze == max_kosz_na_znicze {
		msg := fmt.Sprintf("Blad! Kosz na znicze jest pelny, maksymalna ilosc to: %d", max_kosz_na_znicze)
		panic(msg)
	}

	kosz_na_znicze = kosz_na_znicze + 1
	magazyn_znicze = magazyn_znicze - 1

	fmt.Fprintln(w, "Babka", babka_numer, "pobrala znicz", math.Abs(float64(magazyn_znicze-start_magazyn_znicze)), "z magazynu do kosza.")
}

func babkaWiazanki(w http.ResponseWriter, r *http.Request) {
	babka_numer, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/babka/wiazanki/"))

	if err != nil {
		panic("Blad! Nie mozna pobrac numeru babki!")
	}
	if babka_numer > max_ilosc_babek_wiazanki {
		msg := fmt.Sprintf("Blad! Maksymalna ilosc babek od wiazanek to: %d", max_ilosc_babek_wiazanki)
		panic(msg)
	}
	if magazyn_wiazanki == 0 {
		panic("Blad! Magazyn nie posiada wiecej wiazanek!")
	}
	if kosz_na_wiazanki == max_kosz_na_wiazanki {
		msg := fmt.Sprintf("Blad! Kosz na wiazanki jest pelny, maksymalna ilosc to: %d", max_kosz_na_wiazanki)
		panic(msg)
	}

	kosz_na_wiazanki = kosz_na_wiazanki + 1
	magazyn_wiazanki = magazyn_wiazanki - 1

	fmt.Fprintln(w, "Babka", babka_numer, "pobrala wiazanke", math.Abs(float64(magazyn_wiazanki-start_magazyn_wiazanki)), "z magazynu do kosza.")
}

func poslaniec(w http.ResponseWriter, r *http.Request) {
	poslaniec_numer, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/poslaniec/"))

	if err != nil {
		panic("Blad! Nie mozna pobrac numeru babki!")
	}
	if poslaniec_numer > max_ilosc_poslancow {
		msg := fmt.Sprintf("Blad! Maksymalna ilosc poslancow to: %d", max_ilosc_poslancow)
		panic(msg)
	}
	if kosz_na_wiazanki == 0 && kosz_na_znicze == 0 {
		panic("Blad! Kosz jest pusty!")
	}
	msg := fmt.Sprintf("");
	if kosz_na_znicze >= poslaniec_max_znicze {
		kosz_na_znicze = kosz_na_znicze - poslaniec_max_znicze
		msg = msg + fmt.Sprintf("Poslaniec %d pobiera %d znicze", poslaniec_numer, poslaniec_max_znicze)
	} else {
		msg = msg + fmt.Sprintf("Poslaniec %d pobiera %d znicze", poslaniec_numer, kosz_na_znicze)
		kosz_na_znicze = 0
	}

	if kosz_na_wiazanki >= poslaniec_max_wiazanki {
		kosz_na_wiazanki = kosz_na_wiazanki - poslaniec_max_wiazanki
		msg = msg + fmt.Sprintf("\nPoslaniec %d pobiera %d wiazanki", poslaniec_numer, poslaniec_max_wiazanki)
	} else {
		msg = msg + fmt.Sprintf("\nPoslaniec %d pobiera %d wiazanki", poslaniec_numer, kosz_na_wiazanki)
		kosz_na_wiazanki = 0
	}

	fmt.Fprintln(w, msg)
}
