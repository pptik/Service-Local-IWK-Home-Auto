package consumer

import (
	"encoding/json"
	"go/hioto/pkg/dto"
	"go/hioto/pkg/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
)

var validate = validator.New()

type ConsumerHandler struct {
	deviceService        *service.DeviceService
	controlDeviceService *service.ControlDeviceService
	validator            *validator.Validate
}

func NewConsumerHandler(deviceService *service.DeviceService, controlDeviceService *service.ControlDeviceService) *ConsumerHandler {
	return &ConsumerHandler{
		deviceService:        deviceService,
		controlDeviceService: controlDeviceService,
		validator:            validator.New(),
	}
}

func (h *ConsumerHandler) ControlHandler(message []byte) {
	var controlDeviceDto dto.ControlLocalDto

	if err := json.Unmarshal(message, &controlDeviceDto); err != nil {
		log.Errorf("Failed to unmarshal control message: %v", err)
		return
	}

	if err := validate.Struct(controlDeviceDto); err != nil {
		log.Errorf("Validation error: %v", err)
		return
	}
}

func (h *ConsumerHandler) TestingConsumeAktuator(message []byte) {
	messageString := string(message)

	log.Info(messageString)
}
