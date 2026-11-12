package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	//lenStep                    = 0.65 // средняя длина шага, хз где ее юзать
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// не экспортируемая, пример получаемой строки: "3456,Ходьба,3h00m"
// парсинг строки и перевод данных из нее в считаемные значения
func parseTraining(data string) (int, string, time.Duration, error) {
	values := strings.Split(data, ",") // разбил строку на слайс строк по запятой

	if len(values) != 3 { // если длинна слайса меньше 3 то нули, ошибка и пустая строка
		return 0, "", 0, errors.New("длинна слайса меньше 3")
	}

	counter, err := strconv.Atoi(values[0]) //конверчу первую строку из слайса
	if err != nil || counter <= 0 {         // если ошибка или ко-во шагов меньше или равно нулю то нули, ошибка и пустая строка
		return 0, "", 0, errors.New("ошибка конвертации")
	}

	duration, err := time.ParseDuration(values[2]) // конверчу строку третью в формат time
	if err != nil || duration <= 0 {               // если есть ошибка или время меньше или равно нулю то нули, ошибка и пустая строка
		return 0, "", 0, errors.New("ошибка парсинга time")
	}
	return counter, values[1], duration, nil
}

// не экспортируемая
// по росту и кол-ву шагов считаю пройденнуб дистанцию в км
func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 { // если шагов или рост меньше или равен нулю то верну 0
		return 0
	}

	lengthStep := height * stepLengthCoefficient // считаю длинну шага
	return (lengthStep * float64(steps)) / mInKm // возвращаю пройденную дистанцибю в км
}

// не экспортируемая
// по росту, кол-ву шагов и времени ходьбы считаю среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 { //если шаги\рост\время меньше или равно нулю то перну ноль
		return 0
	}

	dist := distance(steps, height) // считаю дистанцию
	return dist / duration.Hours()  // возвращаю скорость
}

// по данным получаю стоку с инфой в виде строки
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeLoad, timeLoad, err := parseTraining(data) // получаю шаги, тип нагрузки и время нагрузки
	if err != nil {                                       // если ошибка то пустая строка и ошибка
		log.Println(err)
		return "", err
	}

	var calories float64 // объявил переменную калории

	switch { // перебор типов нагрузки и подсчет по типу нагрузки калории
	case typeLoad == "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, timeLoad)
	case typeLoad == "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, timeLoad)
	default:
		calories = 0 //если непонятный тип нагрузки - то калории будут 0 и ошибка
		err = errors.New("неизвестный тип тренировки")
	}

	if err != nil { // если ошибка то пустая строка и ошибка
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height) // получаю дистанцию

	speed := meanSpeed(steps, height, timeLoad) // получаю скорость
	timeLoadHours := timeLoad.Hours()           // время нагрузки

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeLoad, timeLoadHours, dist, speed, calories)
	return result, nil
}

// по росту/весу, кол-ву шагов и времени бега считаю сколько клорий сжег при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 { // если где-то ноль то ошибка
		return 0, errors.New("кривые параметры")
	}

	speed := meanSpeed(steps, height, duration)                // получаю скорость
	calories := (weight * speed * duration.Minutes()) / minInH //получаю калории для бега

	return calories, nil
}

// по росту/весу, кол-ву шагов и времени бега считаю сколько клорий сжег при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 { //если нули то ошибка
		return 0, errors.New("кривые параметры")
	}

	speed := meanSpeed(steps, height, duration)                                             // получаю скорость
	calories := (weight * speed * duration.Minutes()) / minInH * walkingCaloriesCoefficient // получаю калории для шага

	return calories, nil
}
