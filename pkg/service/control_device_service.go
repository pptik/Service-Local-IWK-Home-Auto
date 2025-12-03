package service

import (
	"fmt"
	"go/hioto/config"
	"go/hioto/pkg/dto"
	"go/hioto/pkg/enum"
	messagebroker "go/hioto/pkg/handler/message_broker"
	"go/hioto/pkg/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ControlDeviceService struct {
	db *gorm.DB
}

func NewControlDeviceService(db *gorm.DB) *ControlDeviceService {
	return &ControlDeviceService{
		db: db,
	}
}

func (s *ControlDeviceService) ControlDeviceLocal(controlDto *dto.ControlLocalDto) error {
	var device model.Registration
	value := strings.Split(controlDto.Message, "#")

	if err := s.db.Where("guid = ?", value[0]).First(&device).Error; err != nil {
		log.Errorf("Device not found: %v 💥", err)
		return fiber.NewError(fiber.StatusNotFound, "Device is not found")
	}

	if controlDto.Type == enum.SENSOR {
		s.ControlSensor(value[0], value[1])
		return nil
	}

	location = time.FixedZone("WIB", 7*60*60)

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Errorf("Transaction rollback due to panic: %v 💥", r)
		} else {
			if err := tx.Commit().Error; err != nil {
				log.Errorf("Error committing the transaction: %v 💥", err)
				tx.Rollback()
			}
		}
	}()

	if tx.Error != nil {
		log.Errorf("Error starting transaction: %v 💥", tx.Error)
		return fiber.NewError(fiber.StatusBadRequest, "Error starting transaction")
	}

	status := value[1] == "1"

	if err := tx.Model(&device).Updates(map[string]any{
		"status":     status,
		"updated_at": time.Now().In(location),
	}).Error; err != nil {
		log.Errorf("Error updating registration: %v 💥", err)
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Error updating registration device")
	}

	messagebroker.PublishToRoutingKey(
		config.RMQ_INSTANCE.GetValue(),
		[]byte(controlDto.Message),
		config.EXCHANGE_TOPIC.GetValue(),
		config.AKTUATOR_ROUTING_KEY.GetValue(),
	)

	return nil
}

func (s *ControlDeviceService) ControlSensor(guid, value string) {
	var ruleDevices []model.RuleDevice

	if err := s.db.Where("input_guid = ?", guid).Where("input_value = ?", value).Find(&ruleDevices).Error; err != nil {
		log.Errorf("Failed to fetch rule devices: %v 💥", err)
		return
	}

	if len(ruleDevices) == 0 {
		log.Error("No rule devices found 💥")
		return
	}

	for _, ruleDevice := range ruleDevices {
		var aktuator model.Registration

		location = time.FixedZone("WIB", 7*60*60)

		messageToAktuator := fmt.Sprintf("%s#%s", ruleDevice.OutputGuid, ruleDevice.OutputValue)

		if err := s.db.Where("guid = ?", ruleDevice.OutputGuid).First(&aktuator).Error; err != nil {
			log.Errorf("Failed to fetch aktuator: %v 💥", err)
			continue
		}

		aktuator.Status = ruleDevice.OutputValue
		aktuator.UpdatedAt = time.Now().In(location)

		if err := s.db.Save(&aktuator).Error; err != nil {
			log.Errorf("Failed update aktuator status: %v 💥", err)
			continue
		}

		messagebroker.PublishToRoutingKey(
			config.RMQ_INSTANCE.GetValue(),
			[]byte(messageToAktuator),
			config.EXCHANGE_TOPIC.GetValue(),
			config.AKTUATOR_ROUTING_KEY.GetValue(),
		)

		log.Infof("Sensor rule executed for aktuator %s with value %s ✅", aktuator.Name, ruleDevice.OutputValue)
	}
}
