package predictor

import (
    "testing"
)

func TestWeightedAverage(t *testing.T) {
    // Sin datos → debe retornar 60 por defecto
    result := WeightedAverage([]float64{})
    if result != 60 {
        t.Errorf("esperado 60, obtenido %.0f", result)
    }

    // Con datos — los más recientes pesan más
    data := []float64{30, 45, 60, 90}
    result = WeightedAverage(data)
    if result < 60 {
        t.Errorf("el promedio ponderado debería ser >= 60, obtenido %.0f", result)
    }
}

func TestLinearRegression(t *testing.T) {
    data := []DataPoint{
        {Pieces: 5,  ActualMinutes: 30},
        {Pieces: 10, ActualMinutes: 55},
        {Pieces: 15, ActualMinutes: 80},
        {Pieces: 20, ActualMinutes: 110},
    }

    predict := LinearRegression(data)

    // 10 piezas debería dar ~55 minutos
    result := predict(10)
    if result < 40 || result > 70 {
        t.Errorf("predicción para 10 piezas fuera de rango: %.0f", result)
    }
}