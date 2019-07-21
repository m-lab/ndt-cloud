// Package spec contains constants defined in the ndt7 specification.
package spec

import "time"

// DownloadURLPath selects the download subtest.
const DownloadURLPath = "/ndt/v7/download"

// UploadURLPath selects the upload subtest.
const UploadURLPath = "/ndt/v7/upload"

// SecWebSocketProtocol is the WebSocket subprotocol used by ndt7.
const SecWebSocketProtocol = "net.measurementlab.ndt.v7"

// MinMessageSize is the minimum message size.
const MinMessageSize = 1 << 10

// InitialMessageSize is the initial message size.
const InitialMessageSize = 1 << 13

// MaxMessageSize is the maximum message size.
const MaxMessageSize = 1 << 24

// MeasurementsPerSecond is the number of measurements per second.
const MeasurementsPerSecond = 4

// defaultRuntimeSeconds is the default runtime in seconds.
const defaultRuntimeSecoonds = 10

// NumExpectedMeasurements is the expected number of measurements.
const NumExpectedMeasurements = defaultRuntimeSeconds * MeasurementsPerSecond

// MinMeasurementInterval is the minimum interval between measurements.
const MinMeasurementInterval = (1000 / MeasurementsPerSecond) * time.Millisecond

// DefaultRuntime is the default runtime of a subtest
const DefaultRuntime = defaultRuntimeSeconds * time.Second

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
