package backend

import (
	"errors"
	"net/http"
	"rip/lib/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProjectList godoc
// @Summary Get list of projects
// @Description Get a list of projects filtered by start date, end date, and status
// @Tags Projects
// @Accept json
// @Produce json
// @Param start_date query string false "Start Date in YYYY-MM-DD format"
// @Param end_date query string false "End Date in YYYY-MM-DD format"
// @Param status query string false "Status of the project"
// @Success 200 {array} database.Project "List of projects"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /project [get]
func (app *App) GetProjectList(c *gin.Context) {
	startDateStr := c.Query("start_date") // формат: "YYYY-MM-DD"
	endDateStr := c.Query("end_date")     // формат: "YYYY-MM-DD"
	status := c.Query("status")

	projects, err := app.filterProjects(startDateStr, endDateStr, status)
	if err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to retrieve projects"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

// GetProjectByID godoc
// @Summary Get details of a specific project by ID
// @Description Get detailed information about a project, including associated files, by project ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} database.Project "Detailed information about the project"
// @Success 200 {array} database.File "List of files associated with the project"
// @Failure 400 {object} ErrorResponse "Invalid project ID format"
// @Failure 404 {object} ErrorResponse "Project or files not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /project/{id} [get]
func (app *App) GetProjectByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID format"), err)
		return
	}

	project, err := app.getProjectFirst(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	files, err := app.getFilesForProject(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] files not found for project"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project": project,
		"files":   files,
	})
}

type UpdateProjectRequest struct {
	Status  database.Status `json:"status" example:"draft"`
	Comment *string         `json:"comment" example:"Updated project status to draft"`
}

// UpdateProject godoc
// @Summary Update an existing project
// @Description Update the status and comment of a project by its ID. The user must be the owner of the project to update it.
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body UpdateProjectRequest true "Request payload for updating project"
// @Success 200 {object} gin.H "Successfully updated project"
// @Failure 400 {object} ErrorResponse "Invalid request format or status"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Project not found or project does not belong to the user"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /project/{id} [put]
func (app *App) UpdateProject(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID"), err)
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	project, err := app.getProjectFirst(uint(id))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if requestUserID != project.UserID {
		handleError(c, http.StatusForbidden, errors.New("[err] project does not belong to the user"), err)
		return
	}

	if req.Status != "" {
		project.Status = req.Status
	} else {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid status"))
		return
	}

	if req.Comment != nil {
		project.ModeratorComment = req.Comment
	}

	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project updated successfully",
		"project": project,
	})
}

type SubmitProjectRequest struct {
	FileCodes map[uint]string `json:"file_codes" example:"6:file_code_1,7:file_code_2,8:file_code_3"`
}

// SubmitProject godoc
// @Summary Submit a project by updating files and status
// @Description Submit a project by updating associated file codes and setting its status to "Created". The user must be the owner of the project to submit it.
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body SubmitProjectRequest true "Request payload for submitting project"
// @Success 200 {object} gin.H "Successfully submitted the project"
// @Failure 400 {object} ErrorResponse "Invalid request format or project ID"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Project not found or project does not belong to the user"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /project/{id}/submit [put]
func (app *App) SubmitProject(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

	idStr := c.Param("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID format"), err)
		return
	}

	var req SubmitProjectRequest
	if err = c.ShouldBind(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
		return
	}

	project, err := app.getProjectFirst(uint(projectID))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if requestUserID != project.UserID {
		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
		return
	}

	if req.FileCodes != nil {
		if err := app.updateFilesCode(project.ID, req.FileCodes); err != nil {
			handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update files"), err)
			return
		}
	}

	project.Status = database.Created
	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project status"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project submitted successfully",
		"project": project,
	})
}

type CompleteProjectRequest struct {
	Status  database.Status `json:"status" example:"completed"`
	Comment *string         `json:"comment" example:"Project successfully completed"`
}

// CompleteProject godoc
// @Summary Complete or reject a project with a status and comment
// @Description Mark a project as completed or rejected, and provide an optional comment. The user must be the owner of the project to complete it. The project must have a formation date set to be completed.
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body CompleteProjectRequest true "Request payload for completing or rejecting a project"
// @Success 200 {object} gin.H "Successfully completed or rejected the project"
// @Failure 400 {object} ErrorResponse "Invalid request format, project status, or missing formation date"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Project not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /project/{id}/complete [put]
func (app *App) CompleteProject(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

	idStr := c.Param("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid project ID"), err)
		return
	}

	var req CompleteProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	project, err := app.getProjectFirst(uint(projectID))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if project.FormationTime == nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] project cannot be completed, formation date is missing"))
		return
	}

	project.ModeratorID = &requestUserID

	if req.Status == database.Completed || req.Status == database.Rejected {
		project.Status = req.Status
	} else {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid status"))
		return
	}

	if req.Comment != nil {
		project.ModeratorComment = req.Comment
	}

	err = app.updateAutocheck(project.ID)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed update autocheck"), err)
		return
	}

	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to complete project"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project completed successfully",
		"project": project,
	})
}

type DeleteProjectRequest struct {
	FileCodes map[uint]string `json:"file_codes" example:"6:file_code_1,7:file_code_2,8:file_code_3"`
}

// DeleteProject godoc
// @Summary Delete a project by updating its status to "deleted"
// @Description Marks a project as deleted, but only if the project has no formation date. The user must be the owner of the project to delete it. Optionally, file codes associated with the project can be updated before deletion.
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body DeleteProjectRequest true "Request payload for deleting a project"
// @Success 200 {object} gin.H "Successfully deleted the project"
// @Failure 400 {object} ErrorResponse "Invalid request format or project status, formation date exists"
// @Failure 401 {object} ErrorResponse "Unauthorized access"
// @Failure 404 {object} ErrorResponse "Project not found or project does not belong to the user"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /project/{id} [delete]
func (app *App) DeleteProject(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

	idStr := c.Param("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid id project"), err)
	}

	var req DeleteProjectRequest
	if err := c.ShouldBind(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid data format"), err)
		return
	}

	project, err := app.getProjectFirst(uint(projectID))
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] project not found"), err)
		return
	}

	if requestUserID != project.UserID {
		handleError(c, http.StatusNotFound, errors.New("[err] project does not belong to the user"), err)
		return
	}

	if project.FormationTime != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] project cannot be deleted, formation date found"))
		return
	}

	if req.FileCodes != nil {
		if err := app.updateFilesCode(project.ID, req.FileCodes); err != nil {
			handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update files"), err)
			return
		}
	}

	project.Status = database.Deleted
	if err := app.updateProject(&project); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update project"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project deleted successfully",
		"status":  true,
	})
}
