package camera

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HellsKitchen99/camera-processor/internal/domain"
	rtsp "github.com/bluenviron/gortsplib/v5"
	"github.com/bluenviron/gortsplib/v5/pkg/base"
	"github.com/bluenviron/gortsplib/v5/pkg/format"
	"github.com/bluenviron/gortsplib/v5/pkg/format/rtph264"
	"github.com/bluenviron/mediacommon/v2/pkg/codecs/h264"
	"github.com/pion/rtp"
	"github.com/sirupsen/logrus"
)

type CameraReader struct {
	CamerasUrl []domain.Camera
	Jobs       chan<- domain.FrameJob
	wg         *sync.WaitGroup

	ctx context.Context
}

func NewCameraReader(camerasUrl []domain.Camera, jobs chan<- domain.FrameJob, ctx context.Context, wg *sync.WaitGroup) *CameraReader {
	return &CameraReader{
		CamerasUrl: camerasUrl,
		Jobs:       jobs,
		wg:         wg,

		ctx: ctx,
	}
}

func (c *CameraReader) Run() {
	for i := 0; i < len(c.CamerasUrl); i++ {
		c.wg.Add(1)
		go func(cameraId int, cameraUrl string) {
			defer c.wg.Done()
			c.reader(cameraId, cameraUrl)
		}(c.CamerasUrl[i].ID, c.CamerasUrl[i].URL)
	}
}

func (c *CameraReader) reader(cameraId int, cameraUrl string) {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if err := readCamera(cameraId, cameraUrl, c.Jobs, c.ctx); err != nil {
				if c.ctx.Err() != nil {
					return
				}
				logrus.Errorf("CAMERA %v: %v", cameraId, err)
			}

			select {
			case <-c.ctx.Done():
				return
			case <-time.After(time.Second * 2):
			}
		}
	}
}

func readCamera(cameraId int, cameraUrl string, jobs chan<- domain.FrameJob, ctx context.Context) error {

	// parse URL
	u, err := base.ParseURL(cameraUrl)
	if err != nil {
		return err
	}

	// enable RTSP over TCP
	protocol := rtsp.ProtocolTCP

	// create a client
	client := rtsp.Client{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Protocol: &protocol,
	}

	// connect to rtsp server (camera)
	if err := client.Start(); err != nil {
		return err
	}
	defer client.Close()

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			client.Close()
		case <-done:
			return
		}
	}()
	defer close(done)

	// get stream description
	desc, _, err := client.Describe(u)
	if err != nil {
		return err
	}

	// get h264 format from stream
	var h264Format *format.H264
	media := desc.FindFormat(&h264Format)
	if media == nil {
		return fmt.Errorf("media is nil")
	}

	// create rtp -> h264 decoder
	decoder, err := h264Format.CreateDecoder()
	if err != nil {
		return err
	}

	// create h264 -> image decoder
	h264Decoder := &h264Decoder{}
	if err := h264Decoder.initialize(); err != nil {
		return err
	}
	defer h264Decoder.close()

	// getting sequence and picture parameters sets
	if h264Format.SPS != nil {
		if _, err := h264Decoder.decode([][]byte{h264Format.SPS}); err != nil {
			return err
		}
	}
	if h264Format.PPS != nil {
		if _, err := h264Decoder.decode([][]byte{h264Format.PPS}); err != nil {
			return err
		}
	}

	// setup media path
	if _, err := client.Setup(desc.BaseURL, media, 0, 0); err != nil {
		return err
	}

	lastSent := time.Time{}

	// handle rtp stream
	firstRandomAccess := false
	client.OnPacketRTP(media, h264Format, func(p *rtp.Packet) {
		// get a timing
		_, ok := client.PacketPTS(media, p)
		if !ok {
			fmt.Printf("CAMERA %v: waiting for a timing\n", cameraId)
			return
		}

		// build a frame (rtp -> h264)
		au, err := decoder.Decode(p)
		if err != nil {
			if errors.Is(err, rtph264.ErrMorePacketsNeeded) || errors.Is(err, rtph264.ErrNonStartingPacketAndNoPrevious) {
				return
			}
			logrus.Errorf("CAMERA %v: RTP decode error: %v", cameraId, err)
			return
		}

		// waiting for a random access unit
		if !firstRandomAccess && !h264.IsRandomAccess(au) {
			return
		}
		firstRandomAccess = true

		// build image (h264 -> image)
		image, err := h264Decoder.decode(au)
		if err != nil {
			return
		}

		// check image presence
		if image == nil {
			return
		}

		// sucessful log
		//logrus.Infof("CAMERA %v: Decode frame with size %v and timing\n", cameraId, image.Bounds().Max, pts)

		if lastSent.Add(time.Second).After(time.Now()) {
			return
		}

		lastSent = time.Now()

		select {
		case jobs <- domain.FrameJob{
			CameraID:   cameraId,
			Image:      image,
			ReceivedAt: time.Now(),
		}:
			// successful sending
			//logrus.Infof("CAMERA %v: sent frame successfully", cameraId)

		default:
			// jobs is full
			logrus.Warnf("CAMERA %v: drop frame, jobs queue is full", cameraId)
		}
	})

	// play stream
	if _, err := client.Play(nil); err != nil {
		return err
	}

	// wait
	if err := client.Wait(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
}
