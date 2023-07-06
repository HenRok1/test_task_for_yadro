package main

import (
	"errors"
	"fmt"
	"sort"
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

func (c *Club) IsOpen(t time.Time) bool {
	return t.After(c.OpenTime) && t.Before(c.CloseTime)
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
	fmt.Println(c.OpenTime.Format(time.Stamp))
	for client := range c.CurrentClients {
		fmt.Printf("%s\n", client)
	}
	fmt.Println(c.CloseTime.Format(time.Stamp))

	tables := make([]int, 0)
	for tableNum := range c.TableOccupation {
		tables = append(tables, tableNum)
	}
	sort.Ints(tables)
	for _, tableNum := range tables {
		duration := c.TableOccupation[tableNum]
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		fmt.Printf("%d %02d:%02d\n", tableNum, c.TablePrice*hours, minutes)
	}
}
