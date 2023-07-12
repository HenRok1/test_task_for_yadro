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
	Tables          int //Количество столов
	OpenTime        time.Time
	CloseTime       time.Time
	TablePrice      int
	CurrentClients  map[string]bool //Находится ли в клубе клиент
	WaitingQueue    []string
	TableOccupation map[int]time.Duration //Количество времени за столом
	Revenue         map[int]int
	TableFree       map[int]bool
	ClientTable     map[string]int
	StartTableUse   map[int]time.Time
	EndTableUse     map[int]time.Time
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
		Revenue:         make(map[int]int),
		TableFree:       make(map[int]bool),
		ClientTable:     make(map[string]int),
		StartTableUse:   make(map[int]time.Time),
		EndTableUse:     make(map[int]time.Time),
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

	c.StartTableUse[tableNumber] = t

	c.TableFree[tableNumber] = false

	c.ClientTable[name] = tableNumber

	c.Tables -= 1
	c.CurrentClients[name] = true
	return 0, nil
}

func (c *Club) HandleClientWait(t time.Time, name string) (errNum int, err error) {
	if c.Tables > 0 {
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

	c.EndTableUse[c.ClientTable[name]] = t

	c.TableOccupation[c.ClientTable[name]] += c.EndTableUse[c.ClientTable[name]].Sub(c.StartTableUse[c.ClientTable[name]])

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

func (c *Club) HandleLastClient(t time.Time, name string) {
	delete(c.CurrentClients, name)

	c.TableOccupation[c.ClientTable[name]] = t.Sub(c.StartTableUse[c.ClientTable[name]])

	c.TableFree[c.ClientTable[name]] = true

	c.ClientTable[name] = 0

	c.Tables += 1

	c.WaitingQueue = c.WaitingQueue[:0]

	fmt.Printf("%s %d %s\n", c.CloseTime.Format(time.TimeOnly)[:5], 11, name)
}

func (c *Club) CalculateRevenue() {
	for tableNum, duration := range c.TableOccupation {
		hours := int(duration.Hours())
		if duration-(time.Duration(hours)*time.Hour) > 0 {
			hours++
		}
		c.Revenue[tableNum] += hours * c.TablePrice
	}
}

func (c *Club) PrintClubRevenue() {
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
		fmt.Printf("%d %d %02d:%02d\n", tableNum, c.Revenue[tableNum], hours, minutes)
	}
}

func (c *Club) HandleEvents(scanner *bufio.Scanner) {

	fmt.Println(c.OpenTime.Format(time.TimeOnly)[:5])

	for table := 0; table < c.Tables; table++ {
		c.TableFree[table+1] = true
	}

	for scanner.Scan() {
		line := scanner.Text()
		event := strings.Fields(line)
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

	names := make([]string, 0)

	for name := range c.CurrentClients {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, nameVal := range names {
		c.HandleLastClient(c.CloseTime, nameVal)
	}

	c.CalculateRevenue()

	c.PrintClubRevenue()
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
