// status.go
package constants

type TripStatus string

const (
	TripPublished  TripStatus = "published"
	TripInProgress TripStatus = "in_progress"
	TripCompleted  TripStatus = "completed"
)
