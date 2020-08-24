//Models/AttachmentModel.go
package Models

import "time"

type Attachment struct {
	Id         uint      `json:"id"`
	FileName   string    `json:"file_name" gorm:"not null"`
	Date       time.Time `json:"date" gorm:"not null"`
	UploadType string    `json:"upload_type" gorm:"not null"`
	UserId     uint      `json:"user_id" gorm:"not null"`
}

func (b *Attachment) TableName() string {
	return "attachment"
}
