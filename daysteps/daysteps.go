package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	stepLength    = 0.65 // метра
	mInKm         = 1000.0
)
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: ожидается 2 части, получено %d", len(parts))
	}
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга шагов: %v", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным, получено: %d", steps)
	}
	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка парсинга данных: %v", err)
		return ""
	}

	if steps <= 0 {
		log.Printf("Некорректное количество шагов: %d", steps)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm
	calories, err := WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка вычисления калорий: %v", err)
		return ""
	}
	result := fmt.Sprintf("Количество шагов: %d.\n", steps)
	result += fmt.Sprintf("Дистанция составила %.2f км.\n", distanceKm)
	result += fmt.Sprintf("Вы сожгли %.2f ккал.", calories)

	return result
}