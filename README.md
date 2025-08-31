# KOPICHAT
Kopichat es una API de chatbot, el objetivo principal de este chatbot es hablar contigo y debatir para intentar convencerte de tu punto de vista.
@startuml
' ===== Icons AWS desde el release v20.0 =====
!define AWSPuml https://raw.githubusercontent.com/awslabs/aws-icons-for-plantuml/v20.0/dist
!include AWSPuml/AWSCommon.puml
!include AWSPuml/NetworkingContentDelivery/APIGateway.puml
!include AWSPuml/Compute/Lambda.puml
!include AWSPuml/Database/RDS.puml
!include AWSPuml/NetworkingContentDelivery/VPCNATGateway.puml
!include AWSPuml/NetworkingContentDelivery/VPCInternetGateway.puml

skinparam shadowing false
skinparam dpi 180
left to right direction
hide stereotype

actor "Client" as user

APIGateway(api, "API Gateway", "HTTP API")
Lambda(fn, "Lambda (Container Image)", "return conversation")
RDS(db, "RDS MySQL", "save data")
VPCInternetGateway(igw, "Internet Gateway", "")
VPCNATGateway(nat, "NAT Gateway", "")
cloud openai as "OpenAI API\n(Internet)"

rectangle "VPC (10.0.0.0/16)" {
  frame "Public Subnet" {
    igw
    nat
    
    
  }
  frame "Private Subnet" {
  fn
    db
  }
}

' Flujo de entrada
user --> api : HTTPS
api --> fn  : Invoke

' Acceso a BD dentro de la VPC
fn  --> db  : TCP 3306

' Egreso a Internet (para OpenAI)
''' fn  --> nat --> igw --> openai : HTTPS 443'''
@enduml
