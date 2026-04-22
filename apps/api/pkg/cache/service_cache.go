// Package cache provides cache integration for services.
//
// This file contains cache wrappers for service layer operations,
// implementing the cache-aside pattern for read-heavy operations.
package cache

import (
	"time"
)

// ServiceCacheConfig holds configuration for service-level caching
type ServiceCacheConfig struct {
	// Enable/disable caching for specific operations
	EnableListCache   bool
	EnableDetailCache bool
	EnableStatsCache  bool

	// TTL settings
	ListTTL   time.Duration
	DetailTTL time.Duration
	StatsTTL  time.Duration
}

// DefaultServiceCacheConfig returns default cache configuration
func DefaultServiceCacheConfig() *ServiceCacheConfig {
	return &ServiceCacheConfig{
		EnableListCache:   true,
		EnableDetailCache: true,
		EnableStatsCache:  true,
		ListTTL:           TTLListShort,
		DetailTTL:         TTLDetailMedium,
		StatsTTL:          TTLStatsMedium,
	}
}

// DashboardServiceCacheConfig returns cache configuration optimized for dashboard
func DashboardServiceCacheConfig() *ServiceCacheConfig {
	return &ServiceCacheConfig{
		EnableListCache:   true,
		EnableDetailCache: true,
		EnableStatsCache:  true,
		ListTTL:           TTLDashboardShort, // 1 minute for frequently changing data
		DetailTTL:         TTLDashboardLong,  // 5 minutes for overview data
		StatsTTL:          TTLStatsMedium,    // 3 minutes for aggregated stats
	}
}

// CachedDashboardData represents cached dashboard overview data
type CachedDashboardData struct {
	Data     interface{} `msgpack:"data"`
	CachedAt time.Time   `msgpack:"cached_at"`
	UserID   string      `msgpack:"user_id"`
	Period   string      `msgpack:"period"`
}

// CachedListData represents cached list data with pagination
type CachedListData struct {
	Items      interface{} `msgpack:"items"`
	Total      int64       `msgpack:"total"`
	Page       int         `msgpack:"page"`
	PerPage    int         `msgpack:"per_page"`
	TotalPages int         `msgpack:"total_pages"`
	CachedAt   time.Time   `msgpack:"cached_at"`
}

// CachedDetailData represents cached single entity data
type CachedDetailData struct {
	Data     interface{} `msgpack:"data"`
	CachedAt time.Time   `msgpack:"cached_at"`
}

// CachedAggregateData represents cached aggregate/stats data
type CachedAggregateData struct {
	Data     interface{} `msgpack:"data"`
	CachedAt time.Time   `msgpack:"cached_at"`
	Period   string      `msgpack:"period"`
}

// DashboardCacheService provides caching for dashboard operations
type DashboardCacheService struct {
	cache  *DashboardCache
	config *ServiceCacheConfig
}

// NewDashboardCacheService creates a new dashboard cache service
func NewDashboardCacheService(config *ServiceCacheConfig) *DashboardCacheService {
	if config == nil {
		config = DashboardServiceCacheConfig()
	}
	return &DashboardCacheService{
		cache:  NewDashboardCache(),
		config: config,
	}
}

// GetOverview retrieves cached dashboard overview or returns nil if not found
func (s *DashboardCacheService) GetOverview(userID, period, startDate, endDate string, target interface{}) (bool, error) {
	if !s.config.EnableStatsCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.OverviewKey(userID, period, startDate, endDate)
	return s.cache.ec.Get(key, target)
}

// SetOverview caches dashboard overview data
func (s *DashboardCacheService) SetOverview(userID, period, startDate, endDate string, data interface{}) error {
	if !s.config.EnableStatsCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.OverviewKey(userID, period, startDate, endDate)
	return s.cache.ec.Set(key, data, s.config.StatsTTL)
}

// GetVisitStats retrieves cached visit statistics or returns nil if not found
func (s *DashboardCacheService) GetVisitStats(userID, period, startDate, endDate string, target interface{}) (bool, error) {
	if !s.config.EnableStatsCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.VisitStatsKey(userID, period, startDate, endDate)
	return s.cache.ec.Get(key, target)
}

