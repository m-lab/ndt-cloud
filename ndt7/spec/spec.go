// Package spec contains constants defined in the ndt7 specification.
package spec

import "time"

// DownloadURLPath selects the download subtest.
const DownloadURLPath = "/ndt/v7/download"

// UploadURLPath selects the upload subtest.
const UploadURLPath = "/ndt/v7/upload"

// SecWebSocketProtocol is the WebSocket subprotocol used by ndt7.
const SecWebSocketProtocol = "net.measurementlab.ndt.v7"

// MaxMessageSize is the minimum value of the maximum message size
// that an implementation MAY want to configure. Messages smaller than this
// threshold MUST always be accepted by an implementation.
const MaxMessageSize = 1 << 24

// MaxScaledMessageSize is the maximum value of a scaled binary WebSocket
// message size. This should be <= of MaxMessageSize. The 1<<20 value is
// a good compromise between Go and JavaScript as seen in cloud based tests.
const MaxScaledMessageSize = 1 << 20

// ScalingFraction sets the threshold for scaling binary messages. When
// the current binary message size is <= than 1/scalingFactor of the
// amount of bytes sent so far, we scale the message. This is documented
// in the appendix of the ndt7 specification.
const ScalingFraction = 16

// MinMeasurementInterval is the minimum interval between measurements.
const MinMeasurementInterval = 250 * time.Millisecond

// DefaultRuntime is the default runtime of a subtest
const DefaultRuntime = 10 * time.Second

// MaxRuntime is the maximum runtime of a subtest
const MaxRuntime = 15 * time.Second

// SubtestKind indicates the subtest kind
type SubtestKind string

const (
	// SubtestDownload is a download subtest
	SubtestDownload = SubtestKind("download")

	// SubtestUpload is a upload subtest
	SubtestUpload = SubtestKind("upload")
)
