package uuid_test

import (
	"testing"

	"github.com/notes-bin/calendar-go/pkg/uuid"
)

func TestGenerateUUID(t *testing.T) {
	uuid, err := uuid.Generate()
	if err != nil {
		t.Errorf("生成UUID失败: %v", err)
	}
	t.Logf("生成的UUID: %s", uuid)
}
