package events

import (
	"log"

	pkgevents "github.com/gilabs/crm-healthcare/api/pkg/events"
)

// Helper is a helper for emitting domain events
type Helper struct {
	producer pkgevents.EventProducer
}

// NewHelper creates a new event helper
func NewHelper(producer pkgevents.EventProducer) *Helper {
	return &Helper{
		producer: producer,
	}
}

// EmitLeadCreated emits a LeadCreatedEvent
func (h *Helper) EmitLeadCreated(event *LeadCreatedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("lead.created", event.LeadID, "lead", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create LeadCreatedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitLeadConverted emits a LeadConvertedEvent
func (h *Helper) EmitLeadConverted(event *LeadConvertedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("lead.converted", event.LeadID, "lead", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create LeadConvertedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitLeadStatusChanged emits a LeadStatusChangedEvent
func (h *Helper) EmitLeadStatusChanged(event *LeadStatusChangedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("lead.status_changed", event.LeadID, "lead", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create LeadStatusChangedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitDealCreated emits a DealCreatedEvent
func (h *Helper) EmitDealCreated(event *DealCreatedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("deal.created", event.DealID, "deal", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create DealCreatedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitDealStageChanged emits a DealStageChangedEvent
func (h *Helper) EmitDealStageChanged(event *DealStageChangedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("deal.stage_changed", event.DealID, "deal", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create DealStageChangedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitDealWon emits a DealWonEvent
func (h *Helper) EmitDealWon(event *DealWonEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("deal.won", event.DealID, "deal", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create DealWonEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitDealLost emits a DealLostEvent
func (h *Helper) EmitDealLost(event *DealLostEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("deal.lost", event.DealID, "deal", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create DealLostEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitActivityLogged emits an ActivityLoggedEvent
func (h *Helper) EmitActivityLogged(event *ActivityLoggedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("activity.logged", event.ActivityID, "activity", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create ActivityLoggedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitVisitCompleted emits a VisitCompletedEvent
func (h *Helper) EmitVisitCompleted(event *VisitCompletedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("visit.completed", event.VisitReportID, "visit_report", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create VisitCompletedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitNotificationCreated emits a NotificationCreatedEvent
func (h *Helper) EmitNotificationCreated(event *NotificationCreatedEvent) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": event.UserID,
	}

	e, err := pkgevents.NewEvent("notification.created", event.NotificationID, "notification", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create NotificationCreatedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitTaskCreated emits a TaskCreatedEvent
func (h *Helper) EmitTaskCreated(event *TaskCreatedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("task.created", event.TaskID, "task", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create TaskCreatedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitTaskAssigned emits a TaskAssignedEvent
func (h *Helper) EmitTaskAssigned(event *TaskAssignedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("task.assigned", event.TaskID, "task", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create TaskAssignedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}

// EmitTaskCompleted emits a TaskCompletedEvent
func (h *Helper) EmitTaskCompleted(event *TaskCompletedEvent, userID string) {
	if h.producer == nil {
		return
	}

	metadata := map[string]string{
		"user_id": userID,
	}

	e, err := pkgevents.NewEvent("task.completed", event.TaskID, "task", event, metadata)
	if err != nil {
		log.Printf("ERROR: Failed to create TaskCompletedEvent: %v", err)
		return
	}

	h.producer.PublishAsync(e)
}
