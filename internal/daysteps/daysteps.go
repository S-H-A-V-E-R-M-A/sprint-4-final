package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// не экспортируемая
// по строке с данными получаю значения кол-ва шагов, время прогулки
func parsePackage(data string) (int, time.Duration, error) {
	values := strings.Split(data, ",") // разбил строку на слайс строк по запятой
	if len(values) != 2 {              // если меньше 2 строк в салйсе то ошибка
		return 0, 0, fmt.Errorf("неверный формат данных: %s", data)
	}

	counter, err := strconv.Atoi(values[0]) //конверчу первую строку из слайса
	if err != nil || counter <= 0 {         // если есть ошибка то нули и err
		return 0, 0, errors.New("ошибка парсинга")
	}

	duration, err := time.ParseDuration(values[1]) // конверчу строку вторую в формат time
	if err != nil || duration <= 0 {               // если есть ошибка то нули и err
		return 0, 0, errors.New("ошибка конвертации")
	}

	return counter, duration, nil // возврат данных из ф-ции кол-во шагов и время гуляния
}

// парсинг строки через parsePackage() + по весу/рост получаю строку с дистанцией, кол-вом шагов и пораченных ккал
func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data) //получаю шаги и время

	if err != nil { // если получил вместо шагов и времени ошибку то пустая строка и ошибка
		fmt.Println("ошибка парсинга ", err)
		return ""
	}

	if steps <= 0 { // если шагов меньше или равно нулю то ошибка и пустая стркоа
		fmt.Println("кол-во шаговменьше или равно нулю")
		return ""
	}

	distM := float64(steps) * stepLength // дистанция в метрах

	distKm := distM / float64(mInKm) // дистанция в км

	if distKm <= 0 { // если дист в км меньше 0 то ошибка и пустая строка
		fmt.Println("дистанция должна быть больше 0")
		return ""
	}

	// получаю калории из ф-ции из другого пакета
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil || calories <= 0 { // если ошибка или каллорий меньшу или равно нулю то ошибка и пустая строка
		fmt.Println("ошибка ", err)
		return ""
	}

	// результат
	res := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distKm, calories)

	return res
}
