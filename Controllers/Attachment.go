//Controllers/Attachment.go
package Controllers

import (
	"bytes"
	"fileshare/Models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//GetAttachments ... Get all Attachments
func GetAttachments(c *gin.Context) {
	//look for user
	username := c.Params.ByName("username")
	var user Models.User
	err := Models.GetUserByUsername(&user, username)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username not found!"})
		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	var attachment []Models.Attachment
	err2 := Models.GetAllAttachmentsFromUser(&attachment, user.Id)
	if err2 != nil {
		fmt.Println(err2.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, attachment)
	}
}

//CreateAttachment creates an attachment
func CreateAttachment(c *gin.Context) {
	username := c.Params.ByName("username")
	var user Models.User
	err := Models.GetUserByUsername(&user, username)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username not found!"})
		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	//limit uploads to 100mb
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100*1024*1024)
	//get file and save it to disk
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		fmt.Println("file upload err!" + err.Error())
		return
	}
	filename := header.Filename
	key := username + "/" + filename
	diskpath := os.Getenv("USER_FILE_PATH") + "/" + key
	//
	place := c.Params.ByName("place")
	if place != "FS" {
		if place == "" {
			place = "FS"
		} else {
			place = "s3"
		}
	}

	out, err := createFile(diskpath)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var attachment Models.Attachment
	attachment.FileName = filename
	attachment.UploadType = place
	attachment.UserId = user.Id
	attachment.Date = time.Now()
	//copy to s3 if required, then delete
	if place != "s3" {
		fmt.Println("Saved to fs " + key)
		//err := Models.CreateAttachment(&attachment)
		err := Models.UpdateAttachmentDateByFileNameAndUploadTypeOrCreate(&attachment)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, attachment)
		}
	} else {
		//for asynchronous s3 upload&file delete... maybe not so production-ready
		go func() {
			// Create a single AWS session (we can re use this if we're uploading many files)
			s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
			// if err != nil {
			// 	fmt.Println(err.Error())
			// 	return
			// }
			// Upload
			err = AddFileToS3(s, diskpath, key)
			if err != nil {
				fmt.Println("S3 upload error for " + key + " - " + err.Error())
				return
			}
			fmt.Println("Added to s3 " + key)
			err2 := os.Remove(diskpath)

			if err2 != nil {
				fmt.Println(err2)
				return
			}
			//refresh time.Now() as time could have passed
			attachment.Date = time.Now()
			err3 := Models.UpdateAttachmentDateByFileNameAndUploadTypeOrCreate(&attachment)
			if err3 != nil {
				fmt.Println(err3.Error())

			}
		}()

		c.JSON(http.StatusOK, gin.H{"Status": "upload request sent", "attachment_name": filename})
	}

	//filepath := "http://localhost:8080/file/" + filename
	//c.JSON(http.StatusOK, gin.H{"filepath": filepath})

}

//GetAttachmentByID ... Get the attachment by id
func GetAttachmentByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var attachment Models.Attachment
	err := Models.GetAttachmentByID(&attachment, id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}

	} else {
		c.JSON(http.StatusOK, attachment)
	}
}

//DwAttachmentByID ... Get the attachment by id
func DwAttachmentByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var attachment Models.Attachment
	err := Models.GetAttachmentByID(&attachment, id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	var user Models.User
	err = Models.GetUserByID(&user, fmt.Sprint(attachment.UserId))
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	key := user.Username + "/" + attachment.FileName
	diskpath := os.Getenv("USER_FILE_PATH") + "/" + key
	if attachment.UploadType == "s3" {
		//generate presigned url and redirect user
		s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		urlStr, err := PresignS3(s, user.Username+"/"+attachment.FileName, 300)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		log.Println("Generated presigned url for ", key)
		//c.JSON(http.StatusOK, gin.H{"presigned_url": urlStr})
		//redirect rather than returning json with url
		http.Redirect(c.Writer, c.Request, urlStr, http.StatusTemporaryRedirect)

	} else {
		// c.Header("Content-Description", "File Transfer")
		// c.Header("Content-Transfer-Encoding", "binary")
		// c.Header("Content-Disposition", "attachment; filename="+attachment.FileName)
		// c.Header("Content-Type", "application/octet-stream")
		c.File(diskpath)
	}

}

//DeleteAttachment ... Delete the attachment
func DeleteAttachment(c *gin.Context) {
	id := c.Params.ByName("id")
	//get attachment first
	var att Models.Attachment
	err := Models.GetAttachmentByID(&att, id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment not found!"})
		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	//get user
	var user Models.User
	err = Models.GetUserByID(&user, fmt.Sprint(att.UserId))
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	filepath := os.Getenv("USER_FILE_PATH") + "/" + user.Username + "/" + att.FileName
	if att.UploadType == "s3" {
		//delete from s3... call sdk api
		s, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = DeleteFileFromS3(s, user.Username+"/"+att.FileName)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		fmt.Println("S3 - Deleted file" + user.Username + "/" + att.FileName)
	} else {
		//check if file exists
		if _, err := os.Stat(filepath); err == nil {
			// remove it
			var err = os.Remove(filepath)
			if err != nil {
				fmt.Println(err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else if os.IsNotExist(err) {
			//nevermind if the file isn't there... just log
			fmt.Println(err.Error())

		} else {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		fmt.Println("FS - Deleted file" + user.Username + "/" + att.FileName)
	}
	var attachment Models.Attachment

	err = Models.DeleteAttachment(&attachment, id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, gin.H{"id": id, "deleted": true})
	}
}

func createFile(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, fileDir string, key string) error {

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_BUCKET")),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

//DeleteFileFromS3 deletes file from s3
func DeleteFileFromS3(s *session.Session, key string) error {

	_, err := s3.New(s).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(key),
	})
	return err
}

//PresignS3 gets presigned url
func PresignS3(s *session.Session, key string, durationSec uint) (string, error) {
	// Create S3 service client
	svc := s3.New(s)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(time.Duration(durationSec) * time.Second)

	if err != nil {
		log.Println("Failed to presign request s3", err)
		return "", err
	}

	return urlStr, nil

}
