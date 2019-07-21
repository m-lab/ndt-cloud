// Package measurer contains the measurer. Its job is to collect metrics
// from a socket connection and return them for consumption.
package measurer

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m-lab/ndt-server/bbr"
	"github.com/m-lab/ndt-server/fdcache"
	"github.com/m-lab/ndt-server/logging"
	"github.com/m-lab/ndt-server/ndt7/model"
	"github.com/m-lab/ndt-server/ndt7/results"
	"github.com/m-lab/ndt-server/ndt7/spec"
	"github.com/m-lab/ndt-server/tcpinfox"
)

func getSocketAndPossiblyEnableBBR(conn *websocket.Conn) (*os.File, error) {
	fp := fdcache.GetAndForgetFile(conn.UnderlyingConn())
	// Implementation note: in theory fp SHOULD always be non-nil because
	// now we always register the fp bound to a net.TCPConn. However, in
	// some weird cases it MAY happen that the cache pruning mechanism will
	// remove the fp BEFORE we can steal it. In case we cannot get a file
	// we just abort the test, as this should not happen (TM).
	if fp == nil {
		return nil, errors.New("cannot get file bound to websocket conn")
	}
	err := bbr.Enable(fp)
	if err != nil {
		logging.Logger.WithError(err).Warn("Cannot enable BBR")
		// FALLTHROUGH
	}
	return fp, nil
}

func measure(measurement *model.Measurement, sockfp *os.File) {
	bbrinfo, err := bbr.GetMaxBandwidthAndMinRTT(sockfp)
	if err == nil {
		measurement.BBRInfo = &bbrinfo
	}
	metrics, err := tcpinfox.GetTCPInfo(sockfp)
	if err == nil {
		measurement.TCPInfo = &metrics
	}
}

func measureAndSendToChannel(
	t0, t1 time.Time, sockfp *os.File, dst chan<- model.Measurement,
	connectionInfo *model.ConnectionInfo,
) {
	measurement := model.Measurement{
		ConnectionInfo: connectionInfo,
		Elapsed:        t1.Sub(t0).Seconds(),
	}
	measure(&measurement, sockfp)
	dst <- measurement
}

func loop(
	ctx context.Context, conn *websocket.Conn, resultsfp *results.File,
	dst chan<- model.Measurement,
) {
	logging.Logger.Debug("measurer: start")
	defer logging.Logger.Debug("measurer: stop")
	defer close(dst)
	measurerctx, cancel := context.WithTimeout(ctx, spec.DefaultRuntime)
	defer cancel()
	sockfp, err := getSocketAndPossiblyEnableBBR(conn)
	if err != nil {
		logging.Logger.WithError(err).Warn("getSocketAndPossiblyEnableBBR failed")
		return
	}
	defer sockfp.Close()
	t0 := time.Now()
	resultsfp.StartTest()
	defer resultsfp.EndTest()
	// Liveness: this is probably non blocking because of buffering
	measureAndSendToChannel(t0, time.Now(), sockfp, dst, &model.ConnectionInfo{
		Client: conn.RemoteAddr().String(),
		Server: conn.LocalAddr().String(),
		UUID:   resultsfp.Data.UUID,
	})
	ticker := time.NewTicker(spec.MinMeasurementInterval)
	defer ticker.Stop()
	for {
		select {
		case <-measurerctx.Done(): // Liveness!
			logging.Logger.Debug("measurer: context done")
			return
		case now := <-ticker.C:
			// Liveness: this is probably non blocking because of buffering
			measureAndSendToChannel(t0, now, sockfp, dst, nil)
		}
	}
}

// Start runs the measurement loop in a background goroutine and emits
// the measurements on the returned channel.
//
// Liveness guarantee: the measurer will always terminate after
// a timeout of DefaultRuntime seconds, provided that the consumer
// continues reading from the returned channel.
func Start(
	ctx context.Context, conn *websocket.Conn, resultsfp *results.File,
) <-chan model.Measurement {
	// We use buffering to make sure that the measurer can bufferise the total
	// number of expected measurements without blocking. This allows us to obtain
	// measurements even when the sender blocks for a long time. In principle it
	// may still blocking for some time, but we'll get most measurements.
	//
	// Such measurements won't probably be sent to the client, since the sender
	// would spent most of its time blocked. But it should still be possible for
	// the sender to drain our channel when exiting. This way measurements will
	// be saved on the disk and will be available for analysis.
	dst := make(chan model.Measurement, spec.NumExpectedMeasurements)
	go loop(ctx, conn, resultsfp, dst)
	return dst
}
