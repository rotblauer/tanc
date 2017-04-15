package Utils

import (
	"bufio"
	"log"
	"os"
)

func LoadRsId(rsIdFile string) map[string]string {
	m := make(map[string]string)

	file, err := os.Open(rsIdFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tmp := scanner.Text()
		m[tmp] = tmp
		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return m

}
