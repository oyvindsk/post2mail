
# post2mail - Take HTTP posts and mail the info out

## Usage
 
    (look over the code, change the mail subject & http port perhaps?)
    go build post2mail.go 
    ./post2mail --help

## Using gmail as a smtp server
The default vales are set up for using gmail to send mail through. You can get a unique random password from google called "App passwords". These can be revoked without changing your main password. See under Account >> Security.

## Todo
 - Error-checking, don't just die all the time :)
 - Make the email subject, http port and url configuratable
 - Logging

## Useful Links
 - http://nathanleclaire.com/blog/2013/12/17/sending-email-from-gmail-using-golang/
 - http://www.goinggo.net/2013/06/send-email-in-go-with-smtpsendmail.html

