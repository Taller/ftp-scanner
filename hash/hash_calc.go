package hash

import (
	"crypto/sha256"
	"io"
	"log"
	"os"
)

func CalcHash(fileName string) string {

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%x", h.Sum(nil))
	return string(h.Sum(nil)[:])
}
