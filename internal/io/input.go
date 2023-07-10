package io

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// func ParseFile(filePath string) (numTables int, openTime, closeTime time.Time, payCost int, err error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Не удалось открыть файл:", err)
// 		return
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)

// 	numTables, err = ReadTables(file, scanner)
// 	if err != nil {
// 		return
// 	}

// }

func ReadTables(file *os.File, scanner *bufio.Scanner) (numTables int) {
	scanner.Scan()

	numTables, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal("Ошибка в первой строке: ", scanner.Text())
		// fmt.Println("Ошибка формата в первой строке", numTables)
		// return 0
	}
	return numTables
}

func ReadTime(file *os.File, scanner *bufio.Scanner) (openCloseTimes []string) {
	scanner.Scan()
	openCloseTimes = strings.Split(scanner.Text(), " ")
	if len(openCloseTimes) != 2 {
		log.Fatal("Ошибка во второй строке: ", scanner.Text())

		// fmt.Println("Ошибка формата во второй строке", openCloseTimes)
		// return
	}
	return openCloseTimes
}

func ParseOpenCloseTime(openCloseTimes []string) (openTime, closeTime time.Time) {
	openTime, err := time.Parse("15:04", openCloseTimes[0])
	if err != nil {
		log.Fatal("Ошибка формата во второй строке", openCloseTimes[0])
		// fmt.Println("Ошибка формата во второй строке")
		// return
	}

	// fmt.Printf("Время открытия клуба: %v\n", openTime)

	closeTime, err = time.Parse("15:04", openCloseTimes[1])
	if err != nil {
		log.Fatal("Ошибка формата во второй строке", openCloseTimes[1])
		// return
	}

	return openTime, closeTime
}

func ReadCost(file *os.File, scanner *bufio.Scanner) (payCost int) {
	scanner.Scan()
	payCost, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal("Ошибка формата в третьей строке", scanner.Text())

		// fmt.Println("Ошибка в третьей строке", payCost)
		// return
	}
	return payCost
}
