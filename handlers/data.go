package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

var AWS_S3_ROOT = "https://s3.amazonaws.com/monitoring.dadrix.s3/"

type dataMiddlewareInterface interface {
	GetData(ctx context.Context, arg storage.GetDataParams) (*models.Data, error)
}

func (handler *AppHandler) DataMiddleware(mux chi.Router, db dataMiddlewareInterface) {
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			dataIdParam := chi.URLParamFromCtx(ctx, "dataId")
			dataId, err := primitive.ObjectIDFromHex(dataIdParam)
			if err != nil {
				http.Error(w, "ERR_DATA_MDW_01", http.StatusBadRequest)
				return
			}

			activity := ctx.Value("activity").(*models.Activity)

			data, err := db.GetData(ctx, storage.GetDataParams{
				Id:         dataId,
				ActivityId: activity.Id,
			})
			if err != nil {
				http.Error(w, "ERR_DATA_MDW_02", http.StatusBadRequest)
				return
			}
			if data == nil {
				http.Error(w, "ERR_DATA_MDW_03", http.StatusNotFound)
				return
			}

			ctx = context.WithValue(ctx, "data", data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

type createDataInterface interface {
	CreateData(ctx context.Context, arg storage.CreateDataParams) (*models.Data, error)
	GetDataFromValues(ctx context.Context, arg storage.GetDataFromValuesParams) (*models.Data, error)
}

type CreateDataRequest struct {
	Values map[string]any `json:"values,omitempty"`
}

type CreateDataResponse struct {
	Data models.Data `json:"data"`
}

func (handler *AppHandler) CreateData(mux chi.Router, db createDataInterface) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser := handler.GetAuthenticatedUser(r)

		var input CreateDataRequest
		httpStatus, err := handler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		activity := ctx.Value("activity").(*models.Activity)

		// PRIMARY KEY CHECKING
		var primaryKeyField models.ActivityField
		for i := range activity.Fields {
			field := activity.Fields[i]
			if field.PrimaryKey {
				primaryKeyField = field
			}
		}

		primaryKeyValue := input.Values[primaryKeyField.Id.Hex()]
		if primaryKeyValue == nil {
			http.Error(w, "ERR_DATA_CRT_PRKEY_NOT_FOUND", http.StatusBadRequest)
			return
		}

		dataWithPrKey, err := db.GetDataFromValues(ctx, storage.GetDataFromValuesParams{
			Values: map[string]any{
				primaryKeyField.Id.Hex(): primaryKeyValue,
			},
			ActivityId: activity.Id,
		})
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		if dataWithPrKey != nil {
			http.Error(w, "ERR_DATA_CRT_PRKEY_ALREADY_USED", http.StatusBadRequest)
			return
		}

		// We should ensure that all the data are the type of the one defined in activity
		// values := make(map[string]any)
		// for _, field := range activity.Fields {
		// 	value := input.Values[field.Code]
		// 	castValue, ok := field.IsValid(value)
		// 	if ok {
		// 		values[field.Code] = castValue
		// 	}
		// }
		// for code, value := range input.Values {
		// 	var field models.ActivityFields

		// 	for _, f := range activity.Fields {

		// 	}
		// }

		data, err := db.CreateData(ctx, storage.CreateDataParams{
			Values: input.Values,

			ActivityId: activity.Id,
			CreatedBy: models.DataAuthor{
				Id:   authUser.Id,
				Name: fmt.Sprintf("%s %s", authUser.LastName, authUser.FirstName),
			},
		})
		if err != nil {
			http.Error(w, "ERR_DATA_CRT_01", http.StatusBadRequest)
			return
		}

		response := CreateDataResponse{
			Data: *data,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_CRT_END", http.StatusBadRequest)
			return
		}
	})
}

type updateDataInterface interface {
	UpdateData(ctx context.Context, arg storage.UpdateDataParams) (*models.Data, error)
	GetDataFromValues(ctx context.Context, arg storage.GetDataFromValuesParams) (*models.Data, error)
}

type UpdateDataRequest struct {
	Values map[string]any `json:"values"`
}

type UpdateDataResponse struct {
	Data models.Data `json:"data"`
}

func (appHandler *AppHandler) UpdateData(mux chi.Router, db updateDataInterface) {
	mux.Put("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// authUser := appHandler.GetAuthenticatedUser(r)

		var input UpdateDataRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		activity := ctx.Value("activity").(*models.Activity)
		data := ctx.Value("data").(*models.Data)

		// PRIMARY KEY CHECKING
		var primaryKeyField models.ActivityField
		for i := range activity.Fields {
			field := activity.Fields[i]
			if field.PrimaryKey {
				primaryKeyField = field
			}
		}

		primaryKeyValue := input.Values[primaryKeyField.Id.Hex()]
		if primaryKeyValue == nil {
			http.Error(w, "ERR_DATA_UPDT_PRKEY_NOT_FOUND", http.StatusBadRequest)
			return
		}

		dataWithPrKey, err := db.GetDataFromValues(ctx, storage.GetDataFromValuesParams{
			Values: map[string]any{
				primaryKeyField.Id.Hex(): primaryKeyValue,
			},
			ActivityId: activity.Id,
		})
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		if dataWithPrKey != nil {
			http.Error(w, "ERR_DATA_UPDT_PRKEY_ALREADY_USED", http.StatusBadRequest)
			return
		}

		data, err = db.UpdateData(ctx, storage.UpdateDataParams{
			Id:         data.Id,
			ActivityId: activity.Id,

			Values: input.Values,
		})
		if err != nil {
			http.Error(w, "ERR_DATA_UPDT_FAILED", http.StatusBadRequest)
			return
		}

		response := UpdateDataResponse{
			Data: *data,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_UPDT_ENC_RESP", http.StatusBadRequest)
			return
		}
	})
}

