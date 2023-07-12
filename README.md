# Test task for Yadro

# Инструкция
1) Запустить docker
2) Прописать: ```make docker```
3) Прописать: ```make docker_run```

# Информация
Консольное приложение для обработки ивентов в компьютерном клубе на языке Go. \
В разработке использовался Golang 1.20.5. \
Была использована только стандартная библиотека.

## Что можно было бы исправить
- Разбить структура на 3 структуры: Table, Client и Club для просоты восприятия.

```go
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
```