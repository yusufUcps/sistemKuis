package helper

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"quiz/model"
	"time"

	"io"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

type ExportInterface interface {
    ServiceAccount(email string, key string) *http.Client
    CreateFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error)
    GetDownloadLink(service *drive.Service, fileId string) (string, error)
    UploadFile(fileName string) (string, error)
    ExportMyHistoryScore(history []model.MyHistoryScoreRes, userId uint) (*model.ExportRes, error)
    ExportHistoryScoreMyQuiz(history []model.HistoryScoreMyQuizRes, quiz_id uint) (*model.ExportRes, error)
    ExportHistoryAnswer(history []model.HistoryAnswersRes, user_id uint) (*model.ExportRes, error)
}

type Export struct {
	ClientEmail string
    PrivateKey string
    FolderId string
}

func NewExport(clientEmail string, privateKey string, folderId string) ExportInterface {
	return &Export{
		ClientEmail : clientEmail,
        PrivateKey : privateKey,
        FolderId : folderId,
	}
}

func (e *Export) ServiceAccount(email string, key string) *http.Client {
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


func (e *Export) CreateFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
    f := &drive.File{
        MimeType: mimeType,
        Name:     name,
        Parents:  []string{parentId},
    }
    file, err := service.Files.Create(f).Media(content).Do()
    if err != nil {
        logrus.Error("export: cannot create file, ", err.Error())
        return nil, err
    }
    return file, nil
}


func (e *Export) GetDownloadLink(service *drive.Service, fileId string) (string, error) {
    file, err := service.Files.Get(fileId).Fields("webContentLink").Do()
    if err != nil {
        logrus.Error("export: cannot connect client service, ", err.Error())
        return "", err
    }
    return file.WebContentLink, nil
}

func (e *Export) UploadFile(fileName string) (string, error){
	f, err := os.Open(fileName)
	
	client := e.ServiceAccount(e.ClientEmail, e.PrivateKey)
	srv, err := drive.New(client)
    if err != nil {
        logrus.Error("export: cannot connect client, ", err.Error())
		return "", err
    }

	folderId := e.FolderId
	file, err := e.CreateFile(srv, f.Name(), "application/octet-stream", f, folderId)

    if err != nil {
        logrus.Error("export: cannot CreateFile , ", err.Error())
		return "", err
    }

    downloadLink, err := e.GetDownloadLink(srv, file.Id)
    if err != nil {
        logrus.Error("export: cannot GetDownloadLinkfile, ", err.Error())
        return "", err
    }

	return downloadLink, nil
}

func (e *Export) ExportMyHistoryScore(history []model.MyHistoryScoreRes, userId uint) (*model.ExportRes, error) {

	fileName := fmt.Sprintf("file/History_Score_UserID_%d.csv", userId)
	file, err := os.Create(fileName)
	if err != nil {
        logrus.Error("export: cannot Create file, ", err.Error())
		return nil, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"id", "quiz_id", "title", "right_answer", "wrong_answer", "score", "finish_at"}
	writer.Write(header)

	for _, entry := range history {
		row := []string{
			fmt.Sprint(entry.Id),
			fmt.Sprint(entry.Quiz_id),
			entry.Title,
			fmt.Sprint(entry.Right_answer),
			fmt.Sprint(entry.Wrong_answer),
			fmt.Sprintf("%f", entry.Score),
			entry.Finish_at.Format(time.RFC3339), 
		}
		writer.Write(row)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
        logrus.Error("export: cannot writer file, ", err.Error())
		return nil, err
	}

    link, err := e.UploadFile(fileName)
	if err != nil {
		logrus.Error("export: cannot uploate file, ", err.Error())
        return nil, err
    }

    var res model.ExportRes

    res.Name_file = fmt.Sprintf("History_Score_UserID_%d.csv", userId)
    res.Link_download = link

	return &res, nil

}

func (e *Export) ExportHistoryScoreMyQuiz(history []model.HistoryScoreMyQuizRes, quiz_id uint) (*model.ExportRes, error) {

	fileName := fmt.Sprintf("file/History_Score_QuizID_%d.csv", quiz_id)
	file, err := os.Create(fileName)
	if err != nil {
        logrus.Error("export: cannot Create file, ", err.Error())
		return nil, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"id", "user_id", "name", "right_answer", "wrong_answer", "score", "finish_at"}
	writer.Write(header)

	for _, entry := range history {
		row := []string{
			fmt.Sprint(entry.Id),
			fmt.Sprint(entry.User_id),
			entry.Name,
			fmt.Sprint(entry.Right_answer),
			fmt.Sprint(entry.Wrong_answer),
			fmt.Sprintf("%f", entry.Score),
			entry.Finish_at.Format(time.RFC3339), 
		}
		writer.Write(row)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
        logrus.Error("export: cannot writer file, ", err.Error())
		return nil, err
	}

    link, err := e.UploadFile(fileName)
	if err != nil {
        logrus.Error("export: cannot uploate file, ", err.Error())
        return nil, err
    }

    var res model.ExportRes

    res.Name_file = fmt.Sprintf("History_Score_QuizID_%d.csv", quiz_id)
    res.Link_download = link

	return &res, nil
}

func (e *Export) ExportHistoryAnswer(history []model.HistoryAnswersRes, user_id uint) (*model.ExportRes, error) {

	fileName := fmt.Sprintf("file/History_Answers_UserID_%d.csv", user_id)
	file, err := os.Create(fileName)
	if err != nil {
        logrus.Error("export: cannot Create file, ", err.Error())
		return nil, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"id", "question_id", "question", "option_id", "answer", "is_right"}
	writer.Write(header)

	for _, entry := range history {
		row := []string{
			fmt.Sprint(entry.Id),
			fmt.Sprint(entry.Question_id),
			entry.Question,
			fmt.Sprint(entry.Option_id),
			entry.Answer,
			fmt.Sprintf("%t", entry.Is_right),
		}
		writer.Write(row)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
        logrus.Error("export: cannot writer file, ", err.Error())
		return nil, err
	}

    link, err := e.UploadFile(fileName)
	if err != nil {
        logrus.Error("export: cannot uploate file, ", err.Error())
        return nil, err
    }
    var res model.ExportRes

    res.Name_file = fmt.Sprintf("History_Answers_UserID_%d.csv", user_id)
    res.Link_download = link

	return &res, nil
}



