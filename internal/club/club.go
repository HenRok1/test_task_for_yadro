package club

import (
	"errors"
	"time"
)

// type TableState struct {
// 	TableNumber int
// 	Client      *Client
// 	TotalHours  int
// 	TablePrice  int
// }

// type Client struct {
// 	Name     string
// 	Table    int
// 	QueuePos int
// }

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
