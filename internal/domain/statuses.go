package domain

import "sort"

const (
	RoleAdmin    = "ADMIN"
	RoleOperator = "OPERATOR"
	RoleReviewer = "REVIEWER"
	RoleAuditor  = "AUDITOR"
)

const (
	UserStatusActive   = "ACTIVE"
	UserStatusDisabled = "DISABLED"
)

const (
	DeviceStatusNormal        = "NORMAL"
	DeviceStatusRepairing     = "REPAIRING"
	DeviceStatusReported      = "REPORTED"
	DeviceStatusReadyForPickup = "READY_FOR_PICKUP"
	DeviceStatusPickedUp      = "PICKED_UP"
	DeviceStatusScrapped      = "SCRAPPED"
)

const (
	RepairOrderStatusPending         = "PENDING"
	RepairOrderStatusDispatched      = "DISPATCHED"
	RepairOrderStatusWaitingVisit    = "WAITING_VISIT"
	RepairOrderStatusWaitingDelivery = "WAITING_DELIVERY"
	RepairOrderStatusAccepted        = "ACCEPTED"
	RepairOrderStatusCompleted       = "COMPLETED"
	RepairOrderStatusCancelled       = "CANCELLED"
)

const (
	QuotationStatusDraft         = "DRAFT"
	QuotationStatusWaitingClient = "WAITING_CLIENT"
	QuotationStatusConfirmed     = "CONFIRMED"
	QuotationStatusRejected      = "REJECTED"
	QuotationStatusExpired       = "EXPIRED"
	QuotationStatusConverted     = "CONVERTED"
)

const (
	ApprovalNotRequired = "NOT_REQUIRED"
	ApprovalPending     = "PENDING"
	ApprovalApproved    = "APPROVED"
	ApprovalRejected    = "REJECTED"
)

const (
	RepairExecutionStatusPending      = "PENDING"
	RepairExecutionStatusInProgress   = "IN_PROGRESS"
	RepairExecutionStatusCompleted    = "COMPLETED"
	RepairExecutionStatusWaitingTest  = "WAITING_TEST"
	RepairExecutionStatusTested       = "TESTED"
	RepairExecutionStatusWaitingPickup = "WAITING_PICKUP"
	RepairExecutionStatusDelivered    = "DELIVERED"
)

const (
	WarrantyStatusActive  = "ACTIVE"
	WarrantyStatusExpired = "EXPIRED"
	WarrantyStatusVoid    = "VOID"
)

const (
	FeedbackStatusPending = "PENDING"
	FeedbackStatusVisited = "VISITED"
	FeedbackStatusClosed  = "CLOSED"
)

const (
	ServiceMethodOnSite  = "ON_SITE"
	ServiceMethodDropOff = "DROP_OFF"
	ServiceMethodRemote  = "REMOTE"
)

const (
	UrgencyNormal   = "NORMAL"
	UrgencyUrgent   = "URGENT"
	UrgencyCritical = "CRITICAL"
)

const (
	CustomerLevelNormal     = "NORMAL"
	CustomerLevelVIP        = "VIP"
	CustomerLevelEnterprise = "ENTERPRISE"
)

const (
	PaymentMethodCash        = "CASH"
	PaymentMethodWechat      = "WECHAT"
	PaymentMethodAlipay      = "ALIPAY"
	PaymentMethodBankCard    = "BANK_CARD"
	PaymentMethodCorporate   = "CORPORATE_TRANSFER"
)

const (
	FeedbackMethodPhone = "PHONE"
	FeedbackMethodOnline = "ONLINE"
	FeedbackMethodOnSite = "ON_SITE"
)

var DeviceStatusTransitions = map[string][]string{
	DeviceStatusNormal:         {DeviceStatusReported, DeviceStatusScrapped},
	DeviceStatusReported:       {DeviceStatusRepairing, DeviceStatusNormal},
	DeviceStatusRepairing:      {DeviceStatusReadyForPickup, DeviceStatusScrapped},
	DeviceStatusReadyForPickup: {DeviceStatusPickedUp},
	DeviceStatusPickedUp:       {},
	DeviceStatusScrapped:       {},
}

var RepairOrderStatusTransitions = map[string][]string{
	RepairOrderStatusPending:         {RepairOrderStatusDispatched, RepairOrderStatusCancelled},
	RepairOrderStatusDispatched:      {RepairOrderStatusWaitingVisit, RepairOrderStatusWaitingDelivery, RepairOrderStatusAccepted, RepairOrderStatusCancelled},
	RepairOrderStatusWaitingVisit:    {RepairOrderStatusAccepted, RepairOrderStatusCancelled},
	RepairOrderStatusWaitingDelivery: {RepairOrderStatusAccepted, RepairOrderStatusCancelled},
	RepairOrderStatusAccepted:        {RepairOrderStatusCompleted, RepairOrderStatusCancelled},
	RepairOrderStatusCompleted:       {},
	RepairOrderStatusCancelled:       {},
}

