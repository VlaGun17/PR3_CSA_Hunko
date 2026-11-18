openssl req -x509 -newkey rsa:2048 -nodes -keyout ca.key -out ca.crt -days 365 -subj "/CN=MyRootCA"

openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr -subj "/CN=localhost"
echo subjectAltName=DNS:localhost,IP:127.0.0.1>server-ext.cnf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile server-ext.cnf
del server-ext.cnf

openssl req -newkey rsa:2048 -nodes -keyout client.key -out client.csr -subj "/CN=dispatcher-agent-01"
echo subjectAltName=DNS:dispatcher-agent-01>client-ext.cnf
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -extfile client-ext.cnf
del client-ext.cnf