type getAllDataInterface interface {
	GetAllData(ctx context.Context, arg storage.GetAllDataParams) ([]*models.Data, error)
}

type FieldResponse struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type GetAllDataResponse struct {
	Fields map[string]FieldResponse `json:"fields"`
	Data   []*models.Data           `json:"data"`
}

func (handler *AppHandler) GetAllData(mux chi.Router, db getAllDataInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		activity := ctx.Value("activity").(*models.Activity)

		data, err := db.GetAllData(ctx, storage.GetAllDataParams{
			ActivityId: activity.Id,
		})
		if err != nil {
			http.Error(w, "ERR_DATA_GALL_01", http.StatusBadRequest)
			return
		}

		fields := make(map[string]FieldResponse)
		for _, field := range activity.Fields {
			fields[field.Id.Hex()] = FieldResponse{
				Name: field.Name,
				Type: field.Type,
			}
		}

		response := GetAllDataResponse{
			Fields: fields,
			Data:   data,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_GALL_END", http.StatusBadRequest)
			return
		}
	})
}

type getDataInterface interface {
}

type GetDataResponse struct {
	Data models.Data
}

func (handler *AppHandler) GetData(mux chi.Router, db getDataInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data := ctx.Value("data").(*models.Data)

		response := GetDataResponse{
			Data: *data,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_GONE_END", http.StatusBadRequest)
			return
		}
	})
}

type deleteDataInterface interface {
	DeleteData(ctx context.Context, arg storage.DeleteDataParams) error
}

type deleteFilesS3Interface interface {
	DeleteFile(uploadKey string) error
}

type DeleteDataResponse struct {
	Deleted bool `json:"deleted"`
}

func (handler *AppHandler) DeleteData(mux chi.Router, db deleteDataInterface, s3 deleteFilesS3Interface) {
	mux.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		activity := ctx.Value("activity").(*models.Activity)
		data := ctx.Value("data").(*models.Data)

		uploadFields := make([]string, 0)
		for i := range activity.Fields {
			field := activity.Fields[i]
			if field.Type == "upload" {
				uploadFields = append(uploadFields, field.Id.Hex())
			}
		}
		if len(uploadFields) > 0 {
			deleteUploadFields(uploadFields, data, s3)
		}

		err := db.DeleteData(ctx, storage.DeleteDataParams{
			Id:         data.Id,
			ActivityId: activity.Id,
		})
		if err != nil {
			http.Error(w, "ERR_DATA_DLT_01", http.StatusBadRequest)
			return
		}

		response := DeleteDataResponse{
			Deleted: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_DLT_END", http.StatusBadRequest)
			return
		}
	})
}

func deleteUploadFields(uploadField []string, data *models.Data, s3 deleteFilesS3Interface) {
	var wg sync.WaitGroup

	s3Items := make([]string, 0)
	for i := range uploadField {
		var fieldValues []string

		fieldId := uploadField[i]
		_strValues := fmt.Sprintf("%s", data.Values[fieldId])
		json.Unmarshal([]byte(_strValues), &fieldValues)
		// if err != nil {}

		s3Items = append(s3Items, fieldValues...)
	}

	wg.Add(len(s3Items))
	for i := range s3Items {
		fileKey := strings.TrimPrefix(s3Items[i], AWS_S3_ROOT)

		go func(key string) {
			defer wg.Done()

			s3.DeleteFile(fileKey)
		}(fileKey)
	}
	wg.Wait()

}

var (
	_, b, _, _ = runtime.Caller(0)
	RootPath   = filepath.Join(filepath.Dir(b), "../..")
)

func fileNameWithoutExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func saveFile(file multipart.File, handler *multipart.FileHeader) (*os.File, error) {
	//2. Retrieve file from form-data
	//<Form-id> is the form key that we will read from. Client should use the same form key when uploading the file
	defer file.Close()

	//3. Create a temporary file to our directory
	tempFolderPath := fmt.Sprintf("%s%s", RootPath, "/tmp-files")
	tempFileName := fmt.Sprintf("upload-%s-*%s", fileNameWithoutExtension(handler.Filename), filepath.Ext(handler.Filename))
	tempFile, err := os.CreateTemp(tempFolderPath, tempFileName)
	if err != nil {
		errStr := fmt.Errorf("error in creating the file %s", err)
		fmt.Println(errStr)
		return nil, errStr
	}

	// defer tempFile.Close()

	//4. Write upload file bytes to your new file
	filebytes, err := io.ReadAll(file)
	if err != nil {
		errStr := fmt.Errorf("error in reading the file buffer %s", err)
		fmt.Println(errStr)
		return nil, errStr
	}

	tempFile.Write(filebytes)
	return tempFile, nil
}

