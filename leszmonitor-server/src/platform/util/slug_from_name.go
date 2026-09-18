package util

type SlugFromName struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func (d *SlugFromName) Init(name string) {
	d.Name = name
	d.Slug = SlugFromString(d.Name)
}
