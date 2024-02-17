package utils

import (
	"bytes"
	"encoding/base64"
	"net/smtp"
	"os"
	"regexp"
	"time"
)

// IsNormalEmail 验证邮箱
func IsNormalEmail(email string) bool {
	// 最大长度不能超过320
	if len(email) > 320 {
		return false
	}
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

// SendCaptchaByEmail 发送验证码
func SendCaptchaByEmail(recipientEmail string, imageData *bytes.Buffer) error {
	mailOfQQ := NewMyMailOfQQ()
	err := mailOfQQ.SendImg(recipientEmail, "数字验证码(5分钟后过期)", imageData.Bytes())
	return err
}

// SendTextByEmail 发送文字
func SendTextByEmail(recipientEmail, subject, text string) error {
	return NewMyMailOfQQ().SendText(recipientEmail, text, subject)
}

type MyMail struct {
	username string
	password string
	host     string
	port     string
	auth     smtp.Auth
	buffer   *bytes.Buffer
	boundary string
}

func NewMyMailOfQQ() *MyMail {
	qqMail := new(MyMail)
	qqMail.username = os.Getenv("MailName")
	qqMail.password = os.Getenv("MailPassword")
	qqMail.host = "smtp.qq.com"
	qqMail.port = "25"
	qqMail.auth = smtp.PlainAuth("", qqMail.username, qqMail.password, qqMail.host)
	qqMail.buffer = bytes.NewBuffer(nil)
	qqMail.boundary = "GoBoundary"
	return qqMail
}

func (m *MyMail) writeHeader(header map[string]string) {
	for key, value := range header {
		m.buffer.WriteString(key + ":" + value + "\r\n")
	}
	m.buffer.WriteString("\r\n")
}

func (m *MyMail) writeFile(fileData []byte, filename, contentType string) {
	// 文件前置操作
	m.buffer.WriteString("\r\n--" + m.boundary + "\r\n")
	m.buffer.WriteString("Content-Transfer-Encoding:base64\r\n")
	m.buffer.WriteString("Content-Type:" + contentType + ";name=\"" + filename + "\"\r\n")
	m.buffer.WriteString("Content-ID: <" + filename + "> \r\n\r\n")

	// 文件内容
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(fileData)))
	base64.StdEncoding.Encode(payload, fileData)
	m.buffer.WriteString("\r\n")
	for index, n := 0, len(payload); index < n; index++ {
		m.buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			m.buffer.WriteString("\r\n")
		}
	}
}

func (m *MyMail) writeHtmlTemplate(html string) {
	m.buffer.WriteString("\r\n--" + m.boundary + "\r\n")
	m.buffer.WriteString("Content-Type: text/html; charset=UTF-8 \r\n")
	m.buffer.WriteString(html)
	m.buffer.WriteString("\r\n--" + m.boundary + "--")
}

func (m *MyMail) SendImg(recipientEmail, subject string, fileData []byte) error {
	// header内容
	header := map[string]string{}
	header["From"] = m.username
	header["To"] = recipientEmail
	header["Subject"] = subject
	header["Content-Type"] = "multipart/related;boundary=" + m.boundary
	header["Date"] = time.Now().String()

	// 把header写入buffer
	m.writeHeader(header)

	// 写入文件
	m.writeFile(fileData, "1.png", "image/png")
	htmlTemplate :=
		`
		<html>
		<body>
			<img src="cid:1.png"><br>
		</body>
		</html>
	`

	// 写入html的template
	m.writeHtmlTemplate(htmlTemplate)

	// 发送
	err := smtp.SendMail(m.host+":"+m.port, m.auth, m.username, []string{recipientEmail}, m.buffer.Bytes())
	return err
}

func (m *MyMail) SendText(recipientEmail, text, subject string) error {
	contentType := "Content-Type: text/plain" + "; charset=UTF-8"
	// From需要在To之前
	msg := []byte("From: " + m.username + "\r\nTo: " + recipientEmail + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + text)
	return smtp.SendMail(m.host+":"+m.port, m.auth, m.username, []string{recipientEmail}, msg)
}
