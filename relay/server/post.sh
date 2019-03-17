curl -v -H Content-Type: application/json -X POST -d '{"user":"bob"}' http://localhost:8080/relay/createuser
#curl -v -H "Content-Type: application/json" -X POST http://localhost:8080/relay/createuser?user=bob
curl -v -H Content-Type: application/json -X POST -d '{"type":"tactic","value":"433"}' http://localhost:8080/relay/v1/bob/action
