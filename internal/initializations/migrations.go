package initializations

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"yt-clone-video-processing/internal/dependency"
)

func RunMigrations(dependency *dependency.Dependency) {
	driver, err := postgres.WithInstance(dependency.DBConn, &postgres.Config{})
	if err != nil {
		log.Fatalln("Error while creating postgres instance for migrations", err)
	}

	instance, err := migrate.NewWithDatabaseInstance("file://./database/migrations", "postgres", driver)
	if err != nil {
		log.Fatalln("Error while creating db instance for migrations", err)
	}

	err = instance.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalln("Error while running db migrations", err)
		} else {
			log.Println("No changes to apply")
		}
	}

	log.Println("Finished running migrations")
}
