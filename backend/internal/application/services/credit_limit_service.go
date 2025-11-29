package services

import (
	"fmt"
	"math"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
)

type CreditLimitCalculator struct{}

func NewCreditLimitCalculator() *CreditLimitCalculator {
	return &CreditLimitCalculator{}
}

func (c *CreditLimitCalculator) Calculate(input dto.CreditLimitInput) dto.CreditLimitResult {
	I := input.PredictedIncome
	C := input.ActiveCCMaxLimit
	O := input.OutstandSum
	ovrd := input.OverdueSum
	black := input.BlacklistFlag
	avgMonthlyPayment := input.TurnCurrentCreditAvgV2

	PExist := math.Max(O/36.0, 0.0)

	PDNReg := 0.50
	PDNBank := 0.40

	var denom float64
	if C > 0 {
		denom = avgMonthlyPayment / C
		if math.IsInf(denom, 0) || math.IsNaN(denom) {
			denom = 0.15
		}
	} else {
		denom = 0.15
	}

	MLegal := math.Max(PDNReg*I-PExist, 0.0)
	MBank := math.Max(PDNBank*I-PExist, 0.0)

	addLegal := math.Max(MLegal/denom, 0.0)
	addBank := math.Max(MBank/denom, 0.0)

	LLegal := C + addLegal
	LBank := C + addBank

	if black == 1 || ovrd > 0 {
		LLegal = 0.0
		LBank = 0.0
	}

	return dto.CreditLimitResult{
		LimitLegal:                LLegal,
		RecommendationCreditLimit: LBank,
	}
}

func formatFactors(factors map[string]float64) []string {
	result := make([]string, 0, len(factors))
	for name, value := range factors {
		result = append(result, fmt.Sprintf("%s: %.2f", name, value))
	}
	return result
}

func FormatPositiveFactors(factors map[string]float64) []string {
	return formatFactors(factors)
}

func FormatNegativeFactors(factors map[string]float64) []string {
	return formatFactors(factors)
}
