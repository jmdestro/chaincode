# Chaincode

Métodos disponíveis:
- read_everything
- readTicket
- saveTicket
- saveRCMSTicket
- readIBMAsset
- saveIBMAsset
- getTicketsByRange
- getTicketHistory
- getTicketTimestamp
- getTickets

# read_everything
Ler todos os assets existentes (Employee, Ticket e IBMAsset).

URL:
/api/query

Payload:
```
{
"fcn": "read_everything",
"args": [""]
}
```

Resposta:
- Em caso de sucesso: um JSON com o objeto desejado
```
status 200
```
```
{"tickets":[{"docType":"ticket","id":"ticketId","openDate":"OpenDate","statusTicket":"StatusTicket","queue":"Queue","opsAssignTicket":"OpsAssignTicket","descriptionRequest":"DescriptionRequest","userRequest":"UserRequest","serialNumber":"00000","email":"Email","telephone":"Telephone","cep":"Cep","address":"Address","city":"City","uf":"Uf","opsUpdate":"OpsUpdate","logbook":"Logbook","tokenIn":"TokenIn","tokenOut":"TokenOut"},{"docType":"ticket","id":"ticketId1","openDate":"OpenDate","statusTicket":"StatusTicket","queue":"Queue","opsAssignTicket":"OpsAssignTicket","descriptionRequest":"DescriptionRequest","userRequest":"UserRequest","serialNumber":"00000","email":"Email","telephone":"Telephone","cep":"Cep","address":"Address","city":"City","uf":"Uf","opsUpdate":"OpsUpdate","logbook":"Logbook","tokenIn":"TokenIn","tokenOut":"TokenOut"}],"employee":null,"ibmasset":[{"docType":"ibm_asset","serialNumber":"00000","prod":"Prod","descriptionProduct":"DescriptionProduct","pwHardware":"PwHardware","pwHD":"PwHD","pwOS":"PwOS","opsUpdate":"OpsUpdate"}]}
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# readTicket
Ler os dados de um Ticket específico.

URL:
/api/query

Payload:
```
{
"fcn": "readTicket",
"args": ["ticketId"]
}
```

Resposta:
- Em caso de sucesso: um JSON com o objeto desejado
```
status 200
```
```
{"address":"Address","cep":"Cep","city":"City","descriptionProduct":"DescriptionProduct","descriptionRequest":"DescriptionRequest","docType":"ibm_asset","email":"Email","id":"ticketId","logbook":"Logbook","openDate":"OpenDate","opsAssignTicket":"OpsAssignTicket","opsUpdate":"OpsUpdate","prod":"Prod","pwHD":"PwHD","pwHardware":"PwHardware","pwOS":"PwOS","queue":"Queue","serialNumber":"00000","statusTicket":"StatusTicket","telephone":"Telephone","tokenIn":"TokenIn","tokenOut":"TokenOut","uf":"Uf","userRequest":"UserRequest"}
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# saveTicket
Salvar um ticket. Se existe, os dados são atualizados senão o ticket é criado.

URL:
/api/invoke

Payload:
```
{
"fcn": "saveTicket",
"args": ["ticketId1", "OpenDate", "StatusTicket", "Queue", "OpsAssignTicket", "DescriptionRequest", "UserRequest", "Email", "Telephone", "Cep", "Address", "City", "Uf", "00000", "Prod", "DescriptionProduct", "PwHardware", "PwHD", "PwOS", "OpsUpdate", "Logbook", "TokenIn", "TokenOut"]
}
```

Resposta:
- Em caso de sucesso: uma string com o id do objeto criado
```
status 200
```
```
072985fe4480d7b46a8da0b14517604742d869112ee992b5a59df7d149178712
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# saveRCMSTicket
Salvar um ticket usando o RCMS como base. Se existe, os dados são atualizados conforme regra de status senão o ticket é criado.

URL:
/api/invoke

Payload:
```
{
"fcn": "saveTicket",
"args": ["ticketId1", "OpenDate", "StatusTicket", "Queue", "DescriptionRequest", "UserRequest", "Telephone", "Cep", "Address", "City", "SerialNumber", "Prod", "DescriptionProduct"]
}
```

Resposta:
- Em caso de sucesso: uma string com o id do objeto criado
```
status 200
```
```
072985fe4480d7b46a8da0b14517604742d869112ee992b5a59df7d149178712
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# readIBMAsset
Ler um IBMAsset específico.

