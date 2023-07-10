package club

import (
	"bufio"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Club struct {
	Tables          int
	OpenTime        time.Time
	CloseTime       time.Time
	TablePrice      int
	CurrentClients  map[string]bool
	WaitingQueue    []string
	TableOccupation map[int]time.Duration
	Revenue         int
}

func NewClub(tables int, openTime, closeTime time.Time, tablePrice int) *Club {
	return &Club{
		Tables:          tables,
		OpenTime:        openTime,
		CloseTime:       closeTime,
		TablePrice:      tablePrice,
		CurrentClients:  make(map[string]bool),
		WaitingQueue:    make([]string, 0),
		TableOccupation: make(map[int]time.Duration),
		Revenue:         0,
	}
}

func (c *Club) HandleClientArrival(t time.Time, name string) error {
	if _, ok := c.CurrentClients[name]; ok {
		return errors.New("YouShallNotPass")
	}
	if !c.IsOpen(t) {
		return errors.New("NotOpenYet")
	}

	c.CurrentClients[name] = true
	return nil
}

func (c *Club) IsOpen(timestamp time.Time) bool {
	return timestamp.After(c.OpenTime) && timestamp.Before(c.CloseTime)
}

func (c *Club) HandleClientSeat(t time.Time, name string, tableNumber int) error {
	if _, ok := c.CurrentClients[name]; !ok {
		return errors.New("ClientUnknown")
	}
	if c.TableOccupation[tableNumber] > 0 && c.TableOccupation[tableNumber] <= c.CloseTime.Sub(c.OpenTime) {
		return errors.New("PlaceIsBusy")
	}

	if c.TableOccupation[tableNumber] == 0 {
		duration := t.Sub(c.OpenTime)
		c.TableOccupation[tableNumber] = duration
	} else {
		duration := t.Sub(c.OpenTime) - c.TableOccupation[tableNumber]
		c.TableOccupation[tableNumber] += duration
	}

	c.CurrentClients[name] = true
	return nil
}

func (c *Club) HandleClientWait(t time.Time, name string) error {
	if len(c.WaitingQueue) >= c.Tables {
		c.HandleClientLeave(t, name)
		return errors.New("ICanWaitNoLonger")
	}
	c.WaitingQueue = append(c.WaitingQueue, name)
	return nil
}

func (c *Club) HandleClientLeave(t time.Time, name string) error {
	if _, ok := c.CurrentClients[name]; !ok {
		return errors.New("ClientUnknown")
	}
	delete(c.CurrentClients, name)

	for tableNum, occupiedDuration := range c.TableOccupation {
		if occupiedDuration == 0 {
			if len(c.WaitingQueue) > 0 {
				clientName := c.WaitingQueue[0]
				c.WaitingQueue = c.WaitingQueue[1:]

				c.HandleClientSeat(t, clientName, tableNum)
				break
			}
		}
	}

	return nil
}

func (c *Club) CalculateRevenue() {
	for _, duration := range c.TableOccupation {
		hours := int(duration.Hours())
		if duration-(time.Duration(hours)*time.Hour) > 0 {
			hours++
		}
		c.Revenue += hours * c.TablePrice
	}
}

func (c *Club) PrintClubStatus() {
	fmt.Println(c.OpenTime.Format(time.TimeOnly))
	for val := range c.CurrentClients {
		fmt.Printf("%s\n", val)
	}
	defer fmt.Println(c.CloseTime.Format(time.Stamp))

	tables := make([]int, 0)
	for tableNum := range c.TableOccupation {
		tables = append(tables, tableNum)
	}
	sort.Ints(tables)
	for _, tableNum := range tables {
		duration := c.TableOccupation[tableNum]
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		fmt.Printf("%d %d %02d:%02d\n", tableNum, c.Revenue, c.TablePrice*hours, minutes)
	}
}

func (c *Club) HandleEvents(scanner *bufio.Scanner) {
	c.Tables = 0

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
		c.HandleEventCode(eventCode, event, t, line)
	}

	c.CalculateRevenue()
	// fmt.Println(club.Revenue)
	c.PrintClubStatus()
}

func (c *Club) HandleEventCode(evenCode int, event []string, t time.Time, line string) (err error) {
	switch evenCode {
	case 1: // Клиент пришел
		if len(event) < 3 {
			fmt.Printf("Ошибка формата на строке: %s\n", line)
			return
		}

		err = c.HandleClientArrival(t, event[2])
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

		err = c.HandleClientSeat(t, event[2], tableNum)
		if err != nil {
			fmt.Printf("%s %d %s\n", event[0], 13, err)
		}
	case 3: // Клиент ожидает
		if len(event) < 3 {
			fmt.Printf("Ошибка формата на строке: %s\n", line)
			return
		}

		err = c.HandleClientWait(t, event[2])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	case 4: // Клиент ушел
		if len(event) < 3 {
			fmt.Printf("Ошибка формата на строке: %s\n", line)
			return
		}

		err = c.HandleClientLeave(t, event[2])
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	default:
		fmt.Printf("Неизвестный код события на строке: %s\n", line)
		return
	}
	return nil
}
