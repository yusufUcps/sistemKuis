package helper

import (
	"context"
	"log"
	"net/http"
	"os"
	"quiz/configs"

	"io"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

// ServiceAccount mengembalikan klien HTTP yang terautentikasi untuk penggunaan layanan Google Drive.
func ServiceAccount(email string, key string) *http.Client {
    jwtConfig := &jwt.Config{
        Email:      email,
        PrivateKey: []byte(key),
        Scopes: []string{
            drive.DriveScope,
        },
        TokenURL: google.JWTTokenURL,
    }
    client := jwtConfig.Client(context.Background())
    return client
}

// createFile membuat file baru di Google Drive.
func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
    f := &drive.File{
        MimeType: mimeType,
        Name:     name,
        Parents:  []string{parentId},
    }
    file, err := service.Files.Create(f).Media(content).Do()
    if err != nil {
        log.Println("Could not create file: " + err.Error())
        return nil, err
    }
    return file, nil
}

// getDownloadLink mengembalikan tautan untuk mengunduh file dari Google Drive.
func getDownloadLink(service *drive.Service, fileId string) (string, error) {
    file, err := service.Files.Get(fileId).Fields("webContentLink").Do()
    if err != nil {
        return "", err
    }
    return file.WebContentLink, nil
}

func UploadFile(config configs.ProgramConfig, fileName string) (string, error){
	f, err := os.Open(fileName)
	
	client := ServiceAccount(config.ClientEmail, config.PrivateKey)
	srv, err := drive.New(client)
    if err != nil {
		return "", err
    }

	folderId := config.FolderId
	file, err := createFile(srv, f.Name(), "application/octet-stream", f, folderId)

    if err != nil {
		return "", err
    }

    // Get the download link for the uploaded file
    downloadLink, err := getDownloadLink(srv, file.Id)
    if err != nil {
        return "", err
    }

	return downloadLink, nil
}
