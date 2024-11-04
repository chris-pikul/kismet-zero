package llama

import (
	"errors"
	"fmt"
)

type Model struct {
	handle uintptr
	*Params
}

func (m *Model) pushParams() {
	if m.Params == nil {
		m.Params = NewParams()
	}
	m.Params.Flush()
}

func (m Model) Valid() bool {
	return m.handle != 0
}

func (m *Model) Load(path string) error {
	m.pushParams()
	m.handle = load_model(path)
	if m.handle == 0 {
		return fmt.Errorf("failed to load model at %s", path)
	}
	return nil
}

func (m *Model) Free() {
	if !m.Valid() {
		return
	}
	free_model(m.handle)
	m.handle = 0
}

func (m Model) InferSync(prompt string) (string, error) {
	if !m.Valid() {
		return "", errors.New("invalid model")
	}

	m.pushParams()

	infer_sync(m.handle, prompt)

	if did_error() {
		errStr := get_last_error()
		if errStr != "" {
			return "", errors.New(errStr)
		}
		return "", errors.New("unknown inference error")
	}
	return get_last_output(), nil
}
