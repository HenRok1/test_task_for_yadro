package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Не указан путь к файлу")
		return
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Не удалось открыть файл:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Считываем количество столов
	scanner.Scan()
	tables, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка формата в первой строке")
		return
	}

	// Считываем время начала и окончания работы клуба
	scanner.Scan()
	openCloseTimes := strings.Split(scanner.Text(), " ")
	if len(openCloseTimes) != 2 {
		fmt.Println("Ошибка формата во второй строке")
		return
	}
	openTime, err := time.Parse("15:04", openCloseTimes[0])
	if err != nil {
		fmt.Println("Ошибка формата во второй строке")
		return
	}
	closeTime, err := time.Parse("15:04", openCloseTimes[1])
	if err != nil {
		fmt.Println("Ошибики формата в третьей строке")
		return
	}
}
