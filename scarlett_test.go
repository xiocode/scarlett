/**
 * Author:        Tony.Shao
 * Email:         xiocode@gmail.com
 * Github:        github.com/xiocode
 * File:          scarlett_test.go
 * Description:   testing
 */

package scarlett

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type Users struct {
	ID         int64  `scarlett:"id"`
	UID        int64  `scarlett:"uid"`
	ScreenName string `scarlett:"screen_name"`
}

func TestScarlett(t *testing.T) {
	s := NewScarlett("mysql", "root:299792458@/xtimeline")

	var user Users

	err := s.Query(&user, nil, "SELECT id,uid,screen_name FROM tb_xweibo_user_info LIMIT 1;")
	if err != nil {
		t.Error(err)
	}

	t.Log(user)

}