// SetVisitStats caches visit statistics data
func (s *DashboardCacheService) SetVisitStats(userID, period, startDate, endDate string, data interface{}) error {
	if !s.config.EnableStatsCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.VisitStatsKey(userID, period, startDate, endDate)
	return s.cache.ec.Set(key, data, s.config.StatsTTL)
}

// GetSales retrieves cached sales data or returns nil if not found
func (s *DashboardCacheService) GetSales(userID, period, startDate, endDate string, target interface{}) (bool, error) {
	if !s.config.EnableStatsCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.SalesKey(userID, period, startDate, endDate)
	return s.cache.ec.Get(key, target)
}

// SetSales caches sales data
func (s *DashboardCacheService) SetSales(userID, period, startDate, endDate string, data interface{}) error {
	if !s.config.EnableStatsCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.SalesKey(userID, period, startDate, endDate)
	return s.cache.ec.Set(key, data, s.config.StatsTTL)
}

// InvalidateAll invalidates all dashboard caches
func (s *DashboardCacheService) InvalidateAll() error {
	return s.cache.InvalidateAll()
}

// InvalidateUserDashboard invalidates dashboard cache for a specific user
func (s *DashboardCacheService) InvalidateUserDashboard(userID string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	// Delete all dashboard cache keys containing this user ID
	return s.cache.ec.DeletePattern("dashboard:*:user=" + userID + "*")
}

// AccountCacheService provides caching for account operations
type AccountCacheService struct {
	cache  *AccountCache
	config *ServiceCacheConfig
}

// NewAccountCacheService creates a new account cache service
func NewAccountCacheService(config *ServiceCacheConfig) *AccountCacheService {
	if config == nil {
		config = DefaultServiceCacheConfig()
	}
	return &AccountCacheService{
		cache:  NewAccountCache(),
		config: config,
	}
}

