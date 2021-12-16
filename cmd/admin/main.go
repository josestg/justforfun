package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	logger := log.New(os.Stderr, "admin: ", log.LstdFlags)

	fs := flag.NewFlagSet("Admin Tool CLI", flag.ExitOnError)

	gen := fs.String("gen", "", "Output types")
	alg := fs.String("alg", "", "The algorithm for generating the output")
	src := fs.String("src", "", "Source input")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Printf("error: %v", err)
		return
	}

	switch *gen {
	default:
		fs.PrintDefaults()
	case "hash":
		switch *alg {
		case "bcrypt":
			in := *src
			if len(in) == 0 {
				logger.Print("`src` value is required")
				return
			}

			out, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}

			_, _ = fmt.Fprintf(os.Stdout, "%s\n", out)
		default:
			logger.Print("`alg` is required")
			return
		}
	}
}
