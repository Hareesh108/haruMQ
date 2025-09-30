package storage

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	dir     string
	files   map[string]map[int]*os.File // topic -> partition -> file
	offsets map[string]map[int]int64    // topic -> partition -> offset
}

func NewLog(dir string) (*Log, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Log{
		dir:     dir,
		files:   make(map[string]map[int]*os.File),
		offsets: make(map[string]map[int]int64),
	}, nil
}

func (l *Log) getFile(topic string, partition int) (*os.File, error) {
	if l.files[topic] == nil {
		l.files[topic] = make(map[int]*os.File)
	}
	if l.offsets[topic] == nil {
		l.offsets[topic] = make(map[int]int64)
	}
	if f, ok := l.files[topic][partition]; ok {
		return f, nil
	}
	path := filepath.Join(l.dir, topic+"-"+strconv.Itoa(partition)+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	// Recover offset
	stat, _ := f.Stat()
	if stat.Size() == 0 {
		l.offsets[topic][partition] = 0
	} else {
		offset, err := l.recoverOffset(f)
		if err != nil {
			return nil, err
		}
		l.offsets[topic][partition] = offset
	}
	l.files[topic][partition] = f
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

func (l *Log) Append(topic string, partition int, msg *Message) (int64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := l.getFile(topic, partition)
	if err != nil {
		return 0, err
	}
	msg.Offset = l.offsets[topic][partition]
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
	l.offsets[topic][partition]++
	return msg.Offset, nil
}

func (l *Log) Read(topic string, partition int, offset int64, max int) ([]*Message, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := l.getFile(topic, partition)
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
	for _, partMap := range l.files {
		for _, f := range partMap {
			f.Close()
		}
	}
	return nil
}
