package backend

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetServiceList godoc
// @Summary Get a list of filtered languages and related project information
// @Description Retrieves a list of languages filtered by the specified query and details of the user's most recent draft project, including its ID and count.
// @Tags Services
// @Accept json
// @Produce json
// @Param langname query string false "Language name to filter the list of services"
// @Success 200 {object} gin.H "List of filtered languages, draft project ID, and project count"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Draft project not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /info [get]
func (app *App) GetServiceList(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

	query := c.Query("langname")
	filteredLangs, err := app.GetFilteredLangs(query)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to retrieve filtered languages"), err)
		return
	}

	var count int64 = 0
	var projectID uint = 0
	if requestUserID != uuid.Nil {
		projectID, err = findLastDraft(app, requestUserID)
		if err != nil {
			handleError(c, http.StatusNotFound, errors.New("[err] unable to find the last draft project for the user"), err)
			return
		}

		count, err = app.getProjectCount(projectID)
		if err != nil {
			handleError(c, http.StatusInternalServerError, errors.New("[err] unable to retrieve project count for the draft"), err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"langs":   filteredLangs,
		// "draftID": projectID,
		"count":   count,
	})
}

// GetServiceByID godoc
// @Summary Get details of a language by its ID
// @Description Retrieves the details of a language based on the provided language ID.
// @Tags Services
// @Accept json
// @Produce json
// @Param id path int true "Language ID"
// @Success 200 {object} gin.H "Language details"
// @Failure 400 {object} ErrorResponse "Invalid language ID format"
// @Failure 404 {object} ErrorResponse "Language not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /info/{id} [get]
func (app *App) GetServiceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid language ID format"), err)
		return
	}

	lang, err := app.getLangFirst(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] language not found for the given ID"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"info": lang,
	})
}

type CreateServiceRequest struct {
	Name             string            `json:"name" example:"Example Service"`
	ShortDescription string            `json:"short_description" example:"This is a short description of the service"`
	Description      string            `json:"description" example:"This is a detailed description of the service and its features"`
	Author           *string           `json:"author" example:"John Doe"`
	Year             *string           `json:"year" example:"2024"`
	Version          *string           `json:"version" example:"1.0.0"`
	List             map[string]string `json:"list" example:"item1:10,item2:20,item3:30,item4:40"`
}

// CreateService godoc
// @Summary Create a new language service
// @Description Creates a new language service with the provided details and saves it to the database.
// @Tags Services
// @Accept json
// @Produce json
// @Param body body CreateServiceRequest true "Service creation details"
// @Success 201 {object} gin.H "Service successfully created"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /info [post]
func (app *App) CreateService(c *gin.Context) {
	var input CreateServiceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid input data: unable to parse JSON"), err)
		return
	}

	service := DbLang{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Author:           input.Author,
		Year:             input.Year,
		Version:          input.Version,
		List:             input.List,
	}

	langID, err := app.createLang(&service)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to save service in the database"), err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Service successfully created.",
		"serviceID": langID,
	})
}

type UpdateServiceRequest struct {
	Name             string            `json:"name" example:"Updated Service Name"`
	ShortDescription string            `json:"short_description" example:"Updated short description"`
	Description      string            `json:"description" example:"Updated detailed description of the service"`
	Author           *string           `json:"author" example:"Jane Smith"`
	Year             *string           `json:"year" example:"2025"`
	Version          *string           `json:"version" example:"2.0.0"`
	List             map[string]string `json:"list" example:"item1:15,item2:25,item3:35,item4:45"`
}

// UpdateService godoc
// @Summary Update an existing language service
// @Description Updates the details of an existing language service based on the provided ID and request data.
// @Tags Services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param body body UpdateServiceRequest true "Service update details"
// @Success 200 {object} gin.H "Service successfully updated"
// @Failure 400 {object} ErrorResponse "Invalid input data"
// @Failure 404 {object} ErrorResponse "Service not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /info/{id} [put]
func (app *App) UpdateService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID format"), err)
		return
	}

	var input UpdateServiceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid input data"), err)
		return
	}

	service, err := app.getLangFirst(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	if input.Name != "" {
		service.Name = input.Name
	}
	if input.ShortDescription != "" {
		service.ShortDescription = input.ShortDescription
	}
	if input.Description != "" {
		service.Description = input.Description
	}
	if input.Author != nil {
		service.Author = input.Author
	}
	if input.Year != nil {
		service.Year = input.Year
	}
	if input.Version != nil {
		service.Version = input.Version
	}
	if len(input.List) != 0 {
		service.List = input.List
	}

	if err := app.updateLang(&service); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update service"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service updated successfully",
		"service": service,
	})
}

// UpdateServiceImage godoc
// @Summary Update the image of an existing language service
// @Description Updates the image of an existing language service identified by the provided service ID.
// @Tags Services
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Service ID"
// @Param image formData file true "Service image file"
// @Success 200 {object} gin.H "Service image updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid service ID or missing image file"
// @Failure 404 {object} ErrorResponse "Service not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /info/{id} [post]
func (app *App) UpdateServiceImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID"), err)
		return
	}

	service, err := app.getLangFirst(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] image file is required"), err)
		return
	}

	imageURL, err := app.uploadImageToMinIO(file)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to upload image to MinIO"), err)
		return
	}

	service.ImgLink = &imageURL
	service.Status = false

	if err := app.updateLang(&service); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update service image path"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service image updated successfully",
		"service": service,
	})
}

// DeleteService godoc
// @Summary Delete a language service by ID
// @Description Deletes a language service and its associated image from MinIO. If the service is not found, an error will be returned.
// @Tags Services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} gin.H "Service deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid service ID format"
// @Failure 404 {object} ErrorResponse "Service not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /info/{id} [delete]
func (app *App) DeleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid service ID format"), err)
		return
	}

	service, err := app.getLangFirst(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] service not found"), err)
		return
	}

	if err := app.deleteImageFromMinIO(*service.ImgLink); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to delete image from MinIO"), err)
		return
	}

	if err := app.deleteLang(uint(id)); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to delete service from database"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service deleted successfully",
		"status":  true,
	})
}

type AddServiceRequest struct {
	IDLang uint `json:"id_lang" example:"3"`
}

// AddServiceToDraft godoc
// @Summary Add a service to a project draft
// @Description Adds a service to a draft project. This endpoint expects a service ID and creates a new project for the user, adding the specified service to the draft.
// @Tags Services
// @Accept json
// @Produce json
// @Param request body AddServiceRequest true "Service request data"
// @Success 200 {object} gin.H "Service successfully added to draft"
// @Failure 400 {object} ErrorResponse "Invalid request format or missing fields"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /info/add-service [post]
func (app *App) AddServiceToDraft(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

	var req AddServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid JSON format or missing fields"), err)
		return
	}

	projectID, err := app.createProject(requestUserID)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] unable to create project"), err)
		return
	}

	if err := app.createFile(projectID, req.IDLang); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to create file for service draft"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Service successfully added to draft",
		"projectID": projectID,
	})
}
