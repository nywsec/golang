package dictionnary

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

var words = make([]string, 0, 50)

func Load(filename string) error {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		words = append(words, scanner.Text())

	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func PickWord() string {

	rand.Seed(time.Now().Unix())
	i := rand.Intn(len(words))
	return words[i]
}
