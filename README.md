Referral Service App
=========================
A microservice designed to manage referral program, program membership and member referrals.
Generate and validate referral codes to track referral progress.

Business Use case
=================
- Use the program apis to create referral program say for your new product
- Enroll your partners as members to referral program
  - this generates a unique referral code for program member
- Your partners can use this referral code add their referral in the system
- When the referral takes action like creating account or making purchase, the program reward rules evaluates the action and triggers the reward
- Once reward is triggered the system will issue a payout to member / referral based on their preference (ex. paypal transaction or tango card)

Features
========
V0
- Create and Manage referral programs
- Manage referral program members
  - Support to generate unique referral code for each program member
- Manage member referrals

Future releases
---------------
V1
- Program reward rule configuration
- Rewards API and worklfow support
- Reward Payout preference and transaction support

V2
- Member stats support
  - number of referral shares
  - pending referrals
  - approved referrals
  - rewards

Getting Started
===============

Prerequisites
- Go (v1.18 or later) — install from [official Go site]
  - Uber Fx lib - Dependency injection support

- gRPC and Protocol Buffers — for defining and generating service APIs

- PostgreSQL — for data storage

Installation

1. Clone the repository
```
git clone https://github.com/rajatrao/referral-service.git
cd referral-service
```
2. Install dependecies
```
go get go.uber.org/fx@latest
go get go.uber.org/zap@latest
go get go.uber.org/config@latest
go get github.com/grpc-ecosystem/grpc-gateway/v2/runtime
go get database/sql
go get github.com/jmoiron/sqlx
go get github.com/google/uuid
copy proto files for google apis
```
3. Configure environment
  - Update ./config/*.yaml to add relevant configurations
4. Proto file changed ?
```
protoc --proto_path=proto/ --go_out=proto/referral/ \
     --go_opt=module=referral-service/proto/referral \
     --go-grpc_out=proto/referral/ \
    --go-grpc_opt=module=referral-service/proto/referral \
    --grpc-gateway_out=proto/referral \
    --grpc-gateway_opt=module=referral-service/proto/referral \
    ./proto/referral/referral.proto
```
4. Run the service
```
go build referral-service
./referral-service
```


Usage

This service exposes gRPC and http methods for referral management. See API documentation for details on available routes and response formats
Http calls are forwarded to gRPC methods for execution.

> The project uses the Go Clean Architecture approach to ensure well-structured, maintainable, and testable code. 
>  This style organizes the codebase into distinct layers—such as domain, controllers, and repository —to enforce separation of concerns and adaptability
>   - Domain: Core business logic and entities.
>   - Controllers: Application-specific business rules.
>   - Repository: Data access and persistence interfaces.
>   - Handlers: Handling incoming requests and responses.

1. Referral Program management
   - Add referral program

     request:
     ```
     curl --location --request POST 'http://127.0.0.1:8090/api/v1/programs' \
      --header 'Content-Type: text/plain' \
      --data-raw '{
          "name": "New_User_Referral_Program",
          "title": "New user referral program",
          "active": true
      }'
     ```
     response:
     ```
      {
        "id": "4790309b-d19e-4c46-8677-237bbacc0adc"
      }
     ```
   - View referral programs

     paginated-request:
     ```
      curl --location --request GET 'http://127.0.0.1:8090/api/v1/programs?page=1&size=10'
     ```
     response:
     ```
         {
            "programs": [
                {
                    "id": "b5142d77-2c6b-4dcb-8e78-42db0658550c",
                    "name": "New_Product_Referral_Program",
                    "title": "New product referral program",
                    "active": true,
                    "createdat": "1757270722",
                    "updatedat": "1757270722"
                },
                {
                    "id": "4790309b-d19e-4c46-8677-237bbacc0adc",
                    "name": "New_User_Referral_Program",
                    "title": "New user referral program",
                    "active": true,
                    "createdat": "1757275917",
                    "updatedat": "1757275917"
                }
            ]
        }

     ```
    - Update referral program

       request:
       ```
        curl --location --request PUT 'http://127.0.0.1:8090/api/v1/programs' \
        --header 'Content-Type: text/plain' \
        --data-raw '{
            "active": false,
            "id": "b5142d77-2c6b-4dcb-8e78-42db0658550c"
        }'
       ```
     
       response:
       ```
         {
            "program": {
                "id": "b5142d77-2c6b-4dcb-8e78-42db0658550c",
                "name": "New_Product_Referral_Program",
                "title": "New product referral program",
                "active": false,
                "createdat": "1757270722",
                "updatedat": "1757271003"
            }
        }
       ```

2. Referral program membership management

    - Add program member

        request:
        ```  
          curl --location --request POST 'http://127.0.0.1:8090/api/v1/members' \
          --header 'Content-Type: text/plain' \
          --data-raw '{
              "first_name": "John",
              "last_name":"smith",
              "email": "john@gmail.com",
              "program_id": "b5142d77-2c6b-4dcb-8e78-42db0658550c",
          }'
        ```
    
        response:
        ```
          {
            "id": "fc21290d-4587-423c-83f6-aa2e61089303"
          }
        ```

    - View members

        request:
        ```
          curl --location --request GET 'http://127.0.0.1:8090/api/v1/members'
        ```
      
        response:
         ```
                {
                  "members": [
                      {
                          "id": "fc21290d-4587-423c-83f6-aa2e61089303",
                          "firstName": "john",
                          "lastName": "smith",
                          "email": "john@gmail.com",
                          "programId": "b5142d77-2c6b-4dcb-8e78-42db0658550c",
                          "referralCode": "tfazu",
                          "isActive": true,
                          "createdAt": "1757276260",
                          "updatedAt": "1757276260"
                      },
                      {
                          "id": "2cfa5cc7-e101-4c72-bfb9-81bbc14e4837",
                          "firstName": "alex",
                          "lastName": "manning",
                          "email": "alex@gmail.com",
                          "programId": "b5142d77-2c6b-4dcb-8e78-42db0658550c",
                          "referralCode": "yolqq",
                          "isActive": false,
                          "createdAt": "1757277161",
                          "updatedAt": "1757277161"
                      }
                  ]
              }
         ```

3. Member referrals management

     - Add referral
     
        request:
        ```
          curl --location --request POST 'http://127.0.0.1:8090/api/v1/referrals' \
          --header 'Content-Type: text/plain' \
          --data-raw '{
              "first_name": "carry",
              "last_name":"joshi",
              "email": "carry@gmail.com",
              "phone": "111-111-1111",
              "referral_code":"tfazu"
          }'
        ```
  
        response:
        ```
           {
             "id": "5ee48eeb-7cd0-41f8-83cf-b821d7fadc3d"
           }
        ```

     - View referrals

         request:
    
         ```
          curl --location --request GET 'http://127.0.0.1:8090/api/v1/referrals'
         ```
    
         response:
         ```
            {
                "referrals": [
                    {
                        "id": "5ee48eeb-7cd0-41f8-83cf-b821d7fadc3d",
                        "firstName": "carry",
                        "lastName": "joshi",
                        "email": "carry@gmail.com",
                        "phone": "111-111-1111",
                        "referralCode": "tfazu",
                        "programId": "b5142d77-2c6b-4dcb-8e78-42db0658550c",
                        "referringMemberId": "fc21290d-4587-423c-83f6-aa2e61089303",
                        "status": "pending",
                        "createdAt": "1757288651",
                        "updatedAt": "1757288651"
                    }
                ]
            }
         ```

Data model
```
 - repository/schema.go
```

