# Self-Destruct Share

A modern, secure file sharing application with self-destructing files. Built with Go backend and a stunning HTML/CSS/JavaScript frontend.

## Features

- **Self-Destructing Files**: Files automatically delete after the specified number of downloads
- **Password Protection**: Optional password protection for sensitive files
- **Time-Limited**: Files expire after 5 minutes
- **Modern UI**: Beautiful, responsive design with smooth animations
- **Drag & Drop**: Easy file upload with drag and drop support
- **Download Limits**: Control how many times a file can be downloaded (1-10)
- **Real-time Feedback**: Toast notifications and loading states

## Architecture

### Backend (Go)
- **Framework**: Chi router for HTTP routing
- **Storage**: Redis for temporary file storage
- **Endpoints**:
  - `POST /upload` - Upload a file with optional password and download limit
  - `GET /file/{id}` - Download a file by ID

### Frontend (HTML/CSS/JavaScript)
- **Design**: Modern glassmorphism design with gradient backgrounds
- **Responsive**: Mobile-first responsive design
- **Features**: Drag & drop, file validation, copy to clipboard
- **Pages**: Main upload page and dedicated download page

## Getting Started

### Prerequisites
- Go 1.23.2 or later
- Redis server running on localhost:6379

### Backend Setup
1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Start Redis server
4. Run the backend:
   ```bash
   go run main.go
   ```
   The API will be available at `http://localhost:8080`

### Frontend Setup
1. Open `index.html` in a web browser
2. The frontend will automatically connect to the backend API

## Usage

### Uploading Files
1. Open the application in your browser
2. Drag and drop a file or click to select one
3. Set the number of downloads allowed (1-10)
4. Optionally set a password
5. Click "Upload & Generate Link"
6. Copy the generated link to share

### Downloading Files
1. Use the generated link or visit the download page
2. Enter the file ID
3. Enter the password if required
4. Click "Download File"

## File Structure

```
├── main.go                 # Go backend entry point
├── handlers/               # HTTP handlers
│   ├── upload.go          # File upload handler
│   └── download.go        # File download handler
├── storage/               # Storage layer
│   ├── redis.go          # Redis operations
│   └── types.go          # Data structures
├── utils/                 # Utility functions
│   └── helper.go         # ID generation
├── index.html            # Main frontend page
├── download.html         # Download page
├── styles.css            # CSS styles
├── script.js             # JavaScript functionality
└── README.md             # This file
```

## API Endpoints

### POST /upload
Upload a file with optional parameters.

**Form Data:**
- `file` (required): The file to upload
- `downloads` (optional): Number of downloads allowed (default: 1)
- `password` (optional): Password protection

**Response:**
```
File uploaded--Download:/file/{id}
```

### GET /file/{id}
Download a file by ID.

**Query Parameters:**
- `password` (optional): Password if file is protected

**Response:**
- File download with appropriate headers
- 404 if file not found or expired
- 403 if wrong password
- 410 if no downloads left

## Configuration

### Redis Configuration
The Redis connection is configured in `storage/redis.go`:
```go
rdb = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    DB:   0,
})
```

### File TTL
Files expire after 5 minutes (configurable in `storage/redis.go`):
```go
return rdb.Set(ctx, key, u, time.Minute*5).Err()
```

## Security Features

- **Password Protection**: Optional password for sensitive files
- **Download Limits**: Prevent unlimited downloads
- **Time Expiration**: Files automatically expire after 5 minutes
- **Self-Destruction**: Files delete after final download
- **File Validation**: Client-side file size validation (10MB limit)

## Browser Support

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is open source and available under the MIT License.

## Troubleshooting

### Common Issues

1. **Redis Connection Error**: Ensure Redis is running on localhost:6379
2. **File Upload Fails**: Check file size (max 10MB) and network connection
3. **Download Fails**: Verify file ID and password (if required)
4. **CORS Issues**: The frontend and backend must be on the same origin or CORS must be configured

### Debug Mode
Enable debug logging by setting the log level in the Go application.

## Future Enhancements

- [ ] User authentication
- [ ] File preview
- [ ] Bulk upload
- [ ] Custom expiration times
- [ ] File sharing analytics
- [ ] Mobile app
- [ ] API rate limiting
- [ ] File encryption
