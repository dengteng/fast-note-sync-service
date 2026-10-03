package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeUploadSession struct{ lastActive time.Time }

func (s *fakeUploadSession) GetLastActiveAt() time.Time { return s.lastActive }

// TestCleanupStaleSessionsByLastActive: an upload created long ago but still receiving chunks must survive
// another connection of the same user closing; only idle sessions are removed.
// TestCleanupStaleSessionsByLastActive：创建很久但仍在收分片的上传，不能因同用户别的连接断开被清掉；只清闲置会话。
func TestCleanupStaleSessionsByLastActive(t *testing.T) {
	w := &WebsocketServer{binaryChunkSessions: map[string]map[string]any{
		"1": {
			"active": &fakeUploadSession{lastActive: time.Now().Add(-5 * time.Second)},
			"idle":   &fakeUploadSession{lastActive: time.Now().Add(-11 * time.Minute)},
		},
	}}

	w.cleanupStaleSessions("1", 10*time.Minute)

	assert.Contains(t, w.binaryChunkSessions["1"], "active")
	assert.NotContains(t, w.binaryChunkSessions["1"], "idle")
}
