package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/HenRok1/test_task_for_yadro/internal/club"
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
	numTables, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка формата в первой строке", numTables)
		return
	}
	fmt.Println(numTables)

	// Считываем время начала и окончания работы клуба
	scanner.Scan()
	openCloseTimes := strings.Split(scanner.Text(), " ")
	if len(openCloseTimes) != 2 {
		fmt.Println("Ошибка формата во второй строке", openCloseTimes)
		return
	}

	openTime, err := time.Parse("15:04", openCloseTimes[0])
	if err != nil {
		fmt.Println("Ошибка формата во второй строке")
		return
	}

	fmt.Println(openTime)

	closeTime, err := time.Parse("15:04", openCloseTimes[1])
	if err != nil {
		fmt.Println("Ошибики формата в третьей строке")
		return
	}

	fmt.Println(closeTime)

	scanner.Scan()
	payCost, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка в четвертой строке", payCost)
		return
	}
	fmt.Println(payCost)

	tables := make([]club.TableState, numTables)
	for i := range tables {
		tables[i].TableNumber = i + 1
	}

	for scanner.Scan() {
		line := scanner.Text()
		eventFields := strings.Split(line, " ")
		_, err := time.Parse("15:04", eventFields[0])
		if err != nil {
			fmt.Println("Ошибка в формате времени:", eventFields[0])
			return
		}

		if len(eventFields) != 3 {
			fmt.Println("Ошибка формата: неверное событие", eventFields)
			return
		}
		fmt.Println(eventFields)

		// timestamp := eventFields[0]
		// eventID, err := strconv.Atoi(eventFields[1])
		// if err != nil {
		// 	log.Fatal("Ошибка формата: неверный ID события")
		// }

		// switch eventID {
		// case 1:
		// 	// Обработка события "Клиент пришел"
		// 	clientName := eventFields[2]
		// 	handleClientArrival(clientName, tables, workHours, timestamp)
		// case 2:
		// 	// Обработка события "Клиент сел за стол"
		// 	clientName := eventFields[2]
		// 	tableNumber, err := strconv.Atoi(eventFields[3])
		// 	if err != nil {
		// 		log.Fatal("Ошибка формата: неверный номер стола")
		// 	}
		// 	handleClientSeated(clientName, tableNumber, tables, workHours, pricePerHour, timestamp)
		// case 3:
		// 	// Обработка события "Клиент ожидает"
		// 	clientName := eventFields[2]
		// 	handleClientWaiting(clientName, tables, timestamp)
		// case 4:
		// 	// Обработка события "Клиент ушел"
		// 	clientName := eventFields[2]
		// 	handleClientDeparture(clientName, tables, timestamp)
		// default:
		// 	log.Fatal("Ошибка формата: неизвестный ID события")
		// }
	}

}