type uploadFilesDBInterface interface {
	DeleteData(ctx context.Context, arg storage.DeleteDataParams) error
	AddUploadedFile(ctx context.Context, arg storage.AddUploadedFileParams) (*models.UploadedFile, error)
}

type uploadFilesStorageInterface interface {
	UploadFile(uploadKey string, fileToUpload *os.File) error
}

type UploadFilesResponse struct {
	FileKey string `json:"file_key"`
}

func (appHandler *AppHandler) UploadFiles(mux chi.Router, db uploadFilesDBInterface, s3 uploadFilesStorageInterface) {
	mux.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser := appHandler.GetAuthenticatedUser(r)
		activity := ctx.Value("activity").(*models.Activity)

		// The argument to ParseMultipartForm is the max memory size (in bytes)
		// that will be used to store the file in memory.
		r.ParseMultipartForm(200 << 20) // 200 MB

		file, handler, err := r.FormFile("uploaded-file")
		if err != nil {
			errStr := fmt.Sprintf("Error in reading the file %s\n", err)
			fmt.Println(errStr)
			http.Error(w, "ERR_DATA_UPLF_01", http.StatusBadRequest)
			return
		}

		fileToUpload, err := saveFile(file, handler)
		if err != nil {
			log.Println(err)
			http.Error(w, "ERR_DATA_UPLF_02", http.StatusBadRequest)
			return
		}

		fileKey := fmt.Sprintf("data/%d-%s", time.Now().Unix(), handler.Filename)
		err = s3.UploadFile(
			fileKey,
			fileToUpload,
		)
		if err != nil {
			http.Error(w, "ERR_DATA_UPLF_03", http.StatusBadRequest)
			return
		}

		fullFilePath := fmt.Sprintf("%s%s", AWS_S3_ROOT, fileKey)
		uploadedFile, err := db.AddUploadedFile(ctx, storage.AddUploadedFileParams{
			UploadedBy: authUser.Id,
			ActivityId: activity.Id,
			FileKey:    fullFilePath,
		})
		if err != nil {
			http.Error(w, "ERR_DATA_UPLF_04", http.StatusBadRequest)
			return
		}

		response := UploadFilesResponse{
			FileKey: uploadedFile.FileKey,
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_UPLF_END", http.StatusBadRequest)
			return
		}
	})
}

type getFilesInterface interface {
	GetAllUploadedFiles(ctx context.Context, arg storage.GetAllUploadedFilesParams) ([]*models.UploadedFile, error)
}

type GetAllUploadedFilesResponse struct {
	Files []string `json:"files"`
}

func (appHandler *AppHandler) GetUploadedFiles(mux chi.Router, db getFilesInterface) {
	mux.Get("/upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser := appHandler.GetAuthenticatedUser(r)
		activity := ctx.Value("activity").(*models.Activity)

		uploadedFiles, err := db.GetAllUploadedFiles(ctx, storage.GetAllUploadedFilesParams{
			UploadedBy: authUser.Id,
			ActivityId: activity.Id,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, "ERR_DATA_GUPF_01", http.StatusBadRequest)
			return
		}

		files := make([]string, len(uploadedFiles))
		for i := range uploadedFiles {
			files = append(files, uploadedFiles[i].FileKey)
		}

		response := GetAllUploadedFilesResponse{
			Files: files,
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_GUPF_END", http.StatusBadRequest)
			return
		}
	})
}

type deleteUploadedFileDBInterface interface {
	RemoveUploadedFile(ctx context.Context, arg storage.RemoveUploadedFileParams) error
}

type deleteUploadedFileStorageInterface interface {
	DeleteFile(uploadKey string) error
}

type DeleteUploadedFileResponse struct {
	Deleted bool `json:"deleted"`
}

type DeleteUploadedFileRequest struct {
	FileKey string `json:"file_key"`
}

func (appHandler *AppHandler) DeleteUploadedFile(mux chi.Router, db deleteUploadedFileDBInterface, s3 deleteUploadedFileStorageInterface) {
	mux.Delete("/upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser := appHandler.GetAuthenticatedUser(r)
		activity := ctx.Value("activity").(*models.Activity)

		var input DeleteUploadedFileRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		fileKey := input.FileKey

		err = s3.DeleteFile(fileKey)
		if err != nil {
			http.Error(w, "ERR_DATA_DULF_01", http.StatusBadRequest)
			return
		}

		err = db.RemoveUploadedFile(ctx, storage.RemoveUploadedFileParams{
			UploadedBy: authUser.Id,
			ActivityId: activity.Id,
			FileKey:    fileKey,
		})
		if err != nil {
			http.Error(w, "ERR_DATA_DULF_02", http.StatusBadRequest)
			return
		}

		response := DeleteUploadedFileResponse{
			Deleted: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_DULF_END", http.StatusBadRequest)
			return
		}
	})

}
