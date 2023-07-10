package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/HenRok1/test_task_for_yadro/internal/club"
	"github.com/HenRok1/test_task_for_yadro/internal/io"
)

func main() {
	if len(os.Args) != 2 { 
		fmt.Println("Не указан путь к файлу или некорректный ввод аргументов")
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

	myClub := club.NewClub(numTables, openTime, closeTime, payCost)

	fmt.Println(openCloseTimes[0])

	myClub.HandleEvents(scanner)
}
