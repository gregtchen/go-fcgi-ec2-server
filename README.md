#go-fcgi-ec2-server
Create EC2 Instance
Install apache
Install Go
SFTP into instance and place fcgiserver.go & files in var/www/html
SSH into instance
go run fcgiserver.go -local=":8080"
Open port 8080 on EC2 Instance
visit http://ec2ip:8080/script.html
Any file visited through that address and socket will be through the FCGI

Note: fcgi.serve() can only run on linux environments