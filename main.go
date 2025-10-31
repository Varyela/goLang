package main

import (
	"fmt"
	"fitness-tracker/daysteps"
	"fitness-tracker/spentcalories"
)

func main() {
	dayData := "678,0h50m"
	weight := 70.0  // кг
	height := 1.75  // м
	
	dayInfo := daysteps.DayActionInfo(dayData, weight, height)
	fmt.Println("=== Дневная активность ===")
	fmt.Println(dayInfo)
	fmt.Println()

	trainingData := "3456,Бег,1h30m"
	trainingInfo, err := spentcalories.TrainingInfo(trainingData, weight, height)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("=== Информация о тренировке ===")
		fmt.Println(trainingInfo)
	}
	fmt.Println()

	trainingData2 := "5000,Ходьба,2h15m"
	trainingInfo2, err := spentcalories.TrainingInfo(trainingData2, weight, height)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("=== Информация о тренировке ===")
		fmt.Println(trainingInfo2)
	}
}