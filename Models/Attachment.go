//Models/Attachment.go
package Models

import (
	"fileshare/Config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//
func GetAllAttachmentsFromUser(attachment *[]Attachment, user_id uint) (err error) {
	if err = Config.DB.Where("user_id = ?", user_id).Find(attachment).Error; err != nil {
		return err
	}
	return nil
}

//
func CreateAttachment(attachment *Attachment) (err error) {
	if err = Config.DB.Create(attachment).Error; err != nil {
		return err
	}
	return nil
}

//
func GetAttachmentByIDAndFileNameAndUploadType(attachment *Attachment, id uint, filename string, upload_type string) (err error) {
	if err = Config.DB.Where(&Attachment{Id: id, FileName: filename, UploadType: upload_type}).First(attachment).Error; err != nil {
		return err
	}
	return nil
}

//
func GetAttachmentByID(attachment *Attachment, id string) (err error) {
	if err = Config.DB.Where("id = ?", id).First(attachment).Error; err != nil {
		return err
	}
	return nil
}

// func UpdateUser(user *User, id string) (err error) {
// 	fmt.Println(user)
// 	Config.DB.Save(user)
// 	return nil
// }
//

func UpdateAttachmentDateByFileNameAndUploadTypeOrCreate(attachment *Attachment) (err error) {
	if err = Config.DB.Where(Attachment{FileName: attachment.FileName, UploadType: attachment.UploadType}).Assign(Attachment{Date: time.Now()}).FirstOrCreate(attachment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAttachment(attachment *Attachment, id string) (err error) {
	if err = Config.DB.Where("id = ?", id).Delete(attachment).Error; err != nil {
		return err
	}
	return nil
}
