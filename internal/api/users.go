package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/logger"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/service"
)

func (u *User) ToService() *service.User {
	return &service.User{
		ID:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func UserFromService(u *service.User) User {
	return User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func UsersFromService(usrs []service.User) []User {
	res := make([]User, 0, len(usrs))

	for i := range usrs {
		res = append(res, UserFromService(&usrs[i]))
	}

	return res
}
func (p *GetUserParams) ToService() *service.ListUsersParams {
	res := &service.ListUsersParams{}

	if p.Limit != nil {
		limit := uint64(*p.Limit)
		res.Limit = &limit
	}
	if p.Offset != nil {
		offset := uint64(*p.Offset)
		res.Offset = &offset
	}

	return res
}

func (ct *Controller) GetUser(c *gin.Context, params GetUserParams) {
	log := logger.GetLogger(c)

	srvcUsrs, err := ct.srvc.ListUsers(c.Request.Context(), params.ToService())
	if err != nil {
		log.WithError(err).Error("cannot list users")
		c.JSON(http.StatusInternalServerError, APIError{Message: http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, UsersFromService(srvcUsrs))
}
func (ct *Controller) PostUser(c *gin.Context) {
	log := logger.GetLogger(c)
	var usr User

	if err := c.ShouldBindJSON(&usr); err != nil {
		log.WithError(err).Info("wrong user data")
		c.JSON(http.StatusBadRequest, APIError{Message: http.StatusText(http.StatusBadRequest)})
		return
	}

	newID, err := ct.srvc.CreateUser(c.Request.Context(), usr.ToService())
	if err != nil {
		log.WithError(err).Error("cannot create user")
		c.JSON(http.StatusInternalServerError, APIError{Message: http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, CreatedItem{Id: newID})
}
func (ct *Controller) DeleteUserId(c *gin.Context, id string) {
	log := logger.GetLogger(c).WithField("user_id", id)

	err := ct.srvc.DeleteUserByID(c.Request.Context(), id)
	if err != nil {
		log.WithError(err).Error("cannot delete user")
		c.JSON(http.StatusInternalServerError, APIError{Message: http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, CreatedItem{Id: id})
}
func (ct *Controller) GetUserId(c *gin.Context, id string) {
	log := logger.GetLogger(c).WithField("user_id", id)

	srvcUsr, err := ct.srvc.GetUserByID(c.Request.Context(), id)
	if err != nil {
		log.WithError(err).Error("cannot get user")
		c.JSON(http.StatusInternalServerError, APIError{Message: http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, UserFromService(&srvcUsr))
}
func (ct *Controller) PutUserId(c *gin.Context, id string) {
	log := logger.GetLogger(c).WithField("user_id", id)

	var usr User
	if err := c.ShouldBindJSON(&usr); err != nil {
		log.WithError(err).Info("wrong user data")
		c.JSON(http.StatusBadRequest, APIError{Message: http.StatusText(http.StatusBadRequest)})
		return
	}

	srvcUsr, err := ct.srvc.UpdateUserByID(c.Request.Context(), id, usr.ToService())
	if err != nil {
		log.WithError(err).Error("cannot update user")
		c.JSON(http.StatusInternalServerError, APIError{Message: http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, UserFromService(&srvcUsr))
}
