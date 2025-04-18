package model

type ArticleTag struct {
	Article Article `json:"article"`
	Tag     Tags    `json:"tag"`
}

type Tags struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
