package config

import (
	"fmt"
	"os"
)

type KPIConfig struct {
	SalesRepWeights     KPIWeightConfig
	SalesManagerWeights KPIWeightConfig
	Normalization       KPINormalizationConfig
	GradeBands          []KPIGradeBand
	Diagnostics         KPIDiagnosticsConfig
}

type KPIWeightConfig struct {
	ConversionRate          float64
	RevenueTargetAttainment float64
	VisitCompliance         float64
	OverdueTaskRate         float64
	AverageDealValue        float64
	PipelineMovementScore    float64
	TeamTargetAttainment    float64
	TeamConversionRate      float64
	TerritoryCoverage       float64
	TeamVisitCompliance     float64
	TeamOverdueTaskRate     float64
	BrickPipelineMovement   float64
}

type KPIGradeBand struct {
	Min   float64
	Max   float64
	Label string
	Color string
}

type KPINormalizationConfig struct {
	ConversionRateTarget         float64
	PipelineMovementPerDealFactor float64
	TrendFlatThreshold           float64
}

type KPIDiagnosticsConfig struct {
	LowConversionRateThreshold      float64
	LowConversionMinDealsCreated     int64
	TargetUnderperformThreshold     float64
	LowVisitCadenceThreshold        float64
	HighOverdueTaskRateThreshold    float64
	CoverageDeclineThreshold        float64
	DataQualityBrickMissingThreshold int64
	StagnantPipelineMaxScore         int64
}

func DefaultKPIConfig() KPIConfig {
	return KPIConfig{
		SalesRepWeights: KPIWeightConfig{
			ConversionRate:          0.25,
			RevenueTargetAttainment: 0.25,
			VisitCompliance:         0.20,
			OverdueTaskRate:         0.10,
			AverageDealValue:        0.10,
			PipelineMovementScore:    0.10,
		},
		SalesManagerWeights: KPIWeightConfig{
			TeamTargetAttainment:   0.30,
			TeamConversionRate:     0.20,
			TerritoryCoverage:      0.20,
			TeamVisitCompliance:    0.15,
			TeamOverdueTaskRate:    0.10,
			BrickPipelineMovement: 0.05,
		},
		Normalization: KPINormalizationConfig{
			ConversionRateTarget:         25,
			PipelineMovementPerDealFactor: 4,
			TrendFlatThreshold:           1,
		},
		GradeBands: []KPIGradeBand{
			{Min: 85, Max: 100, Label: "Excellent", Color: "green"},
			{Min: 70, Max: 84, Label: "Good", Color: "blue"},
			{Min: 55, Max: 69, Label: "Needs Improvement", Color: "yellow"},
			{Min: 0, Max: 54, Label: "Critical", Color: "red"},
		},
		Diagnostics: KPIDiagnosticsConfig{
			LowConversionRateThreshold:       15,
			LowConversionMinDealsCreated:     10,
			TargetUnderperformThreshold:      70,
			LowVisitCadenceThreshold:         60,
			HighOverdueTaskRateThreshold:     25,
			CoverageDeclineThreshold:         5,
			DataQualityBrickMissingThreshold: 1,
			StagnantPipelineMaxScore:         0,
		},
	}
}

