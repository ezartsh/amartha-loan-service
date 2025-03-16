# Loan Service

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

This repository contains a REST API to handle loan system.

## Prerequisite
- `Go` Make sure Go binary installed on local system. Visit this link [Go Installation](https://go.dev/doc/install) and follow the instruction to install Go.
- `Mailtrap` Create [Mailtrap](https://mailtrap.io) account to enable sending mail feature.

## Installation
### Clone the project

```
$ git clone https://github.com/ezartsh/amartha-loan-service.git
$ cd amartha-loan-service
```

### Structure
```
├── config
│   ├── env.go           // Configuration
│── handler              // API core handlers
│   ├── loan             // API handler for loan feature
│   │   ├── approve.go   // APIs for handling approve loan
│   │   ├── controller.go     // Controller initiation
│   │   ├── create.go    // APIs for handling create new loan
│   │   ├── disburse.go    // APIs for handling loan disbursement
│   │   ├── dto.go       // DTO/Entities model for store input and output request
│   │   ├── index.go     // Api handler for showing the list of loans and detail specific loan
│   │   ├── invest.go    // APIs for handling invest loan
│   │   └── router.go    // Register APIs endpoints
│── logs // Folder contains all logs that registered when the application started.
│── model                // APIs core model 
│   ├── investor.go      // Core model Investor
│   ├── loan_investor.go      // Core model relation from Loan to Investor
│   ├── loan.go      // Core model Loan
│   ...
└── main.go // Entry point application
```

### Configuration
Copy .env from .env template
```
$ cp .env.example .env
```
* `HTTP_PORT` : Http port when running the application (e.g `3000`)
* `COMPANY_MARGIN_IN_PERCENT` : Percentage margin profit that company gain for each loan transaction. (e.g `2`) means that company margin profit are 2%
* `MINIMUM_LOAN_AMOUNT` : Minimum amount of money (IDR) that borrower can request.
* `MAXIMUM_LOAN_AMOUNT` : Maximum amount of money (IDR) that borrower can request.
* `MAIL_HOST` : Mail host.
* `MAIL_PORT` : Mail port.
* `MAIL_USERNAME` : Mail username.
* `MAIL_PASSWORD` : Mail password.
* `MAIL_EMAIL_SENDER` : Email sender.

For this example project, there is some rules and validation when setting the configuration :
- `COMPANY_MARGIN_IN_PERCENT` value can not be lower than 1%. Assuming that 1% is the minimum margin profit for the company.
- `MINIMUM_LOAN_AMOUNT` value can not be lower than Rp. 10.000 
- `MAXIMUM_LOAN_AMOUNT` value must be greater than `MINIMUM_LOAN_AMOUNT`
---
> For this example. Interest Rate are declared for every new loan created. ROI value will be generated based on the remaining interest deductions and profit margin. (e.g ROI = (Interest Rate - Company Margin Profit))

## Run Service
After setup the configuration, service can be run by executing this command
```
$ go run main.go
```

## Authentication
None

## Database / Entities

### **`Note` : Every data stored in memory.**
It will be erased once the application stoped.

Loan Status are enums in sequence number, which every action can only be trigger if the intention matched with the previous status. e.g. api `invest` can be trigger if loan status is `1` (Approved) :
- `0`: Proposed
- `1`: Approved
- `2`: Invested
- `3`: Disbursed

`EmployeeId`, `BorrowerId` and `InvestorId` are static. The api expect it to be provided on the body request. No user validation access for each api.
The only validation is 1 investor can only invest 1 for each loan.

## REST API
The REST API is described below.
Postman Collection JSON also provided on the root folder of this project `loan-service/Postman.json`

### Get list of Loans

#### Request

`GET /loans`

    curl --location 'http://localhost:3000/loans' --header 'Accept: application/json'

#### Response

    HTTP/1.1 200 OK
    Date: Sun, 16 Mar 2025 07:10:15 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 2

    []

### Create a new Loan (Propose Loan)

#### Request

`POST /loans`

    curl --location 'http://localhost:3000/loans/create' \
    --form 'borrower_id="1"' \
    --form 'principal_amount="10000"' \
    --form 'interest_rate="10"'

#### Response

    HTTP/1.1 201 Created
    Date: Sun, 16 Mar 2025 07:16:31 GMT
    Status: 201 Created
    Connection: close
    Content-Type: application/json
    Content-Length: 249

    {"id":1,"borrower_id":1,"principal_amount":10000,"interest_rate":10,"roi":8.3,"approval_evidence":"","approval_employee_id":0,"approval_date":null,"disburse_evidence":"","disburse_employee_id":0,"disburse_date":null,"loan_investors":[],"status":0,"created_at":"2025-03-16T14:56:41.767056113+07:00","updated_at":"2025-03-16T14:56:41.767056113+07:00"}

### Get a specific Loan

#### Request

`GET /loans/:id`

    curl --location 'http://localhost:3000/loans/1' --header 'Accept: application/json'

#### Response

    HTTP/1.1 200 OK
    Date: Sun, 16 Mar 2025 07:22:26 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 249

    {"id":1,"borrower_id":1,"principal_amount":10000,"interest_rate":10,"roi":8.3,"approval_evidence":"","approval_employee_id":0,"approval_date":null,"disburse_evidence":"","disburse_employee_id":0,"disburse_date":null,"loan_investors":null,"status":0}
    
## Approve Loan

### Request

`PUT /loans/:id/approve`

    curl --location --request PUT 'http://localhost:3000/loans/1/approve' \
    --header 'Accept: application/json' \
    --form 'employee_id="1"' \
    --form 'evidence_picture=@"/home/username/document/approval_evidence.pdf"'

### Response

    HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:31 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 282

    {"id":1,"borrower_id":1,"principal_amount":10000,"interest_rate":10,"roi":8.3,"approval_evidence":"/home/user/loan/upload/approval_evidence/8hdp83912.pdf","approval_employee_id":1,"approval_date":"2025-03-16T14:25:47.379445425+07:00","disburse_evidence":"","disburse_employee_id":0,"disburse_date":null,"loan_investors":[],"status":1,"created_at":"2025-03-16T14:56:41.767056113+07:00","updated_at":"2025-03-16T14:57:39.288747261+07:00"}

## Invest Loan

### Request

`PUT /loans/:id/invest`

    curl --location --request PUT 'http://localhost:3000/loans/1/invest' \
    --header 'Accept: application/json' \
    --form 'investor_id="1"' \
    --form 'amount="2000"'

### Response

    HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:31 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 282

    {"id":1,"borrower_id":1,"principal_amount":10000,"interest_rate":10,"roi":8.3,"approval_evidence":"","approval_employee_id":1,"approval_date":"2025-03-16T14:52:14.918426912+07:00","disburse_evidence":"","disburse_employee_id":0,"disburse_date":null,"loan_investors":[{"id":1,"investor_id":5,"amount":2000,"created_at":"2025-03-16T14:56:41.767056113+07:00","updated_at":"2025-03-16T14:57:39.288747261+07:00"}],"status":1}

## Disburse Loan

### Request

`PUT /loans/:id/disburse`

    curl --location --request PUT 'http://localhost:3000/loans/1/invest' \
    --header 'Accept: application/json' \
    --form 'investor_id="1"' \
    --form 'amount="2000"'

### Response

    HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:31 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 282

    {"id":1,"borrower_id":1,"principal_amount":10000,"interest_rate":10,"roi":8.3,"approval_evidence":"","approval_employee_id":1,"approval_date":"2025-03-16T14:52:14.918426912+07:00","disburse_evidence":"","disburse_employee_id":0,"disburse_date":null,"loan_investors":[{"id":1,"investor_id":5,"amount":2000,"created_at":"2025-03-16T14:56:41.767056113+07:00","updated_at":"2025-03-16T14:57:39.288747261+07:00"}],"status":1}