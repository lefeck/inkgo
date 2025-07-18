package main

import (
	"app/config"
	"app/database"
	"app/model"
	"app/repository"
)

func initDB() (repos repository.Repository, err error) {
	cfg, err := config.Parse("config/app.yaml")
	if err != nil {
		panic("Failed to parse config: " + err.Error())
	}
	mysqldb, err := database.NewMysql(&cfg.DB)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	redisdb, err := database.NewRedis(&cfg.Redis)
	if err != nil {
		panic("Failed to connect to redis: " + err.Error())
	}

	repos = repository.NewRepository(mysqldb, redisdb)

	if cfg.DB.Migrate {
		if err := repos.Migrate(); err != nil {
			panic("Failed to migrate database: " + err.Error())
		}
	}

	return repos, nil
	// Initialize other components like server, controllers, etc.
}

func tag() {
	repos, err := initDB()
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// create tag
	tag, err := repos.Tag().List()
	if err != nil {
		panic("Failed to list tags: " + err.Error())
	}
	for _, t := range tag {
		println("Tag ID:", t.ID, "Name:", t.Name)
	}

	//// delete tag
	//err = repos.Tag().Delete(27) // Assuming 1 is the ID of the tag to delete
	//if err != nil {
	//	panic("Failed to delete tag: " + err.Error())
	//} else {
	//	println("Tag deleted successfully")
	//}

	//// create new tag
	//newTag := &model.Tag{
	//	Name: "Rust",
	//}
	//createdTag, err := repos.Tag().Create(newTag)
	//if err != nil {
	//	panic("Failed to create tag: " + err.Error())
	//} else {
	//	println("Tag created successfully with ID:", createdTag.ID, "Name:", createdTag.Name)
	//}

	// update tag
	updatedTag := &model.Tag{
		ID:   42,
		Name: "Rust Programming",
	}
	updated, err := repos.Tag().Update(updatedTag)
	if err != nil {
		panic("Failed to update tag: " + err.Error())
	} else {
		println("Tag updated successfully with ID:", updated.ID, "Name:", updated.Name)
	}
}

func main() {
	tag()
}
