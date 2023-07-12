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

const errorNum = 13

type Club struct {
	Tables          int
	OpenTime        time.Time
	CloseTime       time.Time
	TablePrice      int
	CurrentClients  map[string]bool
	WaitingQueue    []string
	TableOccupation map[int]time.Duration
	Revenue         int
	TableFree       map[int]bool
	ClientTable     map[string]int
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
		TableFree:       make(map[int]bool),
		ClientTable:     map[string]int{},
	}
}

func (c *Club) HandleClientArrival(t time.Time, name string) (errNum int, err error) {
	if _, ok := c.CurrentClients[name]; ok {
		return errorNum, errors.New("YouShallNotPass")
	}
	if !c.IsOpen(t) {
		return errorNum, errors.New("NotOpenYet")
	}

	c.CurrentClients[name] = true
	return 0, nil
}

func (c *Club) HandleClientSeat(t time.Time, name string, tableNumber int) (errNum int, err error) {
	if _, ok := c.CurrentClients[name]; !ok {
		return errorNum, errors.New("ClientUnknown")
	}
	if !c.TableFree[tableNumber] && c.TableOccupation[tableNumber] <= c.CloseTime.Sub(c.OpenTime) {
		return errorNum, errors.New("PlaceIsBusy")
	}
	//TODO: NEED TO FIX
	if c.TableOccupation[tableNumber] == 0 {
		duration := t.Sub(c.OpenTime)
		c.TableOccupation[tableNumber] = duration
		// fmt.Printf("Duration %v: %v\n", tableNumber, duration)

	} else {
		duration := t.Sub(c.OpenTime) - c.TableOccupation[tableNumber]
		// fmt.Println("Duration 2: ", duration)
		c.TableOccupation[tableNumber] += duration
	}

	c.TableFree[tableNumber] = false

	c.ClientTable[name] = tableNumber

	c.Tables -= 1
	c.CurrentClients[name] = true
	return 0, nil
}

// TODO: переделать
func (c *Club) HandleClientWait(t time.Time, name string) (errNum int, err error) {
	if c.Tables > 0 {
		// c.HandleClientLeave(t, name)
		return errorNum, errors.New("ICanWaitNoLonger")
	}
	c.WaitingQueue = append(c.WaitingQueue, name)
	return 0, nil
}

func (c *Club) HandleClientLeave(t time.Time, name string) error {
	if _, ok := c.CurrentClients[name]; !ok {
		return errors.New("ClientUnknown")
	}
	delete(c.CurrentClients, name)

	c.TableFree[c.ClientTable[name]] = true
	c.ClientTable[name] = 0

	for tableNum := range c.TableOccupation {
		if c.TableFree[tableNum] {
			if len(c.WaitingQueue) > 0 {
				clientName := c.WaitingQueue[0]
				c.WaitingQueue = c.WaitingQueue[1:]

				c.HandleClientSeat(t, clientName, tableNum)
				fmt.Printf("%s %d %s %d\n", t.Format(time.TimeOnly)[:5], 12, clientName, tableNum)
				c.Tables -= 1
				break
			}
		}
	}
	c.Tables += 1

	return nil
}

func (c *Club) HandleLastClient(t time.Time, name string) error {
	for c.ClientTable[name] != 0 {
		err := c.HandleClientLeave(t, name)
		if err != nil {
			return err
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
	fmt.Println(c.CloseTime.Format(time.TimeOnly)[:5])

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
	fmt.Println("c.Tables = ", c.Tables)
	// c.Tables = 0

	for table := 0; table < c.Tables; table++ {
		c.TableFree[table+1] = true
	}

	for scanner.Scan() {
		line := scanner.Text()
		event := strings.Fields(line)
		// Проверка формата входных данных
		if len(event) < 3 {
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
		errNum, err := c.HandleClientArrival(t, event[2])
		if err != nil {
			fmt.Printf("%s %d %s\n", event[0], errNum, err)
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

		errNum, err := c.HandleClientSeat(t, event[2], tableNum)
		if err != nil {
			fmt.Printf("%s %d %s\n", event[0], errNum, err)
		}
	case 3: // Клиент ожидает
		errNum, err := c.HandleClientWait(t, event[2])
		if err != nil {
			fmt.Printf("%s %d %s\n", event[0], errNum, err)
		}
	case 4: // Клиент ушел
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

func (c *Club) IsOpen(timestamp time.Time) bool {
	return timestamp.After(c.OpenTime) && timestamp.Before(c.CloseTime)
}