var QuotationStatusTransitions = map[string][]string{
	QuotationStatusDraft:         {QuotationStatusWaitingClient},
	QuotationStatusWaitingClient: {QuotationStatusConfirmed, QuotationStatusRejected, QuotationStatusExpired},
	QuotationStatusConfirmed:     {QuotationStatusConverted},
	QuotationStatusRejected:      {},
	QuotationStatusExpired:       {},
	QuotationStatusConverted:     {},
}

var RepairExecutionStatusTransitions = map[string][]string{
	RepairExecutionStatusPending:       {RepairExecutionStatusInProgress},
	RepairExecutionStatusInProgress:    {RepairExecutionStatusCompleted, RepairExecutionStatusWaitingTest},
	RepairExecutionStatusCompleted:     {RepairExecutionStatusWaitingTest},
	RepairExecutionStatusWaitingTest:   {RepairExecutionStatusTested},
	RepairExecutionStatusTested:        {RepairExecutionStatusWaitingPickup},
	RepairExecutionStatusWaitingPickup: {RepairExecutionStatusDelivered},
	RepairExecutionStatusDelivered:     {},
}

var WarrantyStatusTransitions = map[string][]string{
	WarrantyStatusActive:  {WarrantyStatusExpired, WarrantyStatusVoid},
	WarrantyStatusExpired: {},
	WarrantyStatusVoid:    {},
}

var FeedbackStatusTransitions = map[string][]string{
	FeedbackStatusPending: {FeedbackStatusVisited, FeedbackStatusClosed},
	FeedbackStatusVisited: {FeedbackStatusClosed},
	FeedbackStatusClosed:  {},
}

func AllowedUserStatuses() []string {
	return []string{UserStatusActive, UserStatusDisabled}
}

func AllowedCustomerLevels() []string {
	return []string{CustomerLevelNormal, CustomerLevelVIP, CustomerLevelEnterprise}
}

func AllowedServiceMethods() []string {
	return []string{ServiceMethodOnSite, ServiceMethodDropOff, ServiceMethodRemote}
}

func AllowedUrgencies() []string {
	return []string{UrgencyNormal, UrgencyUrgent, UrgencyCritical}
}

func AllowedPaymentMethods() []string {
	return []string{
		PaymentMethodCash,
		PaymentMethodWechat,
		PaymentMethodAlipay,
		PaymentMethodBankCard,
		PaymentMethodCorporate,
	}
}

func AllowedFeedbackMethods() []string {
	return []string{FeedbackMethodPhone, FeedbackMethodOnline, FeedbackMethodOnSite}
}

func AllowedDeviceStatuses() []string {
	return sortedKeys(DeviceStatusTransitions)
}

func AllowedRepairOrderStatuses() []string {
	return sortedKeys(RepairOrderStatusTransitions)
}

func AllowedQuotationStatuses() []string {
	return sortedKeys(QuotationStatusTransitions)
}

func AllowedRepairExecutionStatuses() []string {
	return sortedKeys(RepairExecutionStatusTransitions)
}

func AllowedWarrantyStatuses() []string {
	return sortedKeys(WarrantyStatusTransitions)
}

func AllowedFeedbackStatuses() []string {
	return sortedKeys(FeedbackStatusTransitions)
}

func IsAllowedTransition(current string, next string, transitions map[string][]string) bool {
	allowed, ok := transitions[current]
	if !ok {
		return false
	}
	for _, candidate := range allowed {
		if candidate == next {
			return true
		}
	}
	return false
}

func IsAllowedValue(value string, values []string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}

func StatusCatalogue() map[string][]string {
	return map[string][]string{
		"user":              AllowedUserStatuses(),
		"device":            AllowedDeviceStatuses(),
		"repair_order":      AllowedRepairOrderStatuses(),
		"quotation":         AllowedQuotationStatuses(),
		"repair_execution":  AllowedRepairExecutionStatuses(),
		"warranty":          AllowedWarrantyStatuses(),
		"feedback":          AllowedFeedbackStatuses(),
		"service_method":    AllowedServiceMethods(),
		"urgency":           AllowedUrgencies(),
		"customer_level":    AllowedCustomerLevels(),
		"payment_method":    AllowedPaymentMethods(),
		"feedback_method":   AllowedFeedbackMethods(),
	}
}

func RoleCatalogue() []string {
	return []string{RoleAdmin, RoleOperator, RoleReviewer, RoleAuditor}
}

func RoleDisplayNames() map[string]string {
	return map[string]string{
		RoleAdmin:    "Business Admin",
		RoleOperator: "Operator",
		RoleReviewer: "Reviewer",
		RoleAuditor:  "Auditor",
	}
}

func sortedKeys(values map[string][]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
