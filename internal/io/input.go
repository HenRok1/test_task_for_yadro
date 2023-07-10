package io

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadTables(file *os.File, scanner *bufio.Scanner) (numTables int) {
	scanner.Scan()

	numTables, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal("Ошибка в первой строке: ", scanner.Text())

	}
	return numTables
}

func ReadTime(file *os.File, scanner *bufio.Scanner) (openCloseTimes []string) {
	scanner.Scan()
	openCloseTimes = strings.Split(scanner.Text(), " ")
	if len(openCloseTimes) != 2 {
		log.Fatal("Ошибка во второй строке: ", scanner.Text())
	}
	return openCloseTimes
}

func ParseOpenCloseTime(openCloseTimes []string) (openTime, closeTime time.Time) {
	openTime, err := time.Parse("15:04", openCloseTimes[0])
	if err != nil {
		log.Fatal("Ошибка формата во второй строке", openCloseTimes[0])
	}
	closeTime, err = time.Parse("15:04", openCloseTimes[1])
	if err != nil {
		log.Fatal("Ошибка формата во второй строке", openCloseTimes[1])
	}

	return openTime, closeTime
}

func ReadCost(file *os.File, scanner *bufio.Scanner) (payCost int) {
	scanner.Scan()
	payCost, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal("Ошибка формата в третьей строке", scanner.Text())
	}
	return payCost
}
