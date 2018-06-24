package handler

import (
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/xesina/golang-echo-realworld-example-app/db"
	"github.com/xesina/golang-echo-realworld-example-app/store"
	"encoding/json"
	"github.com/xesina/golang-echo-realworld-example-app/model"
	"github.com/xesina/golang-echo-realworld-example-app/user"
	"github.com/xesina/golang-echo-realworld-example-app/article"
)

var (
	d  *gorm.DB
	us user.Store
	as article.Store
	h  *Handler
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func authHeader(token string) string {
	return "Token " + token
}

func setup() {
	d = db.TestDB()
	db.AutoMigrate(d)
	us = store.NewUserStore(d)
	as = store.NewArticleStore(d)
	h = NewHandler(us, as)
	loadFixtures()
}

func tearDown() {
	_ = d.Close()
	if err := db.DropTestDB(); err != nil {
		log.Fatal(err)
	}
}

func responseMap(b []byte, key string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m[key].(map[string]interface{})
}

func loadFixtures() error {
	bio := "user1 bio"
	image := "http://realworld.io/user1.jpg"
	u1 := model.User{
		Username: "user1",
		Email:    "user1@realworld.io",
		Bio:      &bio,
		Image:    &image,
	}
	u1.Password, _ = u1.HashPassword("secret")
	if err := us.Create(&u1); err != nil {
		return err
	}

	bio = "user2 bio"
	image = "http://realworld.io/user2.jpg"
	u2 := model.User{
		Username: "user2",
		Email:    "user2@realworld.io",
		Bio:      &bio,
		Image:    &image,
	}
	u2.Password, _ = u2.HashPassword("secret")
	if err := us.Create(&u2); err != nil {
		return err
	}
	us.AddFollower(&u2, u1.ID)

	//bio = "user2 bio"
	//image = "http://realworld.io/user2.jpg"
	//u = models.User{
	//	Username: "user2",
	//	Email:    "user2@realworld.io",
	//	Bio:      &bio,
	//	Image:    &image,
	//}
	//_ = u.HashPassword("secret")
	//if err := db.Create(&u).Error; err != nil {
	//	return err
	//}
	//db.Model(&u).Association("Followings").Replace(models.Follow{FollowerID: 2, FollowingID: 1})
	//
	//a := models.Article{
	//	Slug:        "article1-slug",
	//	Title:       "article1 title",
	//	Description: "article1 description",
	//	Body:        "article1 body",
	//	AuthorID:    1,
	//	Comments: []models.Comment{
	//		{
	//			UserID:    1,
	//			Body:      "article1 comment1",
	//			ArticleID: 1,
	//		},
	//	},
	//	Favorites: []models.User{
	//		u,
	//	},
	//	Tags: []models.Tag{
	//		{
	//			Tag: "tag1",
	//		},
	//		{
	//			Tag: "tag2",
	//		},
	//	},
	//}
	//if err := db.Create(&a).Error; err != nil {
	//	return err
	//}
	return nil
}
