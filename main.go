package main

import "github.com/gin-gonic/gin"

const dns = "host=db user=postgres password=abc123. dbname=dogstore port=5432 sslmode=disable TimeZone=Europe/Madrid"

func main() {
	r := gin.Default()

	configureDB(dns)

	getInstance().AutoMigrate(new(breed), new(owner), new(dog))

	dogsGroup := r.Group("dogstore")
	{
		dogsGroup.GET(breedsPath, hGetBreed)
		dogsGroup.POST(breedsPath, hAddBreed)
		dogsGroup.GET(breedByNamePath, hGetBreed)
		dogsGroup.DELETE(breedByNamePath, hDelBreed)

		dogsGroup.GET(ownersPath, hGetOwner)
		dogsGroup.POST(ownersPath, hAddOwner)
		dogsGroup.GET(ownerByIdPath, hGetOwner)
		dogsGroup.DELETE(ownerByIdPath, hDelOwner)

		dogsGroup.GET(dogsByBreedPath, hGetDog)
		dogsGroup.POST(dogsByBreedPath, hAddDog)
		dogsGroup.GET(dogByBreedByNamePath, hGetDog)
		dogsGroup.DELETE(dogByBreedByNamePath, hDelDog)

		dogsGroup.GET(dogsByOwnerPath, hGetDog)
		dogsGroup.POST(dogsByOwnerPath, hAddDog)
		dogsGroup.GET(dogByOwnerByNamePath, hGetDog)
		dogsGroup.DELETE(dogByOwnerByNamePath, hDelDog)

		dogsGroup.GET(dogsPath, hGetDog)
		dogsGroup.GET(dogByNamePath, hGetDog)
	}

	r.Run(":8080")
}
