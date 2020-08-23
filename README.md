# userBalanceService
test application for internship

# to build application:
`sudo docker build -t alex https://github.com/saskamegaprogrammist/userBalanceService.git`

# to run application:
`sudo docker run -p 5000:5000 --name alex -t alex`

# API

# Add funds
"/funds/add" **POST**

### Answers

- 200 - OK
- 400 - Bad Request
- 500 - Internal error

### JSON example

{"user_id": 1, "sum": 114.3}

## CURL request example

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"user_id": 1, "sum": 114.3}' \
   http://localhost:5000/funds/add
   
# Withdraw funds
"/funds/withdraw" **POST**

### Answers

- 200 - OK
- 400 - Bad Request
- 402 - Not enough funds
- 500 - Internal error

### JSON example

{"user_id": 4, "sum": 11}

## CURL request example

curl --header "Content-Type: application/json" \ 
 --request POST \
  --data '{"user_id": 4, "sum": 11}' \
  http://localhost:5000/funds/withdraw

# Get balance
"/funds/get" **POST**

### Answers

- 200 - OK
- 400 - Bad Request
- 500 - Internal error

### Query params

- currency 

    *"USD"* - currency type

### JSON example

{"user": 2}
    
## CURL request example

 curl --header "Content-Type: application/json"  \
    --request POST \
    --data '{"user": 2}' \
     http://localhost:5000/funds/get?currency=USD
     
### JSON answer example

{"user_id":4,"balance":4.02427764,"currency":"USD"}         
     

# Transfer funds
"/funds/transfer" **POST**

### Answers

- 200 - OK
- 400 - Bad Request
- 402 - Not enough funds
- 500 - Internal error

### JSON example

{"user_id": 1, "sum": 100, "user_from_id": 2}

## CURL request example

curl --header "Content-Type: application/json"  \
 --request POST \
 --data '{"user_id": 1, "sum": 11423.32, "user_from_id": 2}' \
 http://localhost:5000/funds/transfer

# Get transaction list
"/funds/details" **POST**

### Answers

- 200 - OK
- 400 - Bad Request
- 500 - Internal error

### Query params

- limit 

    *"10"* - number
- since

    *"2020-01-02T15:04:05.999999-07:00"* - timestamp
- sort

    *"sum"* - sorting by transaction sum 
    
    *"date"* - sorting by time of transaction creation
- desc  

    *"true"* - ordering by most resent / with biggest sum
    
    *"false"* - ordering by oldest / with lowest sum

### JSON example

{"user": 2}
    
## CURL request example

curl --header "Content-Type: application/json" \
  --request POST\
   --data '{"user": 4}'  \
   http://localhost:5000/funds/details?sort=sum
   
### JSON answer example

[{"user_id":1,"user_from_id":0,"operation_type":1,"sum":100,"balance":100,"balance_from":0,"created":"2020-08-02T00:10:09.887457+03:00"},
{"user_id":1,"user_from_id":0,"operation_type":1,"sum":200,"balance":300,"balance_from":0,"created":"2020-08-03T00:10:09.887457+03:00"}]