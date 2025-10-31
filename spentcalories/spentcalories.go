package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	stepLengthCoefficient       = 0.414
	mInKm                       = 1000.0
	minInH                      = 60.0
	walkingCaloriesCoefficient  = 0.789
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 3 части, получено %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга шагов: %v", err)
	}

	trainingType := strings.TrimSpace(parts[1])

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	
	dist := distance(steps, height)
	hours := duration.Hours()
	
	if hours == 0 {
		return 0
	}
	
	return dist / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным, получено: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным, получено: %.2f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным, получено: %.2f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной, получено: %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным, получено: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным, получено: %.2f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным, получено: %.2f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной, получено: %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient
	
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Printf("Ошибка парсинга тренировки: %v", err)
		return "", err
	}

	var calories float64
	var speed float64
	var dist float64

	switch strings.ToLower(trainingType) {
	case "бег", "running":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		speed = meanSpeed(steps, height, duration)
		dist = distance(steps, height)
		
	case "ходьба", "walking":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		speed = meanSpeed(steps, height, duration)
		dist = distance(steps, height)
		
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", trainingType)
	}

	result := fmt.Sprintf("Тип тренировки: %s\n", trainingType)
	result += fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
	result += fmt.Sprintf("Дистанция: %.2f км.\n", dist)
	result += fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
	result += fmt.Sprintf("Сожгли калорий: %.2f", calories)

	return result, nil
}