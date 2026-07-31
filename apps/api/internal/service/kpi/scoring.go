package kpi

import (
	"math"
	"sort"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	domainkpi "github.com/gilabs/crm-healthcare/api/internal/domain/kpi"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

func kpiConfig() config.KPIConfig {
	if config.AppConfig != nil {
		return config.AppConfig.KPI
	}
	return config.DefaultKPIConfig()
}

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func pctValue(numerator, denominator int64) *float64 {
	if denominator == 0 {
		return nil
	}
	v := (float64(numerator) / float64(denominator)) * 100
	return &v
}

func normalizePercentage(value *float64) *float64 {
	if value == nil {
		return nil
	}
	v := clampFloat(*value, 0, 100)
	return &v
}

func normalizeInverseRate(value *float64) *float64 {
	if value == nil {
		return nil
	}
	v := clampFloat(100-*value, 0, 100)
	return &v
}

func normalizeAverageDealValue(current, benchmark *float64) *float64 {
	if current == nil || benchmark == nil || *benchmark <= 0 {
		return nil
	}
	ratio := (*current / *benchmark) * 100
	v := clampFloat(ratio, 0, 100)
	return &v
}

func normalizePipelineMovement(score int64, dealsCreated int64) *float64 {
	if dealsCreated == 0 {
		return nil
	}
	benchmark := float64(dealsCreated) * kpiConfig().Normalization.PipelineMovementPerDealFactor
	if benchmark <= 0 {
		return nil
	}
	v := clampFloat((float64(score)/benchmark)*100, 0, 100)
	return &v
}

func gradeForScore(score float64) string {
	for _, band := range kpiConfig().GradeBands {
		if score >= band.Min && score <= band.Max {
			return band.Label
		}
	}
	return kpiConfig().GradeBands[len(kpiConfig().GradeBands)-1].Label
}

func trendFromScores(current, previous *float64) domainkpi.KPITrend {
	if current == nil || previous == nil {
		return domainkpi.KPITrend{Direction: "flat"}
	}
	delta := *current - *previous
	direction := "flat"
	if math.Abs(delta) >= kpiConfig().Normalization.TrendFlatThreshold {
		if delta > 0 {
			direction = "up"
		} else {
			direction = "down"
		}
	}
	return domainkpi.KPITrend{
		PreviousCompositeScore: previous,
		Delta:                  &delta,
		Direction:              direction,
	}
}

func computeWeightedScore(values map[string]*float64, weights map[string]float64) float64 {
	total := 0.0
	for key, weight := range weights {
		value := values[key]
		if value == nil {
			continue
		}
		total += *value * weight
	}
	return total
}

func (s *Service) salesRepWeights() map[string]float64 {
	weights := kpiConfig().SalesRepWeights
	return map[string]float64{
		"conversion_rate":           weights.ConversionRate,
		"revenue_target_attainment": weights.RevenueTargetAttainment,
		"visit_compliance":          weights.VisitCompliance,
		"overdue_task_rate":         weights.OverdueTaskRate,
		"average_deal_value":        weights.AverageDealValue,
		"pipeline_movement_score":   weights.PipelineMovementScore,
	}
}

func (s *Service) salesManagerWeights() map[string]float64 {
	weights := kpiConfig().SalesManagerWeights
	return map[string]float64{
		"team_target_attainment": weights.TeamTargetAttainment,
		"team_conversion_rate":   weights.TeamConversionRate,
		"territory_coverage":     weights.TerritoryCoverage,
		"team_visit_compliance":  weights.TeamVisitCompliance,
		"team_overdue_task_rate":  weights.TeamOverdueTaskRate,
		"brick_pipeline_movement": weights.BrickPipelineMovement,
	}
}

func (s *Service) buildSalesRepEvaluation(current domainkpi.SalesRepScorecard, previous *domainkpi.SalesRepScorecard) (domainkpi.KPIEvaluation, []domainkpi.KPIDiagnostic) {
	conversion := normalizePercentage(current.ConversionRate)
	revenueAttainment := normalizePercentage(current.RevenueTargetAttainment)
	visitCompliance := normalizePercentage(current.VisitCompliance)
	overdueRate := normalizeInverseRate(current.OverdueTaskRate)
	dealValue := (*float64)(nil)
	if current.AverageDealValue != nil {
		dealValue = normalizeAverageDealValue(current.AverageDealValue, current.AverageDealValue)
	}
	pipelineMovement := normalizePipelineMovement(current.PipelineMovementScore, current.DealsCreated)

	values := map[string]*float64{
		"conversion_rate":           conversion,
		"revenue_target_attainment": revenueAttainment,
		"visit_compliance":          visitCompliance,
		"overdue_task_rate":         overdueRate,
		"average_deal_value":        dealValue,
		"pipeline_movement_score":   pipelineMovement,
	}
	currentScore := computeWeightedScore(values, s.salesRepWeights())
	var previousScore *float64
	if previous != nil {
		prevConversion := normalizePercentage(previous.ConversionRate)
		prevRevenue := normalizePercentage(previous.RevenueTargetAttainment)
		prevVisit := normalizePercentage(previous.VisitCompliance)
		prevOverdue := normalizeInverseRate(previous.OverdueTaskRate)
		prevDealValue := (*float64)(nil)
		if previous.AverageDealValue != nil {
			prevDealValue = normalizeAverageDealValue(previous.AverageDealValue, previous.AverageDealValue)
		}
		prevPipeline := normalizePipelineMovement(previous.PipelineMovementScore, previous.DealsCreated)
		prevScore := computeWeightedScore(map[string]*float64{
			"conversion_rate":           prevConversion,
			"revenue_target_attainment": prevRevenue,
			"visit_compliance":          prevVisit,
			"overdue_task_rate":         prevOverdue,
			"average_deal_value":        prevDealValue,
			"pipeline_movement_score":   prevPipeline,
		}, s.salesRepWeights())
		previousScore = &prevScore
	}
	currentScoreValue := currentScore
	trend := trendFromScores(&currentScoreValue, previousScore)
	revenueTarget := int64(0)
	if current.RevenueTargetAttainment != nil && *current.RevenueTargetAttainment > 0 {
		revenueTarget = int64(math.Round(float64(current.TotalRevenue) / (*current.RevenueTargetAttainment / 100)))
	}
	return domainkpi.KPIEvaluation{
		CompositeScore: currentScore,
		Grade:          gradeForScore(currentScore),
		Trend:          trend,
		TargetGap: domainkpi.KPITargetGap{
			Revenue: domainkpi.KPITargetGapItem{Target: revenueTarget, Actual: current.TotalRevenue, GapPercent: current.RevenueTargetAttainment, Status: targetStatus(current.RevenueTargetAttainment)},
			Deals:   domainkpi.KPITargetGapItem{Target: current.DealsCreated, Actual: current.TotalDealsClosed, Status: "met"},
		},
	}, buildSalesDiagnostics(current, nil, currentScore, previousScore, nil)
}

func valueOrZero(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func targetStatus(attainment *float64) string {
	if attainment == nil {
		return "unknown"
	}
	if *attainment > 100 {
		return "above"
	}
	if *attainment >= 100 {
		return "met"
	}
	return "below"
}

func buildSalesDiagnostics(current domainkpi.SalesRepScorecard, teamSummary *domainkpi.SalesManagerTeamSummary, currentScore float64, previousScore *float64, brickCoverage *float64) []domainkpi.KPIDiagnostic {
	cfg := kpiConfig().Diagnostics
	diagnostics := make([]domainkpi.KPIDiagnostic, 0)
	if current.ConversionRate != nil && current.DealsCreated >= cfg.LowConversionMinDealsCreated && *current.ConversionRate < cfg.LowConversionRateThreshold {
		diagnostics = append(diagnostics, domainkpi.KPIDiagnostic{Code: "LOW_CONVERSION", Severity: "warning", Message: "Conversion rate rendah"})
	}
	if current.RevenueTargetAttainment != nil && *current.RevenueTargetAttainment < cfg.TargetUnderperformThreshold {
		diagnostics = append(diagnostics, domainkpi.KPIDiagnostic{Code: "TARGET_UNDERPERFORM", Severity: "critical", Message: "Pencapaian target revenue masih rendah"})
	}
	if current.VisitCompliance != nil && *current.VisitCompliance < cfg.LowVisitCadenceThreshold {
		diagnostics = append(diagnostics, domainkpi.KPIDiagnostic{Code: "LOW_VISIT_CADENCE", Severity: "warning", Message: "Visit compliance di bawah standar"})
	}
	if current.OverdueTaskRate != nil && *current.OverdueTaskRate > cfg.HighOverdueTaskRateThreshold {
		diagnostics = append(diagnostics, domainkpi.KPIDiagnostic{Code: "HIGH_OVERDUE_TASK", Severity: "critical", Message: "Task overdue terlalu tinggi"})
	}
	if current.PipelineMovementScore <= cfg.StagnantPipelineMaxScore && current.DealsCreated > 0 {
		diagnostics = append(diagnostics, domainkpi.KPIDiagnostic{Code: "STAGNANT_PIPELINE", Severity: "warning", Message: "Pipeline tidak bergerak signifikan"})
	}
	if previousScore != nil && teamSummary != nil && teamSummary.TeamTargetAttainment != nil {
		_ = teamSummary
	}
	if brickCoverage != nil && *brickCoverage < 100 {
		_ = brickCoverage
	}
	if len(diagnostics) == 0 {
		return nil
	}
	return diagnostics
}

func (s *Service) buildManagerEvaluation(current float64, previous *float64) domainkpi.SalesManagerEvaluation {
	trend := trendFromScores(&current, previous)
	return domainkpi.SalesManagerEvaluation{
		CompositeScore: current,
		Grade:          gradeForScore(current),
		Trend:          trend,
	}
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func uniqueUsers(users []user.User) []user.User {
	seen := make(map[string]struct{}, len(users))
	result := make([]user.User, 0, len(users))
	for _, u := range users {
		if _, exists := seen[u.ID]; exists {
			continue
		}
		seen[u.ID] = struct{}{}
		result = append(result, u)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}