// GetList retrieves cached account list or returns nil if not found
func (s *AccountCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches account list data
func (s *AccountCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached account detail or returns nil if not found
func (s *AccountCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches account detail data
func (s *AccountCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// InvalidateOnWrite invalidates account caches after write operations
func (s *AccountCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	// Invalidate list caches (all pagination variants)
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	// Invalidate specific detail cache
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
	}
	// Invalidate stats caches
	if err := s.cache.InvalidateStats(); err != nil {
		return err
	}
	// Invalidate related dashboard caches
	return InvalidateRelatedEntities("dashboard")
}

// LeadCacheService provides caching for lead operations
type LeadCacheService struct {
	cache  *LeadCache
	config *ServiceCacheConfig
}

// NewLeadCacheService creates a new lead cache service
func NewLeadCacheService(config *ServiceCacheConfig) *LeadCacheService {
	if config == nil {
		config = DefaultServiceCacheConfig()
	}
	return &LeadCacheService{
		cache:  NewLeadCache(),
		config: config,
	}
}

// GetList retrieves cached lead list or returns nil if not found
func (s *LeadCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches lead list data
func (s *LeadCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached lead detail or returns nil if not found
func (s *LeadCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches lead detail data
func (s *LeadCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// InvalidateOnWrite invalidates lead caches after write operations
func (s *LeadCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
	}
	if err := s.cache.InvalidateStats(); err != nil {
		return err
	}
	return InvalidateRelatedEntities("dashboard")
}

// ContactCacheService provides caching for contact operations
type ContactCacheService struct {
	cache  *ContactCache
	config *ServiceCacheConfig
}

// NewContactCacheService creates a new contact cache service
func NewContactCacheService(config *ServiceCacheConfig) *ContactCacheService {
	if config == nil {
		config = DefaultServiceCacheConfig()
	}
	return &ContactCacheService{
		cache:  NewContactCache(),
		config: config,
	}
}

// GetList retrieves cached contact list or returns nil if not found
func (s *ContactCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches contact list data
func (s *ContactCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached contact detail or returns nil if not found
func (s *ContactCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches contact detail data
func (s *ContactCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// InvalidateOnWrite invalidates contact caches after write operations
func (s *ContactCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
	}
	return nil
}

// ActivityCacheService provides caching for activity operations
type ActivityCacheService struct {
	cache  *ActivityCache
	config *ServiceCacheConfig
}

// NewActivityCacheService creates a new activity cache service
func NewActivityCacheService(config *ServiceCacheConfig) *ActivityCacheService {
	if config == nil {
		config = DefaultServiceCacheConfig()
	}
	return &ActivityCacheService{
		cache:  NewActivityCache(),
		config: config,
	}
}

// GetList retrieves cached activity list or returns nil if not found
func (s *ActivityCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches activity list data
func (s *ActivityCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached activity detail or returns nil if not found
func (s *ActivityCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches activity detail data
func (s *ActivityCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// InvalidateOnWrite invalidates activity caches after write operations
func (s *ActivityCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
	}
	if err := s.cache.InvalidateStats(); err != nil {
		return err
	}
	return InvalidateRelatedEntities("dashboard")
}

// TaskCacheService provides caching for task operations
type TaskCacheService struct {
	cache  *TaskCache
	config *ServiceCacheConfig
}

// NewTaskCacheService creates a new task cache service
func NewTaskCacheService(config *ServiceCacheConfig) *TaskCacheService {
	if config == nil {
		config = DefaultServiceCacheConfig()
	}
	return &TaskCacheService{
		cache:  NewTaskCache(),
		config: config,
	}
}

// GetList retrieves cached task list or returns nil if not found
func (s *TaskCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches task list data
func (s *TaskCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached task detail or returns nil if not found
func (s *TaskCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches task detail data
func (s *TaskCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// InvalidateOnWrite invalidates task caches after write operations
func (s *TaskCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
	}
	return InvalidateRelatedEntities("dashboard")
}

// VisitReportCacheService provides caching for visit report operations
type VisitReportCacheService struct {
	cache  *VisitReportCache
	config *ServiceCacheConfig
}

// NewVisitReportCacheService creates a new visit report cache service
func NewVisitReportCacheService(config *ServiceCacheConfig) *VisitReportCacheService {
	if config == nil {
		config = DefaultServiceCacheConfig()
	}
	return &VisitReportCacheService{
		cache:  NewVisitReportCache(),
		config: config,
	}
}

// GetList retrieves cached visit report list or returns nil if not found
func (s *VisitReportCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches visit report list data
func (s *VisitReportCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached visit report detail or returns nil if not found
func (s *VisitReportCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches visit report detail data
func (s *VisitReportCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// InvalidateOnWrite invalidates visit report caches after write operations
func (s *VisitReportCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
	}
	if err := s.cache.InvalidateStats(); err != nil {
		return err
	}
	return InvalidateRelatedEntities("dashboard")
}

// BrickCacheService provides caching for brick operations
type BrickCacheService struct {
	cache  *BrickCache
	config *ServiceCacheConfig
}

// NewBrickCacheService creates a new brick cache service
func NewBrickCacheService(config *ServiceCacheConfig) *BrickCacheService {
	if config == nil {
		config = &ServiceCacheConfig{
			EnableListCache:   true,
			EnableDetailCache: true,
			EnableStatsCache:  true,
			ListTTL:           TTLReferenceData, // Brick data rarely changes
			DetailTTL:         TTLReferenceData,
			StatsTTL:          TTLStatsLong,
		}
	}
	return &BrickCacheService{
		cache:  NewBrickCache(),
		config: config,
	}
}

// GetList retrieves cached brick list or returns nil if not found
func (s *BrickCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches brick list data
func (s *BrickCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached brick detail or returns nil if not found
func (s *BrickCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches brick detail data
func (s *BrickCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// GetGeoJSON retrieves cached brick GeoJSON or returns nil if not found
func (s *BrickCacheService) GetGeoJSON(brickID string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.GeoJSONKey(brickID)
	return s.cache.ec.Get(key, target)
}

// SetGeoJSON caches brick GeoJSON data
func (s *BrickCacheService) SetGeoJSON(brickID string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.GeoJSONKey(brickID)
	return s.cache.ec.Set(key, data, TTLGeoJSON)
}

// InvalidateOnWrite invalidates brick caches after write operations
func (s *BrickCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
		// Also invalidate GeoJSON cache
		key := s.cache.GeoJSONKey(id)
		if err := s.cache.ec.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// ScheduleCacheService provides caching for schedule operations
type ScheduleCacheService struct {
	cache  *ScheduleCache
	config *ServiceCacheConfig
}

// NewScheduleCacheService creates a new schedule cache service
func NewScheduleCacheService(config *ServiceCacheConfig) *ScheduleCacheService {
	if config == nil {
		config = DefaultServiceCacheConfig()
	}
	return &ScheduleCacheService{
		cache:  NewScheduleCache(),
		config: config,
	}
}

// GetList retrieves cached schedule list or returns nil if not found
func (s *ScheduleCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches schedule list data
func (s *ScheduleCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached schedule detail or returns nil if not found
func (s *ScheduleCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches schedule detail data
func (s *ScheduleCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// InvalidateOnWrite invalidates schedule caches after write operations
func (s *ScheduleCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
	}
	return nil
}

// DealCacheService provides caching for deal operations
type DealCacheService struct {
	cache  *DealCache
	config *ServiceCacheConfig
}

// NewDealCacheService creates a new deal cache service
func NewDealCacheService(config *ServiceCacheConfig) *DealCacheService {
	if config == nil {
		config = DefaultServiceCacheConfig()
	}
	return &DealCacheService{
		cache:  NewDealCache(),
		config: config,
	}
}

// GetList retrieves cached deal list or returns nil if not found
func (s *DealCacheService) GetList(page, perPage int, filters map[string]interface{}, target interface{}) (bool, error) {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.GetList(key, target)
}

// SetList caches deal list data
func (s *DealCacheService) SetList(page, perPage int, filters map[string]interface{}, data interface{}) error {
	if !s.config.EnableListCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.ListKey(page, perPage, filters)
	return s.cache.SetList(key, data, s.config.ListTTL)
}

// GetDetail retrieves cached deal detail or returns nil if not found
func (s *DealCacheService) GetDetail(id string, target interface{}) (bool, error) {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return false, nil
	}
	return s.cache.GetDetail(id, target)
}

// SetDetail caches deal detail data
func (s *DealCacheService) SetDetail(id string, data interface{}) error {
	if !s.config.EnableDetailCache || !s.cache.IsEnabled() {
		return nil
	}
	return s.cache.SetDetail(id, data, s.config.DetailTTL)
}

// GetSummary retrieves cached deal summary or returns nil if not found
func (s *DealCacheService) GetSummary(userID string, target interface{}) (bool, error) {
	if !s.config.EnableStatsCache || !s.cache.IsEnabled() {
		return false, nil
	}
	key := s.cache.SummaryKey(userID)
	return s.cache.ec.Get(key, target)
}

// SetSummary caches deal summary data
func (s *DealCacheService) SetSummary(userID string, data interface{}) error {
	if !s.config.EnableStatsCache || !s.cache.IsEnabled() {
		return nil
	}
	key := s.cache.SummaryKey(userID)
	return s.cache.ec.Set(key, data, s.config.StatsTTL)
}

// InvalidateOnWrite invalidates deal caches after write operations
func (s *DealCacheService) InvalidateOnWrite(id string) error {
	if !s.cache.IsEnabled() {
		return nil
	}
	if err := s.cache.InvalidateList(); err != nil {
		return err
	}
	if id != "" {
		if err := s.cache.InvalidateDetail(id); err != nil {
			return err
		}
	}
	// Invalidate all summaries
	if err := s.cache.ec.DeletePattern(PrefixDealSummary + "*"); err != nil {
		return err
	}
	return InvalidateRelatedEntities("dashboard")
}
