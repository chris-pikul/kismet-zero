package exfer

import "context"

// Provider is the common interface for any external inference providers.
type Provider interface {
	Generate(opts *GenerateOptions) (GenerateResult, error)
	GenerateAsync(ctx context.Context, opts *GenerateOptions, fragmentChan chan<- GenerateFragment) error
}
