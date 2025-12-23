// status.go
package constants

type TripStatus string

const (
	TripPublished  TripStatus = "published"
	TripInProgress TripStatus = "in_progress"
	TripCompleted  TripStatus = "completed"
)

type BookingStatus string

const (
	BookingPending  = "pending"   // заявка отправлена
	BookingApproved = "approved"  // водитель принял
	BookingRejected = "rejected"  // водитель отклонил
)
