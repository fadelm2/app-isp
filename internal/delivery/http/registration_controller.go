package http

import (
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type RegistrationController struct {
	Log     *logrus.Logger
	UseCase *usecase.RegistrationUseCase
}

func NewRegistrationController(log *logrus.Logger, useCase *usecase.RegistrationUseCase) *RegistrationController {
	return &RegistrationController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *RegistrationController) Create(ctx *fiber.Ctx) error {
	var request model.CreateRegistrationRequest

	contentType := ctx.Get("Content-Type")
	if contentType != "application/json" {
		request.FullName = ctx.FormValue("full_name")
		request.NIK = ctx.FormValue("nik")
		request.BirthPlace = ctx.FormValue("birth_place")
		request.BirthDate = ctx.FormValue("birth_date")
		request.Gender = ctx.FormValue("gender")
		request.Email = ctx.FormValue("email")
		request.Phone = ctx.FormValue("phone")
		request.InstallationAddress = ctx.FormValue("installation_address")
		request.BillingAddress = ctx.FormValue("billing_address")
		request.PackageID = ctx.FormValue("package_id")
		request.Latitude, _ = strconv.ParseFloat(ctx.FormValue("latitude"), 64)
		request.Longitude, _ = strconv.ParseFloat(ctx.FormValue("longitude"), 64)
		request.Notes = ctx.FormValue("notes")
		request.Province = ctx.FormValue("province")
		request.City = ctx.FormValue("city")
		request.District = ctx.FormValue("district")
		request.Village = ctx.FormValue("village")

		uploadDir := "./storage/uploads"
		os.MkdirAll(uploadDir+"/ktp", os.ModePerm)
		os.MkdirAll(uploadDir+"/selfie", os.ModePerm)
		os.MkdirAll(uploadDir+"/house", os.ModePerm)
		os.MkdirAll(uploadDir+"/installation", os.ModePerm)
		os.MkdirAll(uploadDir+"/documents", os.ModePerm)

		fileKeys := []string{"ktp", "selfie", "house", "installation", "supporting_doc"}
		for _, key := range fileKeys {
			file, err := ctx.FormFile(key)
			if err == nil {
				filename := uuid.NewString() + filepath.Ext(file.Filename)
				subFolder := key
				if key == "supporting_doc" {
					subFolder = "documents"
				}
				dest := filepath.Join(uploadDir, subFolder, filename)
				if err := ctx.SaveFile(file, dest); err == nil {
					pathVal := "/storage/uploads/" + subFolder + "/" + filename
					switch key {
					case "ktp":
						request.KtpPath = pathVal
					case "selfie":
						request.SelfiePath = pathVal
					case "house":
						request.HousePath = pathVal
					case "installation":
						request.InstallationPath = pathVal
					case "supporting_doc":
						request.SupportingDocPath = pathVal
					}
				}
			}
		}
	} else {
		if err := ctx.BodyParser(&request); err != nil {
			return fiber.ErrBadRequest
		}
	}

	response, err := c.UseCase.Create(ctx.UserContext(), &request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[model.RegistrationResponse]{Data: response})
}

func (c *RegistrationController) List(ctx *fiber.Ctx) error {
	request := &model.SearchRegistrationRequest{
		Search: ctx.Query("search", ""),
		Status: ctx.Query("status", ""),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	response, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.RegistrationResponse]{
		Data:   response,
		Paging: paging,
	})
}

func (c *RegistrationController) UpdateStatus(ctx *fiber.Ctx) error {
	request := new(model.UpdateRegistrationStatusRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}
	request.ID = ctx.Params("registrationId")

	auth := middleware.GetUser(ctx)
	response, err := c.UseCase.UpdateStatus(ctx.UserContext(), request, auth.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[model.RegistrationResponse]{Data: response})
}
