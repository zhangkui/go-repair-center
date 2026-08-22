package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type OrderNumberService struct {
	repairPrefix   string
	feedbackPrefix string
	quotationPrefix string
	now            func() time.Time
}

type ParsedOrderNumber struct {
	Prefix    string    `json:"prefix"`
	Date      time.Time `json:"date"`
	Sequence  int       `json:"sequence"`
	Raw       string    `json:"raw"`
}

func NewOrderNumberService() *OrderNumberService {
	return &OrderNumberService{
		repairPrefix:    "R",
		feedbackPrefix:  "F",
		quotationPrefix: "Q",
		now:             time.Now,
	}
}

func (s *OrderNumberService) GenerateRepairOrderNumber(sequence int) string {
	return s.generate(s.repairPrefix, s.now(), sequence)
}

func (s *OrderNumberService) GenerateFeedbackNumber(sequence int) string {
	return s.generate(s.feedbackPrefix, s.now(), sequence)
}

func (s *OrderNumberService) GenerateQuotationNumber(sequence int) string {
	return s.generate(s.quotationPrefix, s.now(), sequence)
}

func (s *OrderNumberService) GenerateRepairOrderNumberAt(date time.Time, sequence int) string {
	return s.generate(s.repairPrefix, date, sequence)
}

func (s *OrderNumberService) SequenceKey(date time.Time) string {
	return fmt.Sprintf("sequence:%s", date.Format("20060102"))
}

func (s *OrderNumberService) ValidateRepairOrderNumber(value string) bool {
	parsed, err := s.Parse(value)
	return err == nil && parsed.Prefix == s.repairPrefix
}

func (s *OrderNumberService) Parse(value string) (*ParsedOrderNumber, error) {
	value = strings.TrimSpace(strings.ToUpper(value))
	if len(value) != 16 {
		return nil, fmt.Errorf("invalid order number length: %s", value)
	}
	prefix := value[:1]
	datePart := value[1:9]
	sequencePart := value[9:]

	date, err := time.Parse("20060102", datePart)
	if err != nil {
		return nil, fmt.Errorf("invalid order date: %w", err)
	}
	sequence, err := strconv.Atoi(sequencePart)
	if err != nil {
		return nil, fmt.Errorf("invalid order sequence: %w", err)
	}
	return &ParsedOrderNumber{
		Prefix:   prefix,
		Date:     date,
		Sequence: sequence,
		Raw:      value,
	}, nil
}

func (s *OrderNumberService) NormalizeIdempotencyKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.ReplaceAll(value, "_", "-")
	return strings.Trim(value, "-")
}

func (s *OrderNumberService) BuildDispatchReference(orderNumber string, technicianID int64) string {
	return fmt.Sprintf("%s-T%06d", orderNumber, technicianID)
}

func (s *OrderNumberService) BuildWarrantyReference(orderNumber string, iteration int) string {
	if iteration <= 0 {
		iteration = 1
	}
	return fmt.Sprintf("%s-W%02d", orderNumber, iteration)
}

func (s *OrderNumberService) BuildReworkNumber(originalOrder string, index int) string {
	if index <= 0 {
		index = 1
	}
	return fmt.Sprintf("%s-RW%02d", originalOrder, index)
}

func (s *OrderNumberService) DayBucket(value string) string {
	parsed, err := s.Parse(value)
	if err != nil {
		return ""
	}
	return parsed.Date.Format("2006-01-02")
}

func (s *OrderNumberService) Compare(left string, right string) int {
	leftParsed, leftErr := s.Parse(left)
	rightParsed, rightErr := s.Parse(right)
	if leftErr != nil || rightErr != nil {
		return strings.Compare(left, right)
	}
	if leftParsed.Date.Before(rightParsed.Date) {
		return -1
	}
	if leftParsed.Date.After(rightParsed.Date) {
		return 1
	}
	if leftParsed.Sequence < rightParsed.Sequence {
		return -1
	}
	if leftParsed.Sequence > rightParsed.Sequence {
		return 1
	}
	return 0
}

func (s *OrderNumberService) NextSequence(existing []string, date time.Time) int {
	maxSequence := 0
	for _, item := range existing {
		parsed, err := s.Parse(item)
		if err != nil {
			continue
		}
		if parsed.Date.Format("20060102") != date.Format("20060102") {
			continue
		}
		if parsed.Sequence > maxSequence {
			maxSequence = parsed.Sequence
		}
	}
	return maxSequence + 1
}

func (s *OrderNumberService) BatchGenerate(prefix string, date time.Time, start int, count int) []string {
	if count <= 0 {
		return []string{}
	}
	items := make([]string, 0, count)
	for index := 0; index < count; index++ {
		items = append(items, s.generate(prefix, date, start+index))
	}
	return items
}

func (s *OrderNumberService) SuggestSearchKeywords(orderNumber string) []string {
	value := strings.TrimSpace(strings.ToUpper(orderNumber))
	if value == "" {
		return nil
	}
	keywords := []string{value}
	if parsed, err := s.Parse(value); err == nil {
		keywords = append(keywords,
			parsed.Date.Format("20060102"),
			parsed.Date.Format("2006-01-02"),
			fmt.Sprintf("%08d", parsed.Sequence),
		)
	}
	return deduplicateStrings(keywords)
}

func (s *OrderNumberService) ToHumanLabel(orderNumber string) string {
	parsed, err := s.Parse(orderNumber)
	if err != nil {
		return orderNumber
	}
	return fmt.Sprintf("%s %s #%d", parsed.Prefix, parsed.Date.Format("2006-01-02"), parsed.Sequence)
}

func (s *OrderNumberService) Prefixes() []string {
	return []string{s.repairPrefix, s.feedbackPrefix, s.quotationPrefix}
}

func (s *OrderNumberService) IsKnownPrefix(prefix string) bool {
	prefix = strings.ToUpper(strings.TrimSpace(prefix))
	for _, item := range s.Prefixes() {
		if item == prefix {
			return true
		}
	}
	return false
}

func (s *OrderNumberService) GuessPrefix(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if len(value) == 0 {
		return ""
	}
	prefix := value[:1]
	if s.IsKnownPrefix(prefix) {
		return prefix
	}
	return ""
}

func (s *OrderNumberService) generate(prefix string, date time.Time, sequence int) string {
	if sequence < 0 {
		sequence = 0
	}
	return fmt.Sprintf("%s%s%07d", strings.ToUpper(prefix), date.Format("20060102"), sequence)
}

func deduplicateStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	items := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		items = append(items, value)
	}
	return items
}
