package debug

// Config defines layout and sizing for debug overlay windows.
type Config struct {
	Margin         int
	ViewportWidth  int
	ViewportHeight int

	StatsWidth    int
	EntitiesWidth int
	MaxEntities   int
}

// normalize fills zero-valued fields with sensible defaults.
func (c *Config) normalize() {
	if c.Margin <= 0 {
		c.Margin = 16
	}
	if c.ViewportWidth <= 0 {
		c.ViewportWidth = 640
	}
	if c.ViewportHeight <= 0 {
		c.ViewportHeight = 360
	}
	if c.StatsWidth <= 0 {
		c.StatsWidth = 260
	}
	if c.EntitiesWidth <= 0 {
		c.EntitiesWidth = 360
	}
	if c.MaxEntities <= 0 {
		c.MaxEntities = 10
	}
}