URL:
/api/query

Payload:
```
{
"fcn": "readIBMAsset",
"args": ["00000"]
}
```

Resposta:
- Em caso de sucesso: um JSON com o objeto desejado
```
status 200
```
```
{"ibmAsset":{"docType":"ibm_asset","serialNumber":"00000","prod":"Prod","descriptionProduct":"DescriptionProduct","pwHardware":"PwHardware","pwHD":"PwHD","pwOS":"PwOS","tickets":null,"opsUpdate":"OpsUpdate"},"ticketsDetails":null}
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# saveIBMAsset
Salvar um IBMAsset. Se já existe, atualiza os dados caso contrário cria um novo.

URL:
/api/invoke

Payload:
```
{
"fcn": "saveIBMAsset",
"args": ["00000", "Prod", "DescriptionProduct", "PwHardware", "PwHD", "PwOS", "ticketId1", "OpsUpdate"]
}
```

Resposta:
- Em caso de sucesso: uma string com o id do objeto criado
```
status 200
```
```
3ec00f5c44835b7d21189db57f08a0060f70f31c0111d7c69aa2b474629197ff
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# getTicketsByRange
Obter todos os tickets dentro de uma faixa de números.

URL:
/api/query

Payload:
```
{
	"fcn": "getTicketsByRange",
	"args": ["ticket", ""]
}
```

Resposta:
- Em caso de sucesso: um JSON com o objeto desejado
```
status 200
```
```
[{"Key":"ticketId", "Record":{"address":"Address","cep":"Cep","city":"City","descriptionRequest":"DescriptionRequest","docType":"ticket","email":"Email","id":"ticketId","logbook":"Logbook","openDate":"OpenDate","opsAssignTicket":"OpsAssignTicket","opsUpdate":"OpsUpdate","queue":"Queue","serialNumber":"00000","statusTicket":"StatusTicket","telephone":"Telephone","tokenIn":"TokenIn","tokenOut":"TokenOut","uf":"Uf","userRequest":"UserRequest"}},{"Key":"ticketId1", "Record":{"address":"Address","cep":"Cep","city":"City","descriptionRequest":"DescriptionRequest","docType":"ticket","email":"Email","id":"ticketId1","logbook":"Logbook","openDate":"OpenDate","opsAssignTicket":"OpsAssignTicket","opsUpdate":"OpsUpdate","queue":"Queue","serialNumber":"00000","statusTicket":"StatusTicket","telephone":"Telephone","tokenIn":"TokenIn","tokenOut":"TokenOut","uf":"Uf","userRequest":"UserRequest"}}]
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# getTicketHistory
Mostrar o histório (blocochain) de um ticket.

URL:
/api/query

Payload:
```
{
"fcn": "getTicketHistory",
"args": ["ticketId"]
}
```

Resposta:
- Em caso de sucesso: um JSON com o objeto desejado
```
status 200
```
```
[{"txId":"72955f517222720c096e7d2ee50a592e3d7295acd1237474cd34df19a8d66092","value":{"docType":"ticket","id":"ticketId","openDate":"OpenDate","statusTicket":"StatusTicket","queue":"Queue","opsAssignTicket":"OpsAssignTicket","descriptionRequest":"DescriptionRequest","userRequest":"UserRequest","serialNumber":"ticketId","email":"Email","telephone":"Telephone","cep":"Cep","address":"Address","city":"City","uf":"Uf","opsUpdate":"OpsUpdate","logbook":"Logbook","tokenIn":"TokenIn","tokenOut":"TokenOut"},"timestamp":"2017-08-02 22:08:58.484 +0000 UTC"},{"txId":"7901b02bd7efd079aa1281ee036dedc6dd03fecb3d1eae3e9947d45ead02a627","value":{"docType":"ticket","id":"ticketId","openDate":"OpenDate","statusTicket":"StatusTicket","queue":"Queue","opsAssignTicket":"OpsAssignTicket","descriptionRequest":"DescriptionRequest","userRequest":"UserRequest","serialNumber":"00000","email":"Email","telephone":"Telephone","cep":"Cep","address":"Address","city":"City","uf":"Uf","opsUpdate":"OpsUpdate","logbook":"Logbook","tokenIn":"TokenIn","tokenOut":"TokenOut"},"timestamp":"2017-08-02 22:12:59.566 +0000 UTC"}]
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# getTicketTimestamp
Retorna o último timestamp de um ticket.

