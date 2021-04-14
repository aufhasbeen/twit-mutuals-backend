package models

type jsonUser struct {
	// json loaded
	name                    string
	screen_name             string
	profile_image_url_https string
	description             string

	// generated
	link string
}
