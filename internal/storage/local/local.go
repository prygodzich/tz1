package local

import (
	"context"
	"fmt"
	"sync"
	"targetads/internal/apperrs"
	"targetads/internal/logger"
)

type TypeContent int

const (
	Unknown TypeContent = iota
	Video
	Js
)

type Files struct {
	Files map[string][]byte
	mu    sync.RWMutex
}

type Repository struct {
	Content map[TypeContent]*Files
	Default map[TypeContent][]byte
}

func NewRepository() *Repository {
	files := make(map[TypeContent]*Files)
	files[Js] = &Files{
		Files: make(map[string][]byte),
	}
	files[Video] = &Files{
		Files: make(map[string][]byte),
	}
	return &Repository{
		Content: files,
		Default: make(map[TypeContent][]byte),
	}
}

func (r *Repository) Init(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Debug("init local storage")
	files := make(map[TypeContent]*Files)
	files[Js] = &Files{
		Files: make(map[string][]byte),
	}
	files[Video] = &Files{
		Files: make(map[string][]byte),
	}
	r.Content = files
	r.Default = make(map[TypeContent][]byte)
}

func (r *Repository) SetDefault(ctx context.Context, typeContent TypeContent, data []byte) {
	r.Default[typeContent] = data
	log := logger.FromContext(ctx)
	log.Debug("set default file ", typeContent)
}

func (r *Repository) Set(ctx context.Context, typeContent TypeContent, name string, data []byte) {
	r.Content[typeContent].mu.Lock()
	defer r.Content[typeContent].mu.Unlock()
	r.Content[typeContent].Files[name] = data
	log := logger.FromContext(ctx)
	log.Debug(fmt.Sprintf("set file type=%d name=%s len=%d", typeContent, name, len(r.Content[typeContent].Files)))
}

func (r *Repository) GetDefault(ctx context.Context, typeContent TypeContent) []byte {
	log := logger.FromContext(ctx)
	log.Debug("get default file ", typeContent)
	return r.Default[typeContent]
}

func (r *Repository) Get(ctx context.Context, typeContent TypeContent, name string) ([]byte, error) {
	cont := r.Content[typeContent]
	cont.mu.RLock()
	defer cont.mu.RUnlock()
	data, exists := cont.Files[name]
	if !exists {
		return nil, fmt.Errorf("not found file %s: %w", name, apperrs.ErrNotFound)
	}
	log := logger.FromContext(ctx)
	log.Debug("get file", name)
	return data, nil
}

func (r *Repository) Clear(ctx context.Context) {
	for _, cont := range r.Content {
		cont.mu.Lock()
		defer cont.mu.Unlock()
		cont.Files = make(map[string][]byte)
	}
}
