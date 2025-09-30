package storage

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	dir     string
	files   map[string]*os.File
	offsets map[string]int64
}

func NewLog(dir string) (*Log, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Log{
		dir:     dir,
		files:   make(map[string]*os.File),
		offsets: make(map[string]int64),
	}, nil
}

func (l *Log) getFile(topic string) (*os.File, error) {
	if f, ok := l.files[topic]; ok {
		return f, nil
	}
	path := filepath.Join(l.dir, topic+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	// Recover offset
	stat, _ := f.Stat()
	if stat.Size() == 0 {
		l.offsets[topic] = 0
	} else {
		offset, err := l.recoverOffset(f)
		if err != nil {
			return nil, err
		}
		l.offsets[topic] = offset
	}
	l.files[topic] = f
	return f, nil
}

func (l *Log) recoverOffset(f *os.File) (int64, error) {
	var offset int64 = 0
	_, err := f.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	for {
		lenBuf := make([]byte, 4)
		_, err := f.Read(lenBuf)
		if err != nil {
			break
		}
		msgLen := binary.BigEndian.Uint32(lenBuf)
		_, err = f.Seek(int64(msgLen), 1)
		if err != nil {
			break
		}
		offset++
	}
	return offset, nil
}

func (l *Log) Append(topic string, msg *Message) (int64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := l.getFile(topic)
	if err != nil {
		return 0, err
	}
	msg.Offset = l.offsets[topic]
	b, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(b)))
	if _, err := f.Write(lenBuf); err != nil {
		return 0, err
	}
	if _, err := f.Write(b); err != nil {
		return 0, err
	}
	if err := f.Sync(); err != nil {
		return 0, err
	}
	l.offsets[topic]++
	return msg.Offset, nil
}

func (l *Log) Read(topic string, offset int64, max int) ([]*Message, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := l.getFile(topic)
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	var msgs []*Message
	var cur int64 = 0
	for {
		lenBuf := make([]byte, 4)
		_, err := f.Read(lenBuf)
		if err != nil {
			break
		}
		msgLen := binary.BigEndian.Uint32(lenBuf)
		msgBuf := make([]byte, msgLen)
		_, err = f.Read(msgBuf)
		if err != nil {
			break
		}
		if cur >= offset {
			var m Message
			if err := json.Unmarshal(msgBuf, &m); err == nil {
				msgs = append(msgs, &m)
			}
			if len(msgs) >= max {
				break
			}
		}
		cur++
	}
	return msgs, nil
}

func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, f := range l.files {
		f.Close()
	}
	return nil
}
