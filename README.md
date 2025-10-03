# FileOrcha

A modern, secure file sharing application with self-destructing files and blockchain proof capabilities. Built with Go backend and a clean, responsive frontend.

## Features

### Core Features
- **Self-Destructing Files** - Files automatically delete after the specified number of downloads
- **Custom Expiry Times** - Set expiration from 1 minute to 7 days
- **Password Protection** - Optional password protection for sensitive files
- **Download Limits** - Control how many times a file can be downloaded (1-10)
- **Custom URLs** - Create memorable links with custom slugs (e.g., /file/my-document)
- **Large File Support** - Upload files up to 50MB

### Advanced Features
- **QR Code Generation** - Instant QR codes for mobile sharing
- **Live Statistics** - Real-time countdown timers and download tracking
- **File Preview** - Preview images, text files, and PDFs before upload
- **Smart File Type Detection** - Automatic icon assignment for 20+ file types
- **Expiry Presets** - Quick selection buttons (5 minutes, 1 hour, 1 day, 7 days)

### Security & Privacy
- **Time-Limited Links** - All files expire automatically
- **One-Time Downloads** - Option for single-use links
- **Password Encryption** - Secure password hashing
- **No Permanent Storage** - Files deleted from server after expiry

### User Experience
- **Modern Dark Theme** - Professional, eye-friendly interface
- **Light Mode Support** - Toggle between dark and light themes
- **Drag & Drop Upload** - Easy file selection
- **Copy to Clipboard** - One-click link copying
- **Mobile Responsive** - Works seamlessly on all devices
- **Real-Time Notifications** - Toast alerts for all actions
- **Keyboard Shortcuts** - Power user features

## Architecture

### Backend (Go)
- **Framework** - Chi router for HTTP routing with middleware support
- **Storage** - Redis for temporary file storage with automatic expiration
- **File Processing** - Efficient chunked upload handling
- **API Endpoints**
  - `POST /upload` - Upload file with optional password and download limit
  - `GET /file/{id}` - Download file by ID or custom slug

### Frontend (HTML/CSS/JavaScript)
- **Design** - Modern dark theme with gradient accents
- **Responsive** - Mobile-first responsive design
- **Features** - Drag & drop, file validation, QR codes, live timers
- **Pages** - Main upload page and dedicated download page

## Getting Started

### Prerequisites
- Go 1.23.2 or later
- Redis server running on localhost:6379

### Local Development

1. Clone the repository
```bash
git clone https://github.com/Morizz00/fileorcha.git
cd fileorcha
```

2. Install dependencies
```bash
go mod tidy
```

3. Start Redis server
```bash
redis-server
```

4. Run the application
```bash
go run main.go
```

The application will be available at `http://localhost:8080`

5. Open `index.html` in your browser or serve the frontend separately

## Usage

### Uploading Files
1. Open the application in your browser
2. Drag and drop a file or click to select one
3. Configure settings:
   - Set number of downloads allowed (1-10)
   - Choose expiry time (1 minute to 7 days)
   - Optionally set a password
   - Optionally create a custom URL slug
4. Click "Upload & Generate Link"
5. Copy the generated link or download the QR code

### Downloading Files
1. Use the generated link or visit the download page
2. Enter the file ID if required
3. Enter the password if the file is protected
4. Click "Download File"
5. The file will download and the download counter will decrease

## Configuration

### Environment Variables
- `PORT` - Server port (default: 8080)
- `REDIS_URL` - Redis connection string (default: localhost:6379)

### File Limits
- Maximum file size: 50MB
- Maximum downloads per file: 10
- Maximum expiry time: 7 days (10,080 minutes)

## API Documentation

### POST /upload
Upload a file with optional parameters.

**Form Data:**
- `file` (required) - The file to upload
- `downloads` (optional) - Number of downloads allowed (default: 1, max: 10)
- `expiry` (optional) - Expiry time in minutes (default: 5, max: 10080)
- `password` (optional) - Password protection
- `slug` (optional) - Custom URL slug (lowercase, numbers, hyphens only)

**Response:**
```
File uploaded--Download:/file/{id}
```

### GET /file/{id}
Download a file by ID or custom slug.

**Query Parameters:**
- `password` (optional) - Password if file is protected

**Response:**
- File download with appropriate headers
- 404 if file not found or expired
- 403 if wrong password
- 410 if no downloads remaining

## Deployment

### Render.com (Recommended)

1. Push your code to GitHub

2. Create a Redis instance on Render
   - Go to render.com
   - Create new Redis instance
   - Copy the Internal Redis URL

3. Create a Web Service
   - Connect your GitHub repository
   - Settings:
     - Environment: Go
     - Build Command: `go build -o main`
     - Start Command: `./main`
   - Add environment variable:
     - `REDIS_URL` = your Redis Internal URL

4. Deploy and access your live URL

### Railway.app

1. Push to GitHub
2. Connect repository to Railway
3. Add Redis service
4. Set environment variable: `REDIS_URL=redis://redis:6379`
5. Deploy automatically

## Project Structure

```
fileorcha/
├── main.go              # Go backend entry point
├── handlers/            # HTTP request handlers
│   ├── upload.go       # File upload handler
│   └── download.go     # File download handler
├── storage/            # Storage layer
│   ├── redis.go       # Redis operations
│   └── types.go       # Data structures
├── utils/              # Utility functions
│   └── helper.go      # ID generation
├── index.html         # Main upload page
├── download.html      # Download page
├── styles.css         # Application styles
├── script.js          # Frontend functionality
├── go.mod             # Go dependencies
└── README.md          # This file
```

## Technical Details

### File Storage
Files are stored in Redis with the following structure:
- Key: file ID or custom slug
- Value: JSON containing file data, metadata, and settings
- TTL: Automatic expiration based on user-specified time

### Security Features
- Password hashing for protected files
- Custom URL validation (alphanumeric and hyphens only)
- Automatic file deletion after expiry or download limit
- No permanent file storage
- CORS configuration for cross-origin requests

### Performance
- Efficient Redis operations with connection pooling
- Chunked file uploads for large files
- Client-side validation before upload
- Optimized asset delivery

## Browser Support

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+
- Mobile browsers (iOS Safari, Chrome Mobile)

## Roadmap

### Planned Features
- End-to-end encryption with zero-knowledge architecture
- Blockchain proof of transfer (Polygon network)
- Live file preview (PDFs, images, videos)
- Browser extension for quick sharing
- Mobile applications (iOS/Android)
- API webhooks for integrations
- Screenshot-to-link desktop application

See FUTURE_FEATURES.md for detailed roadmap.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the MIT License.

## Support

For issues, questions, or feature requests, please open an issue on GitHub.

## Acknowledgments

Built with modern web technologies and best practices for security and user experience.
