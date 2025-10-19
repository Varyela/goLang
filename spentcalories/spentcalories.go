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

// parseTraining парсит строку с данными о тренировке
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 3 части, получено %d", len(parts))
	}

	// Парсим количество шагов
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга шагов: %v", err)
	}

	// Получаем тип тренировки
	trainingType := strings.TrimSpace(parts[1])

	// Парсим продолжительность
	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	return steps, trainingType, duration, nil
}

// distance вычисляет дистанцию в километрах
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость
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

// RunningSpentCalories вычисляет потраченные калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность входных параметров
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

// WalkingSpentCalories вычисляет потраченные калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность входных параметров
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

// TrainingInfo возвращает информацию о тренировке
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

	// Форматируем результат
	result := fmt.Sprintf("Тип тренировки: %s\n", trainingType)
	result += fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
	result += fmt.Sprintf("Дистанция: %.2f км.\n", dist)
	result += fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
	result += fmt.Sprintf("Сожгли калорий: %.2f", calories)

	return result, nil
}