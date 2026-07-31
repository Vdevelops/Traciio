// Auto-generated types for KPI feature — follow BE contract exactly

export type KPIGrade = "Excellent" | "Good" | "Needs Improvement" | "Critical";
export type TrendDirection = "up" | "down" | "flat";
export type DiagnosticSeverity = "info" | "warning" | "critical";

export interface DiagnosticFlag {
  code: string;
  severity: DiagnosticSeverity;
  message: string;
  brickId?: string;
}

export interface TrendInfo {
  previousCompositeScore: number;
  delta: number;
  direction: TrendDirection;
}

export interface KPIMeta {
  brickMissingCount: number;
  brickInferredCount: number;
  generatedAt: string;
}

export interface SalesRepScorecard {
  totalDealsClosed: number;
  totalRevenue: number;
  dealsCreated: number;
  conversionRate: number | null;
  averageDealValue: number | null;
  visitCompleted: number;
  visitPlanned: number;
  visitCompliance: number | null;
  tasksCompleted: number;
  overdueTaskRate: number | null;
  revenueTargetAttainment: number | null;
  dealTargetAttainment: number | null;
  pipelineMovementScore: number;
}

export interface TargetGapItem {
  target: number;
  actual: number;
  gapPercent: number;
  status: "above" | "met" | "below";
}

export interface SalesRepEvaluation {
  compositeScore: number;
  grade: KPIGrade;
  trend: TrendInfo;
  targetGap: {
    revenue: TargetGapItem;
    deals: TargetGapItem;
  };
}

export interface SalesRepKPIResponse {
  scope: { userId: string; role: "sales_rep"; startDate: string; endDate: string };
  scorecard: SalesRepScorecard;
  evaluation: SalesRepEvaluation;
  diagnostics: DiagnosticFlag[];
  meta: KPIMeta;
}

export interface TeamSummary {
  totalRepsCount: number;
  totalDealsClosed: number;
  totalRevenue: number;
  teamConversionRate: number | null;
  teamVisitCompliance: number | null;
  teamOverdueTaskRate: number | null;
  teamTargetAttainment: number | null;
}

export interface TeamBreakdownItem {
  userId: string;
  name: string;
  compositeScore: number;
  grade: KPIGrade;
  totalRevenue: number;
  conversionRate: number | null;
  rank: number;
}

export interface BrickBreakdownItem {
  brickId: string;
  name: string;
  coveragePenetration: number | null;
  totalRevenue: number;
  repsCount: number;
  compositeScore: number;
}

export interface SalesManagerKPIResponse {
  scope: { managerId: string; role: "sales_manager"; startDate: string; endDate: string; bricks: string[] };
  teamSummary: TeamSummary;
  evaluation: { compositeScore: number; grade: KPIGrade; trend: TrendInfo };
  teamBreakdown: TeamBreakdownItem[];
  brickBreakdown: BrickBreakdownItem[];
  topBottomPerformers: { top: string[]; bottom: string[] };
  diagnostics: DiagnosticFlag[];
  meta: KPIMeta;
}

export interface KPIDateRangeParams {
  startDate: string;
  endDate: string;
  compareWithPrevious?: boolean;
}

export interface SalesRepKPIParams extends KPIDateRangeParams {
  userId?: string;
}

export interface SalesManagerKPIParams extends KPIDateRangeParams {
  managerId?: string;
  brickId?: string;
  includeTeamBreakdown?: boolean;
}
