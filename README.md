##Idea
Currently I have subscribed to TIKR which costs me $20/month to get fundamental data like earnings, cashflow etc. This project is to bring that cost down by having a self-implemented solution using AlphaVantage API and my impeccable coding skills in golang with the help of github copilot.

- Write a application which takes a post request with a json payload like
`
{"symbol": "META", "function": "cash_flow", "api_key":"xxxxx"}
`
- Store the data in mongo-db which will be deployed in the same cluster
- Use Grafana to visualize the data stored in mongo-db
myapp/
├── main.go
├── internal/
│   ├── api/
│   │   └── handler.go
│   ├── db/
│   │   └── mongo.go
│   ├── models/
│   │   └── data.go
│   ├── services/
│   │   └── alpha_vantage.go
├── pkg/
│   └── utils/
│       └── utils.go
├── go.mod
├── go.sum
└── README.md
