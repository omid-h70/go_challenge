package mail

import (
	"github.com/stretchr/testify/require"
	"go_challenge/util"
)

func TestSendEmailWithGmail(t *testing.T) {

	if testing.Short() {
		//if short flag present skip it
		t.Skip()
	}

	cfg, err := util.LoadConfig("..")
	require.NoError(err)

	sender := NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddr, cfg.EmailSenderPass)
	subject := "A Test Mail"
	content := `
	<h1>Hello World </h1>
	<p> This is Unit test from go app, never give up dude !</p>
	`
	to := []string{cfg.EmailSenderName}
	attachFiles := []string{"../README.md"}
	err = sender.SendEmail(subject, content, nil, nil, attachFiles)
	require.NoError(err)
}
