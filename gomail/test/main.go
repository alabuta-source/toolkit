package main

import "github.com/alabuta-source/toolkit/gomail"

func main() {
	testCOnfig := &gomail.EmailSenderConfig{
		Password:     "mgwthszohqrlcrrj",
		From:         "suporte@alabuta.com",
		ServerConfig: "smtp.gmail.com",
		Port:         587,
	}

	sender := gomail.NewEmailSender(testCOnfig)
	err := sender.SendEmail("joaquim.borges1993@gmail.com",
		gomail.WithSubject("test options"),
		gomail.WithBody("donttttttt"),
	)

	if err != nil {
		print(err)
	}
}
