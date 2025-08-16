package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/shared-drawboard/internal/models"
)

func GenerateUniqueID() string {
	return fmt.Sprintf("client-%d", time.Now().UnixNano())
}

func ParseEventData(rawData []byte) (models.Event, error) {
	// First pass to get event type
	var baseEvent struct {
		Type models.EventType `json:"type"`
		Tool string           `json:"tool"`
		Data json.RawMessage  `json:"data"`
	}

	if err := json.Unmarshal(rawData, &baseEvent); err != nil {
		return models.Event{}, err
	}

	var event models.Event
	event.Type = baseEvent.Type
	event.Tool = baseEvent.Tool

	// Second pass for specific data types
	switch baseEvent.Type {
	case models.FreehandDraw:
		var data models.FreehandDrawData
		if err := json.Unmarshal(baseEvent.Data, &data); err != nil {
			return models.Event{}, err
		}
		event.Data = data

	case models.ShapeCreate:
		var data models.ShapeCreateData
		if err := json.Unmarshal(baseEvent.Data, &data); err != nil {
			return models.Event{}, err
		}
		event.Data = data

	case models.TextAdd:
		var data models.TextAddData
		if err := json.Unmarshal(baseEvent.Data, &data); err != nil {
			return models.Event{}, err
		}
		event.Data = data

	case models.ObjectDelete:
		var data models.ObjectDeleteData
		if err := json.Unmarshal(baseEvent.Data, &data); err != nil {
			return models.Event{}, err
		}
		event.Data = data

	case models.BoardClear:
		// No data payload needed
		event.Data = nil

	default:
		return models.Event{}, errors.New("unknown event type")
	}

	return event, nil
}