URL:
/api/query

Payload:
```
{
"fcn": "getTicketTimestamp",
"args": ["ticketId"]
}
```

Resposta:
- Em caso de sucesso: um JSON com o objeto desejado
```
status 200
```
```
"2017-08-02 22:08:58.484 +0000 UTC"
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```

# getTickets
Retorna todos os tickets existentes.

URL:
/api/query

Payload:
```
{
"fcn": "getTickets",
"args": ["status"]
}
```
```
args é opcional. Se passar um string vazia, a resposta será a lista de todos os tickets, se passar um status específico, retornará somente os tickets com esse status.
```

Resposta:
- Em caso de sucesso: um JSON com o objeto desejado
```
status 200
```
```
[{"address":"Address","cep":"00000-000","city":"SP","descriptionProduct":"DescriptionProduct2","descriptionRequest":"DescriptionRequest","docType":"ticket","email":"Email","hasWarranty":"","ibmLocation":"TU","id":"ticketId3","logbook":"","openDate":"2017/08/08","opsAssignTicket":"OpsAssignTicket","opsUpdate":"cassior@br.ibm.com","prod":"Prod2","pwHD":"","pwHardware":"","pwOS":"","queue":"Queue","repairType":"repair1, repair2","serialNumber":"00000","statusTicket":"ACK","telephone":"Telephone","timestamp":"2017-08-10 13:19:58.42 +0000 UTC","tokenIn":"tokenIn","tokenOut":"","uf":"SP","userRequest":"UserRequest"},{"address":"Address","cep":"00000-000","city":"SP","descriptionProduct":"DescriptionProduct4","descriptionRequest":"DescriptionRequest","docType":"ticket","email":"Email","hasWarranty":"","ibmLocation":"TU","id":"ticketId4","logbook":"","openDate":"2017/08/08","opsAssignTicket":"OpsAssignTicket","opsUpdate":"cassior@br.ibm.com","prod":"Prod4","pwHD":"PwdHD4","pwHardware":"PwdHardware4","pwOS":"PwdOS4","queue":"Queue","repairType":"repair type,repair1,repair2","serialNumber":"00001","statusTicket":"ACK","telephone":"Telephone","timestamp":"2017-08-10 16:31:50.887 +0000 UTC","tokenIn":"tokenIn","tokenOut":"","uf":"SP","userRequest":"UserRequest"},{"address":"Address","cep":"00000-000","city":"SP","descriptionProduct":"DescriptionProduct4","descriptionRequest":"DescriptionRequest","docType":"ticket","email":"Email","hasWarranty":"","ibmLocation":"TU","id":"ticketId5","logbook":"","openDate":"2017/08/08","opsAssignTicket":"OpsAssignTicket","opsUpdate":"cassior@br.ibm.com","prod":"Prod4","pwHD":"PwdHD5","pwHardware":"PwdHardware5","pwOS":"PwdOS5","queue":"Queue","repairType":"repair1, repair2","serialNumber":"00005","statusTicket":"ACK","telephone":"Telephone","timestamp":"2017-08-10 16:34:48.365 +0000 UTC","tokenIn":"tokenIn","tokenOut":"","uf":"SP","userRequest":"UserRequest"}]
```

- Em caso de erro: um texto com a mensagem de erro dentro de "message"
```
status 400
```
```
chaincode error (status: 500, message: blablabla)
```
