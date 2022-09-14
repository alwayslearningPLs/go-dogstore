package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	breedsPath      = "/breeds"
	breedByNamePath = breedsPath + "/:" + breedParamName

	ownersPath    = "/owners"
	ownerByIdPath = ownersPath + "/:" + ownerParamName

	dogsPath      = "/dogs"
	dogByNamePath = dogsPath + "/:" + dogParamName

	dogsByBreedPath      = breedByNamePath + dogsPath
	dogByBreedByNamePath = dogsByBreedPath + "/:" + dogParamName

	dogsByOwnerPath      = ownerByIdPath + dogsPath
	dogByOwnerByNamePath = dogsByOwnerPath + "/:" + dogParamName

	breedParamName = "breed-name"
	ownerParamName = "owner-id"
	dogParamName   = "dog-name"

	ownerQueryName = "owner"
	breedQueryName = "breed"
)

var (
	MimeOffered = []string{gin.MIMEJSON}

	errOwnerNotSpecified = errors.New("owner not specified")
	errBreedNotSpecified = errors.New("breed not specified")
)

func hAddBreed(c *gin.Context) {
	var b breed

	if err := c.BindJSON(&b); err != nil {
		errResp(c, err, http.StatusBadRequest)
		return
	}

	if err := addBreed(c.Request.Context(), &b); err != nil {
		// TODO fix err so add, get, del custom error is returned
		errResp(c, err, http.StatusInternalServerError)
		return
	}

	c.Negotiate(http.StatusCreated, gin.Negotiate{
		Offered: MimeOffered,
		Data:    b,
	})
}

func hGetBreed(c *gin.Context) {
	r, err := getBreed(c.Request.Context(), breed{Name: c.Param(breedParamName)})
	if err != nil {
		errResp(c, err, http.StatusNotFound)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: MimeOffered,
		Data:    r,
	})
}

func hDelBreed(c *gin.Context) {
	row, err := delBreed(c.Request.Context(), breed{Name: c.Param(breedParamName)})
	if err != nil {
		errResp(c, err, http.StatusNotFound)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: MimeOffered,
		Data:    ResponseWrap{Code: http.StatusOK, Msg: strconv.Itoa(int(row)) + " rows deleted successfully"},
	})
}

func hAddOwner(c *gin.Context) {
	var o owner

	if err := c.BindJSON(&o); err != nil {
		errResp(c, err, http.StatusBadRequest)
		return
	}

	if err := addOwner(c.Request.Context(), &o); err != nil {
		errResp(c, err, http.StatusInternalServerError)
		return
	}

	c.Negotiate(http.StatusCreated, gin.Negotiate{
		Offered: MimeOffered,
		Data:    o,
	})
}

func hGetOwner(c *gin.Context) {
	r, err := getOwner(c.Request.Context(), owner{ID: c.Param(ownerParamName)})
	if err != nil {
		errResp(c, err, http.StatusBadRequest)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: MimeOffered,
		Data:    r,
	})
}

func hDelOwner(c *gin.Context) {
	row, err := delOwner(c.Request.Context(), owner{ID: c.Param(ownerParamName)})
	if err != nil {
		errResp(c, err, http.StatusNotFound)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: MimeOffered,
		Data:    ResponseWrap{Code: http.StatusOK, Msg: strconv.Itoa(int(row)) + " rows deleted successfully"},
	})
}

func hAddDog(c *gin.Context) {
	var p dog

	if err := c.BindJSON(&p); err != nil {
		errResp(c, err, http.StatusBadRequest)
		return
	}

	if !readParamQueryName(c, &p.Owner.ID, ownerParamName, ownerQueryName) {
		errResp(c, errOwnerNotSpecified, http.StatusBadRequest)
		return
	}

	if !readParamQueryName(c, &p.Breed.Name, breedParamName, breedQueryName) {
		errResp(c, errBreedNotSpecified, http.StatusBadRequest)
		return
	}

	if err := addDog(c.Request.Context(), &p); err != nil {
		errResp(c, err, http.StatusInternalServerError)
		return
	}

	c.Negotiate(http.StatusCreated, gin.Negotiate{
		Offered: MimeOffered,
		Data:    p,
	})
}

func hGetDog(c *gin.Context) {
	var p dog

	p.Name = c.Param(dogParamName)
	readParamQueryName(c, &p.Owner.ID, ownerParamName, ownerQueryName)
	readParamQueryName(c, &p.Breed.Name, breedParamName, breedQueryName)

	r, err := getDog(c.Request.Context(), p)
	if err != nil {
		errResp(c, err, http.StatusNotFound)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: MimeOffered,
		Data:    r,
	})
}

func hDelDog(c *gin.Context) {
	var p dog

	p.Name = c.Param(dogParamName)
	readParamQueryName(c, &p.Owner.ID, ownerParamName, ownerQueryName)
	readParamQueryName(c, &p.Breed.Name, breedParamName, breedQueryName)

	r, err := delDog(c.Request.Context(), p)
	if err != nil {
		errResp(c, err, http.StatusNotFound)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: MimeOffered,
		Data:    strconv.Itoa(int(r)) + " rows deleted successfully",
	})
}

func errResp(c *gin.Context, err error, status int) {
	c.Negotiate(status, gin.Negotiate{
		Offered: MimeOffered,
		Data: ResponseWrap{
			Code: status,
			Msg:  err.Error(),
		},
	})
}

func readParamQueryName(c *gin.Context, target *string, paramName, queryName string) bool {
	if *target = c.Param(paramName); *target == "" {
		if *target = c.Query(queryName); *target == "" {
			return false
		}
	}
	return true
}
