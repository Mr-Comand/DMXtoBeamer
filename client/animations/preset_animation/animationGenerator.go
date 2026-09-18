package preset_animation

type AnimationGenerator interface {
	Create(config AnimationParameters) AnimationInterface
	Unload()
}
