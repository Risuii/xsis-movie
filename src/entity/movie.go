package entity

type Movie struct {
	ModelID
	ModelLogTime
	MovieData
}

type MovieData struct {
	Title       string  `db:"title"`
	Description string  `db:"description"`
	Rating      float32 `db:"rating"`
	Image       string  `db:"image"`
}
