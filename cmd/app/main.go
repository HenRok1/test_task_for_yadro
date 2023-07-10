package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/HenRok1/test_task_for_yadro/internal/club"
	"github.com/HenRok1/test_task_for_yadro/internal/io"
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
	numTables := io.ReadTables(file, scanner)

	// Считываем время начала и окончания работы клуба
	openCloseTimes := io.ReadTime(file, scanner)

	openTime, closeTime := io.ParseOpenCloseTime(openCloseTimes)

	//Считываение стоимости часа в клубе
	payCost := io.ReadCost(file, scanner)

	club := club.NewClub(numTables, openTime, closeTime, payCost)

	/////////////////////////////////////////////////////////////////////////////////////

	fmt.Println(openCloseTimes[0])

	for scanner.Scan() {
		line := scanner.Text()
		event := strings.Fields(line)
		// Проверка формата входных данных
		if len(event) < 2 {
			fmt.Printf("Ошибка формата на строке: %s\n", line)
			return
		}

		t, err := time.Parse("15:04", event[0])
		if err != nil {
			fmt.Printf("Ошибка формата на строке: %s\n", line)
			return
		}

		eventCode, err := strconv.Atoi(event[1])
		if err != nil {
			fmt.Printf("Ошибка формата на строке: %s\n", line)
			return
		}
		fmt.Println(line)

		switch eventCode {
		case 1: // Клиент пришел
			if len(event) < 3 {
				fmt.Printf("Ошибка формата на строке: %s\n", line)
				return
			}

			err = club.HandleClientArrival(t, event[2])
			if err != nil {
				fmt.Printf("%s %d %s\n", event[0], 13, err)
			}
		case 2: // Клиент сел за стол
			if len(event) < 4 {
				fmt.Printf("Ошибка формата на строке: %s\n", line)
				return
			}

			tableNum, err := strconv.Atoi(event[3])
			if err != nil {
				fmt.Printf("Ошибка формата на строке: %s\n", line)
				return
			}

			err = club.HandleClientSeat(t, event[2], tableNum)
			if err != nil {
				fmt.Printf("%s %d %s\n", event[0], 13, err)
			}
		case 3: // Клиент ожидает
			if len(event) < 3 {
				fmt.Printf("Ошибка формата на строке: %s\n", line)
				return
			}

			err = club.HandleClientWait(t, event[2])
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		case 4: // Клиент ушел
			if len(event) < 3 {
				fmt.Printf("Ошибка формата на строке: %s\n", line)
				return
			}

			err = club.HandleClientLeave(t, event[2])
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		default:
			fmt.Printf("Неизвестный код события на строке: %s\n", line)
			return
		}
	}

	club.CalculateRevenue()
	// fmt.Println(club.Revenue)
	club.PrintClubStatus()
}
