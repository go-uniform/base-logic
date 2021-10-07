ssh-keygen -t rsa -b 4096 -m PEM -f resources/jwt.key
openssl rsa -in resources/jwt.key -pubout -outform PEM -out resources/jwt.pub
rm resources/jwt.key.pub