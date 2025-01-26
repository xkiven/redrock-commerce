package search

import (
	"database/sql"
)

// ValidateToken 验证token
func ValidateToken(token string) (bool, error) {

	db, err := sql.Open("mysql", "root:xkw510724@tcp(127.0.0.1:3306)/redrock_ecommerce?charset=utf8")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	//fmt.Println("获取到的token:", token)
	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE token=?", token).Scan(&count)
	//fmt.Println("获取到的count:", count)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}
