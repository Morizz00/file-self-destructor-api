---

📦 Self-Destructing File Share API

A minimal Go API for uploading files that delete themselves — after a single download or 5 minutes.
Perfect for sharing secrets, cursed contracts, or sensitive memes.


---

🚀 Features

🔥 Upload files via REST

💥 Auto-deletes after 1 download or 5 minutes, whichever comes first

⚡ Powered by Redis for in-memory speed and TTL-based cleanup

📬 Returns a one-time-use download link



---

🛠️ Tech Stack

Layer	Tech

Language	Go
Router	chi
Storage	Redis (with TTL)
HTTP Client	Postman / Curl



---

📂 Folder Structure

selfdestruct-share-api/
│
├── main.go                 # Router setup
├── handlers/
│   ├── upload.go           # POST /upload
│   └── download.go         # GET /file/:id
├── storage/
│   └── redis.go            # Redis logic
├── utils/
│   └── helper.go           # ID generator
├── go.mod
└── README.md               # this file


---

📬 API Endpoints

POST /upload

Upload a file using multipart/form-data

Field name: file


curl -F "file=@secret.txt" http://localhost:8080/upload

Response:

✅ File uploaded!
Download link: /file/4a2f9d (valid for 5 mins or 1 download)


---

GET /file/:id

Downloads the file

Deletes it immediately after download


curl http://localhost:8080/file/4a2f9d --output recovered.txt

If expired or already downloaded:

404 Not Found


---

🧪 Running the App

1. Start Redis

redis-server

Or with Docker:

docker run -p 6379:6379 redis


---

2. Run the Go app

go run main.go


---

🔐 Security Notes

This is designed for short-term, secure transfers

Once a file is downloaded, it's gone forever

Don’t use this for long-term storage or tax returns (unless you like chaos)



---

🧙 Future Features (In Progress)

[x] Metadata support (filename, MIME)

[ ] Password protection

[ ] Download count limits

[ ] Custom aliases

[ ] Email-delivery of download links



---

✨ Inspired by

Burn After Reading

Telegram's self-destructing messages




---

