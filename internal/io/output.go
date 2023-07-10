package io

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/HenRok1/test_task_for_yadro/internal/club"
)

func HandleEvents(scanner *bufio.Scanner, club *club.Club) {
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
		HandleEventCode(eventCode, event, t, line, club)
	}

	club.CalculateRevenue()
	// fmt.Println(club.Revenue)
	club.PrintClubStatus()
}

func HandleEventCode(evenCode int, event []string, t time.Time, line string, club *club.Club) (err error) {
	switch evenCode {
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
			return nil
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
	return nil
}
