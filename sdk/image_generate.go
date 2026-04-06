package sdk

import (
	"context"
	"fmt"
)

// --- Generation ---

type imageGenerateConfig struct {
	Params ImageGenerationParams
}

// ImageGenerateOption configures an image generation request.
type ImageGenerateOption func(*imageGenerateConfig)

func WithImageGenerationModel(model *ImageGenerationModel) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.Model = model }
}

func WithImagePrompt(prompt string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.Prompt = prompt }
}

func WithImageN(n int) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.N = &n }
}

func WithImageSize(size string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.Size = size }
}

func WithImageQuality(quality string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.Quality = quality }
}

func WithImageStyle(style string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.Style = style }
}

func WithImageResponseFormat(format string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.ResponseFormat = format }
}

func WithImageBackground(background string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.Background = background }
}

func WithImageOutputFormat(format string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.OutputFormat = format }
}

func WithImageOutputCompression(compression int) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.OutputCompression = &compression }
}

func WithImageModeration(moderation string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.Moderation = moderation }
}

func WithImageUser(user string) ImageGenerateOption {
	return func(c *imageGenerateConfig) { c.Params.User = user }
}

func buildImageGenerateConfig(options []ImageGenerateOption) (*imageGenerateConfig, ImageGenerationProvider, error) {
	cfg := &imageGenerateConfig{}
	for _, opt := range options {
		opt(cfg)
	}
	if cfg.Params.Model == nil {
		return nil, nil, fmt.Errorf("twilightai: image generation model is required (use WithImageGenerationModel)")
	}
	if cfg.Params.Model.Provider == nil {
		return nil, nil, fmt.Errorf("twilightai: image generation model %q has no provider", cfg.Params.Model.ID)
	}
	if cfg.Params.Prompt == "" {
		return nil, nil, fmt.Errorf("twilightai: prompt is required (use WithImagePrompt)")
	}
	return cfg, cfg.Params.Model.Provider, nil
}

// GenerateImage generates images from a text prompt.
func (c *Client) GenerateImage(ctx context.Context, options ...ImageGenerateOption) (*ImageResult, error) {
	cfg, prov, err := buildImageGenerateConfig(options)
	if err != nil {
		return nil, err
	}
	return prov.DoGenerate(ctx, &cfg.Params)
}

// --- Edit ---

type imageEditConfig struct {
	Params ImageEditParams
}

// ImageEditOption configures an image edit request.
type ImageEditOption func(*imageEditConfig)

func WithImageEditModel(model *ImageEditModel) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.Model = model }
}

func WithEditImages(images ...ImageInput) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.Images = images }
}

func WithEditPrompt(prompt string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.Prompt = prompt }
}

func WithEditMask(mask *ImageInput) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.Mask = mask }
}

func WithEditN(n int) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.N = &n }
}

func WithEditSize(size string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.Size = size }
}

func WithEditQuality(quality string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.Quality = quality }
}

func WithEditBackground(background string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.Background = background }
}

func WithEditOutputFormat(format string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.OutputFormat = format }
}

func WithEditOutputCompression(compression int) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.OutputCompression = &compression }
}

func WithEditInputFidelity(fidelity string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.InputFidelity = fidelity }
}

func WithEditModeration(moderation string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.Moderation = moderation }
}

func WithEditResponseFormat(format string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.ResponseFormat = format }
}

func WithEditUser(user string) ImageEditOption {
	return func(c *imageEditConfig) { c.Params.User = user }
}

func buildImageEditConfig(options []ImageEditOption) (*imageEditConfig, ImageEditProvider, error) {
	cfg := &imageEditConfig{}
	for _, opt := range options {
		opt(cfg)
	}
	if cfg.Params.Model == nil {
		return nil, nil, fmt.Errorf("twilightai: image edit model is required (use WithImageEditModel)")
	}
	if cfg.Params.Model.Provider == nil {
		return nil, nil, fmt.Errorf("twilightai: image edit model %q has no provider", cfg.Params.Model.ID)
	}
	if cfg.Params.Prompt == "" {
		return nil, nil, fmt.Errorf("twilightai: prompt is required (use WithEditPrompt)")
	}
	return cfg, cfg.Params.Model.Provider, nil
}

// EditImage edits or extends images given a prompt.
func (c *Client) EditImage(ctx context.Context, options ...ImageEditOption) (*ImageResult, error) {
	cfg, prov, err := buildImageEditConfig(options)
	if err != nil {
		return nil, err
	}
	return prov.DoEdit(ctx, &cfg.Params)
}
