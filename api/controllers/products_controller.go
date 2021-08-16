package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/google/uuid"
	"strings"
	"github.com/joho/godotenv"
	"os"
	"io"

	"github.com/gorilla/mux"
	"../auth"
	"../models"
	"../responses"
	"../utils/formaterror"
)

func (server *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {

	/*
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	product := models.Product{}
	err = json.Unmarshal(body, &product)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	*/

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	err = godotenv.Load()
	// <num> limit filesize <num> in MB
	limit_size, err := strconv.ParseInt(os.Getenv("IMG_LIMIT"),10,32)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	r.ParseMultipartForm(limit_size)
	product := models.Product{}

	product.Code = r.Form.Get("code")
	product.Name = r.Form.Get("name")
	product.UserID = uid
	product.Description = r.Form.Get("description")
	res, err := strconv.ParseInt(r.Form.Get("default_price"),10,32)

	if err != nil {
		fmt.Printf(r.Form.Get("default_price"))
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	product.DefaultPrice = uint32(res)

	file, handler, err := r.FormFile("myImage")
	imageName := ""

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)

	if err == nil {
		imageName = "product_"+strings.Replace(uuid.New().String(), "-", "", -1) + ".png"
		
		f, err := os.OpenFile(os.Getenv("IMG_DIR")+"/images/products/"+imageName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		io.Copy(f, file)
	}
	
	product.Image = imageName
	defer file.Close()

	product.Prepare()
	err = product.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	
	if uid != product.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	productCreated, err := product.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, productCreated.ID))
	responses.JSON(w, http.StatusCreated, productCreated)
}

func (server *Server) GetProducts(w http.ResponseWriter, r *http.Request) {

	product := models.Product{}

	products, err := product.FindAllProducts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, products)
}

func (server *Server) GetProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	product := models.Product{}

	productReceived, err := product.FindPostByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, productReceived)
}

func (server *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	product := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", pid).Take(&product).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// If a user attempt to update a post not belonging to him
	if uid != product.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	/*
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	productUpdate := models.Product{}
	err = json.Unmarshal(body, &productUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	*/

	err = godotenv.Load()
	// <num> limit filesize <num> in MB
	limit_size, err := strconv.ParseInt(os.Getenv("IMG_LIMIT"),10,32)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	r.ParseMultipartForm(limit_size)
	productUpdate := models.Product{}

	productUpdate.Code = r.Form.Get("code")
	productUpdate.Name = r.Form.Get("name")
	productUpdate.UserID = uid
	productUpdate.Description = r.Form.Get("description")
	res, err := strconv.ParseInt(r.Form.Get("default_price"),10,32)

	if err != nil {
		fmt.Printf(r.Form.Get("default_price"))
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productUpdate.DefaultPrice = uint32(res)

	// when flag set to true if user want remove or change image of product
	if r.Form.Get("change_image") == "true"{

		file, handler, err := r.FormFile("myImage")
		imageName := ""

		fmt.Printf("Uploaded File: %+v\n", handler.Filename)

		if err == nil {
			imageName = "product_"+strings.Replace(uuid.New().String(), "-", "", -1) + ".png"

			f, err := os.OpenFile(os.Getenv("IMG_DIR")+"/images/products/"+imageName, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			io.Copy(f, file)

			err = os.Remove(os.Getenv("IMG_DIR")+"/images/products/"+product.Image)

			if err != nil {
				fmt.Println(err)
			}
		}
		
		productUpdate.Image = imageName
		defer file.Close()
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != productUpdate.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	productUpdate.Prepare()
	err = productUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productUpdate.ID = product.ID //this is important to tell the model the post id to update, the other update field are set above

	productUpdated, err := productUpdate.UpdateAPost(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, productUpdated)
}

func (server *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	product := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", pid).Take(&product).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	filename_image := product.Image

	// Is the authenticated user, the owner of this post?
	if uid != product.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	_, err = product.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// after remove post from database, remove image from disk
	err = os.Remove(os.Getenv("IMG_DIR")+"/images/products/"+filename_image)

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}