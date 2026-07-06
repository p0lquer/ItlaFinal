package predictor

import "math"

// WeightedAverage calcula el tiempo estimado dando más peso
// a los datos más recientes. Los últimos registros pesan más.
func WeightedAverage(historicalMinutes []float64) float64 {
	if len(historicalMinutes) == 0 {
		return 60 // default: 60 minutos
	}

	totalWeight := 0.0
	weightedSum := 0.0
	// n := len(historicalMinutes)

	for i, minutes := range historicalMinutes {
		// El dato más reciente (último índice) tiene peso mayor
		weight := float64(i + 1)
		weightedSum += minutes * weight
		totalWeight += weight
	}

	return math.Round(weightedSum / totalWeight)
}

// LinearRegression predice el tiempo según cantidad de piezas.
// Usa la fórmula: y = a + b*x
// Donde x = piezas, y = tiempo estimado
func LinearRegression(data []DataPoint) func(pieces int) float64 {
	if len(data) < 2 {
		return func(pieces int) float64 { return 60 }
	}

	n := float64(len(data))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for _, p := range data {
		x := float64(p.Pieces)
		y := p.ActualMinutes
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	denominator := n*sumX2 - sumX*sumX
	// valida si el denominador es cero para evitar división por cero
	if denominator == 0 {
		// If all historical jobs had the exact same number of pieces,
		// we can't find a slope. Just return the average time.
		averageTime := sumY / n
		return func(pieces int) float64 { return averageTime }
	}

	// Calcular pendiente (b) e intercepto (a)
	b := (n*sumXY - sumX*sumY) / denominator
	a := (sumY - b*sumX) / n

	return func(pieces int) float64 {
		result := a + b*float64(pieces)
		if result < 20 {
			return 20 // mínimo 20 minutos
		}
		return math.Round(result)
	}
}

// DataPoint representa un registro histórico para regresión
type DataPoint struct {
	Pieces        int
	ActualMinutes float64
}
