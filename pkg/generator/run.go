package generator

func Run(svc string) error {
	if err := RunGoctl(svc); err != nil {
		return err
	}
	if err := Rearrange(svc); err != nil {
		return err
	}
	return Overlay(svc)
}
