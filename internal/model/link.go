package model

type Link struct {
	Id   int64  `pg:"id,pk"`
	Link string `pg:"link,unique"`
	Code string `pg:"code"`
}