func loadKPIConfig() KPIConfig {
	def := DefaultKPIConfig()
	return KPIConfig{
		SalesRepWeights: KPIWeightConfig{
			ConversionRate:          getEnvAsFloat("KPI_SALES_REP_CONVERSION_WEIGHT", def.SalesRepWeights.ConversionRate),
			RevenueTargetAttainment: getEnvAsFloat("KPI_SALES_REP_REVENUE_TARGET_WEIGHT", def.SalesRepWeights.RevenueTargetAttainment),
			VisitCompliance:         getEnvAsFloat("KPI_SALES_REP_VISIT_COMPLIANCE_WEIGHT", def.SalesRepWeights.VisitCompliance),
			OverdueTaskRate:         getEnvAsFloat("KPI_SALES_REP_OVERDUE_TASK_WEIGHT", def.SalesRepWeights.OverdueTaskRate),
			AverageDealValue:        getEnvAsFloat("KPI_SALES_REP_AVG_DEAL_VALUE_WEIGHT", def.SalesRepWeights.AverageDealValue),
			PipelineMovementScore:    getEnvAsFloat("KPI_SALES_REP_PIPELINE_MOVEMENT_WEIGHT", def.SalesRepWeights.PipelineMovementScore),
		},
		SalesManagerWeights: KPIWeightConfig{
			TeamTargetAttainment:   getEnvAsFloat("KPI_SALES_MANAGER_TEAM_TARGET_WEIGHT", def.SalesManagerWeights.TeamTargetAttainment),
			TeamConversionRate:     getEnvAsFloat("KPI_SALES_MANAGER_TEAM_CONVERSION_WEIGHT", def.SalesManagerWeights.TeamConversionRate),
			TerritoryCoverage:      getEnvAsFloat("KPI_SALES_MANAGER_TERRITORY_COVERAGE_WEIGHT", def.SalesManagerWeights.TerritoryCoverage),
			TeamVisitCompliance:    getEnvAsFloat("KPI_SALES_MANAGER_VISIT_COMPLIANCE_WEIGHT", def.SalesManagerWeights.TeamVisitCompliance),
			TeamOverdueTaskRate:    getEnvAsFloat("KPI_SALES_MANAGER_OVERDUE_TASK_WEIGHT", def.SalesManagerWeights.TeamOverdueTaskRate),
			BrickPipelineMovement: getEnvAsFloat("KPI_SALES_MANAGER_PIPELINE_MOVEMENT_WEIGHT", def.SalesManagerWeights.BrickPipelineMovement),
		},
		Normalization: KPINormalizationConfig{
			ConversionRateTarget:         getEnvAsFloat("KPI_NORMALIZE_CONVERSION_TARGET", def.Normalization.ConversionRateTarget),
			PipelineMovementPerDealFactor: getEnvAsFloat("KPI_NORMALIZE_PIPELINE_MOVE_PER_DEAL", def.Normalization.PipelineMovementPerDealFactor),
			TrendFlatThreshold:           getEnvAsFloat("KPI_NORMALIZE_TREND_FLAT_THRESHOLD", def.Normalization.TrendFlatThreshold),
		},
		GradeBands: []KPIGradeBand{
			{
				Min:   getEnvAsFloat("KPI_GRADE_EXCELLENT_MIN", def.GradeBands[0].Min),
				Max:   getEnvAsFloat("KPI_GRADE_EXCELLENT_MAX", def.GradeBands[0].Max),
				Label: getEnv("KPI_GRADE_EXCELLENT_LABEL", def.GradeBands[0].Label),
				Color: getEnv("KPI_GRADE_EXCELLENT_COLOR", def.GradeBands[0].Color),
			},
			{
				Min:   getEnvAsFloat("KPI_GRADE_GOOD_MIN", def.GradeBands[1].Min),
				Max:   getEnvAsFloat("KPI_GRADE_GOOD_MAX", def.GradeBands[1].Max),
				Label: getEnv("KPI_GRADE_GOOD_LABEL", def.GradeBands[1].Label),
				Color: getEnv("KPI_GRADE_GOOD_COLOR", def.GradeBands[1].Color),
			},
			{
				Min:   getEnvAsFloat("KPI_GRADE_NEEDS_IMPROVEMENT_MIN", def.GradeBands[2].Min),
				Max:   getEnvAsFloat("KPI_GRADE_NEEDS_IMPROVEMENT_MAX", def.GradeBands[2].Max),
				Label: getEnv("KPI_GRADE_NEEDS_IMPROVEMENT_LABEL", def.GradeBands[2].Label),
				Color: getEnv("KPI_GRADE_NEEDS_IMPROVEMENT_COLOR", def.GradeBands[2].Color),
			},
			{
				Min:   getEnvAsFloat("KPI_GRADE_CRITICAL_MIN", def.GradeBands[3].Min),
				Max:   getEnvAsFloat("KPI_GRADE_CRITICAL_MAX", def.GradeBands[3].Max),
				Label: getEnv("KPI_GRADE_CRITICAL_LABEL", def.GradeBands[3].Label),
				Color: getEnv("KPI_GRADE_CRITICAL_COLOR", def.GradeBands[3].Color),
			},
		},
		Diagnostics: KPIDiagnosticsConfig{
			LowConversionRateThreshold:       getEnvAsFloat("KPI_DIAG_LOW_CONVERSION_RATE", def.Diagnostics.LowConversionRateThreshold),
			LowConversionMinDealsCreated:     int64(getEnvAsInt("KPI_DIAG_LOW_CONVERSION_MIN_DEALS", int(def.Diagnostics.LowConversionMinDealsCreated))),
			TargetUnderperformThreshold:       getEnvAsFloat("KPI_DIAG_TARGET_UNDERPERFORM", def.Diagnostics.TargetUnderperformThreshold),
			LowVisitCadenceThreshold:         getEnvAsFloat("KPI_DIAG_LOW_VISIT_CADENCE", def.Diagnostics.LowVisitCadenceThreshold),
			HighOverdueTaskRateThreshold:     getEnvAsFloat("KPI_DIAG_HIGH_OVERDUE_TASK_RATE", def.Diagnostics.HighOverdueTaskRateThreshold),
			CoverageDeclineThreshold:         getEnvAsFloat("KPI_DIAG_COVERAGE_DECLINE", def.Diagnostics.CoverageDeclineThreshold),
			DataQualityBrickMissingThreshold: int64(getEnvAsInt("KPI_DIAG_BRICK_MISSING_MIN", int(def.Diagnostics.DataQualityBrickMissingThreshold))),
			StagnantPipelineMaxScore:         int64(getEnvAsInt("KPI_DIAG_STAGNANT_PIPELINE_MAX_SCORE", int(def.Diagnostics.StagnantPipelineMaxScore))),
		},
	}
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	var value float64
	if _, err := fmt.Sscanf(valueStr, "%f", &value); err != nil {
		return defaultValue
	}
	return value
}
