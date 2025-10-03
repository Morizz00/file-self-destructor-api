// Global variables
const API_BASE_URL = 'http://localhost:8080';
let currentFileId = null;
let currentPassword = '';

// DOM elements
const uploadSection = document.getElementById('uploadSection');
const successSection = document.getElementById('successSection');
const downloadSection = document.getElementById('downloadSection');
const uploadForm = document.getElementById('uploadForm');
const downloadForm = document.getElementById('downloadForm');
const fileInput = document.getElementById('fileInput');
const fileInputLabel = document.querySelector('.file-input-label');
const uploadBtn = document.getElementById('uploadBtn');
const downloadBtn = document.getElementById('downloadBtn');
const copyBtn = document.getElementById('copyBtn');
const uploadAnotherBtn = document.getElementById('uploadAnotherBtn');
const testDownloadBtn = document.getElementById('testDownloadBtn');
const uploadModeBtn = document.getElementById('uploadModeBtn');
const downloadModeBtn = document.getElementById('downloadModeBtn');
const toastContainer = document.getElementById('toastContainer');

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    initializeEventListeners();
    showUploadSection();
    initializeDarkMode();
});

// Event Listeners
function initializeEventListeners() {
    // Upload form
    uploadForm.addEventListener('submit', handleUpload);
    
    // Download form
    downloadForm.addEventListener('submit', handleDownload);
    
    // File input
    fileInput.addEventListener('change', handleFileSelect);
    
    // Drag and drop
    fileInputLabel.addEventListener('dragover', handleDragOver);
    fileInputLabel.addEventListener('dragleave', handleDragLeave);
    fileInputLabel.addEventListener('drop', handleDrop);
    
    // Copy button
    copyBtn.addEventListener('click', copyToClipboard);
    
    // Copy ID button
    document.getElementById('copyIdBtn').addEventListener('click', copyFileId);

    const slugInput = document.getElementById('customSlug');
    slugInput.addEventListener('input',function(){
        this.value = this.value.toLowerCase().replace(/[^a-z0-9-]/g, '');
    
    // Visual feedback
    if (this.value && !/^[a-z0-9-]+$/.test(this.value)) {
        this.style.borderColor = '#ff4444';
    } else {
        this.style.borderColor = '';
    }
});
    
    // Download QR button
    document.getElementById('downloadQRBtn').addEventListener('click', downloadQRCode);
    
    // Navigation buttons
    uploadAnotherBtn.addEventListener('click', showUploadSection);
    testDownloadBtn.addEventListener('click', testDownload);
    uploadModeBtn.addEventListener('click', showUploadSection);
    downloadModeBtn.addEventListener('click', showDownloadSection);
    
    // Close preview button
    document.getElementById('closePreview').addEventListener('click', hidePreview);
    
    // Dark mode toggle
    document.getElementById('darkModeToggle').addEventListener('change', toggleDarkMode);
    
    // Expiry preset buttons
    document.querySelectorAll('.preset-btn').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.preventDefault();
            const minutes = this.getAttribute('data-minutes');
            document.getElementById('expiry').value = minutes;
            
            // Visual feedback
            document.querySelectorAll('.preset-btn').forEach(b => b.classList.remove('active'));
            this.classList.add('active');
        });
    });
}

// File handling
function handleFileSelect(event) {
    const file = event.target.files[0];
    if (file) {
        updateFileInputLabel(file);
        previewFile(file);
    }
}

function handleDragOver(event) {
    event.preventDefault();
    fileInputLabel.classList.add('dragover');
}

function handleDragLeave(event) {
    event.preventDefault();
    fileInputLabel.classList.remove('dragover');
}

function handleDrop(event) {
    event.preventDefault();
    fileInputLabel.classList.remove('dragover');
    
    const files = event.dataTransfer.files;
    if (files.length > 0) {
        fileInput.files = files;
        updateFileInputLabel(files[0]);
        previewFile(files[0]);
    }
}

function updateFileInputLabel(file) {
    const fileName = file.name;
    const fileSize = formatFileSize(file.size);
    
    fileInputLabel.innerHTML = `
        <i class="fas fa-file"></i>
        <span>${fileName}</span>
        <small>${fileSize}</small>
    `;
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Upload functionality
async function handleUpload(event) {
    event.preventDefault();
    
    const file = fileInput.files[0];
    const downloads = document.getElementById('downloads').value;
    const expiry = document.getElementById('expiry').value;
    const password = document.getElementById('password').value;
    
    try {
        validateFile(file);
    } catch (error) {
        showToast(error.message, 'error');
        return;
    }
    
    if (downloads < 1 || downloads > 10) {
        showToast('Downloads must be between 1 and 10', 'error');
        return;
    }
    
    if (expiry < 1 || expiry > 10080) {
        showToast('Expiry time must be between 1 minute and 7 days', 'error');
        return;
    }
    
    uploadBtn.disabled = true;
    uploadBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Uploading...';
    
    try {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('downloads', downloads);
        formData.append('expiry', expiry);
        if (password) {
            formData.append('password', password);
        }

        const customSlug = document.getElementById('customSlug').value.trim();
        if (customSlug) {
            if (!/^[a-z0-9-]+$/.test(customSlug)) {
                showToast('Custom link must contain only lowercase letters, numbers, and hyphens', 'error');
                uploadBtn.disabled = false;
                uploadBtn.innerHTML = '<i class="fas fa-rocket"></i> Upload & Generate Link';
                return;
            }
            formData.append('slug', customSlug);
        }

        
        
        const responseText = await uploadWithProgress(formData);
        console.log('Upload response:', responseText); // Debug log
        const fileId = extractFileId(responseText);
        console.log('Extracted file ID:', fileId); // Debug log
        
        if (fileId) {
            currentFileId = fileId;
            currentPassword = password;
            showSuccessSection(file, downloads, expiry, password);
            showToast('File uploaded successfully!', 'success');
        } else {
            console.error('Failed to extract file ID. Response was:', responseText);
            throw new Error('Failed to extract file ID from response');
        }
        
    } catch (error) {
        console.error('Upload error:', error);
        showToast(`Upload failed: ${error.message}`, 'error');
    } finally {
        uploadBtn.disabled = false;
        uploadBtn.innerHTML = '<i class="fas fa-rocket"></i> Upload & Generate Link';
    }
}
function startCountdown(mins, eId) {
    const el = document.getElementById(eId);
    let tot = mins * 60;

    const timer = setInterval(() => {
        const hours = Math.floor(tot / 3600);
        const minutes = Math.floor((tot % 3600) / 60);
        const secs = Math.floor(tot % 60);

        const timeString = hours > 0
            ? `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
            : `${minutes}:${secs.toString().padStart(2, '0')}`;

        el.textContent = timeString;

        if (tot <= 0) {
            clearInterval(timer);
            el.textContent = 'Expired';
            el.style.color = '#ff4444';
        }
        tot--;

    }, 1000);
    return timer;
}
function extractFileId(responseText) {
    // Backend returns: "File uploaded--Download:/file/{id}\n"
    const match = responseText.match(/\/file\/([a-z0-9-]+)/);
    return match ? match[1] : null;
}

// Download functionality
async function handleDownload(event) {
    event.preventDefault();
    
    const fileId = document.getElementById('fileId').value.trim();
    const password = document.getElementById('downloadPassword').value;
    
    if (!fileId) {
        showToast('Please enter a file ID', 'error');
        return;
    }
    
    downloadBtn.disabled = true;
    downloadBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Downloading...';
    
    try {
        let url = `${API_BASE_URL}/file/${fileId}`;
        if (password) {
            url += `?password=${encodeURIComponent(password)}`;
        }
        
        const response = await fetch(url);
        
        if (!response.ok) {
            if (response.status === 404) {
                throw new Error('File not found or expired');
            } else if (response.status === 403) {
                throw new Error('Wrong password');
            } else if (response.status === 410) {
                throw new Error('No downloads left');
            } else {
                throw new Error(`Download failed: ${response.statusText}`);
            }
        }
        
        const contentDisposition = response.headers.get('Content-Disposition');
        let filename = 'download';
        if (contentDisposition) {
            const filenameMatch = contentDisposition.match(/filename="(.+)"/);
            if (filenameMatch) {
                filename = filenameMatch[1];
            }
        }
        const blob = await response.blob();
        const downloadUrl = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = downloadUrl;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(downloadUrl);
        
        showToast('File downloaded successfully!', 'success');
        
    } catch (error) {
        console.error('Download error:', error);
        showToast(`Download failed: ${error.message}`, 'error');
    } finally {
        // Reset button state
        downloadBtn.disabled = false;
        downloadBtn.innerHTML = '<i class="fas fa-download"></i> Download File';
    }
}

// UI Navigation
function showUploadSection() {
    uploadSection.style.display = 'block';
    successSection.style.display = 'none';
    downloadSection.style.display = 'none';
    resetUploadForm();
}

function showSuccessSection(file, downloads, expiry, password) {
    uploadSection.style.display = 'none';
    successSection.style.display = 'block';
    downloadSection.style.display = 'none';
    
    // Update success section content
    document.getElementById('fileName').textContent = file.name;
    document.getElementById('fileSize').textContent = formatFileSize(file.size);
    document.getElementById('downloadsLeft').textContent = downloads;
    
    // Update expiry display
    const expiryText = expiry === '1' ? '1 minute' : `${expiry} minutes`;
    document.querySelector('.download-info .info-item:last-child span strong').textContent = expiryText;
    
    const shareLink = `${window.location.origin}/download.html?id=${currentFileId}`;
    document.getElementById('shareLink').value = shareLink;
    
    // Display the file ID
    document.getElementById('fileIdDisplay').value = currentFileId;

    // Generate QR Code
    const qrcodeContainer = document.getElementById('qrcode');
    qrcodeContainer.innerHTML = ''; // Clear previous QR
    new QRCode(qrcodeContainer, {
        text: shareLink,
        width: 200,
        height: 200,
        colorDark: '#000000',
        colorLight: '#ffffff'
    });
    // Calculate expiry time and start countdown
const expiryMinutes = parseInt(expiry);
const expiryTime = new Date(Date.now() + expiryMinutes * 60 * 1000);

function updateCountdown() {
    const now = new Date();
    const remaining = expiryTime - now;
    
    if (remaining <= 0) {
        document.getElementById('expiryCountdown').textContent = 'Expired';
        document.getElementById('expiryCountdown').style.color = '#ff4444';
        return;
    }
    
    const hours = Math.floor(remaining / (1000 * 60 * 60));
    const minutes = Math.floor((remaining % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((remaining % (1000 * 60)) / 1000);
    
    const formatted = hours > 0 
        ? `${hours}h ${minutes}m ${seconds}s`
        : `${minutes}m ${seconds}s`;
    
    document.getElementById('expiryCountdown').textContent = formatted;
}

// Update countdown every second
updateCountdown();
const countdownInterval = setInterval(updateCountdown, 1000);

// Set upload time
const now = new Date();
const timeStr = now.toLocaleTimeString('en-US', { 
    hour: '2-digit', 
    minute: '2-digit' 
});
document.getElementById('uploadTime').textContent = timeStr;
}

function showDownloadSection() {
    uploadSection.style.display = 'none';
    successSection.style.display = 'none';
    downloadSection.style.display = 'block';
    resetDownloadForm();
}

function resetUploadForm() {
    uploadForm.reset();
    fileInputLabel.innerHTML = `
        <i class="fas fa-plus"></i>
        <span>Choose file or drag & drop</span>
    `;
    currentFileId = null;
    currentPassword = '';
    hidePreview();

    document.getElementById('qrcode').innerHTML = '';
}

// Hide preview function
function hidePreview() {
    const previewContainer = document.getElementById('filePreviewContainer');
    previewContainer.style.display = 'none';
}

function resetDownloadForm() {
    downloadForm.reset();
}

// Utility functions
function copyToClipboard() {
    const shareLink = document.getElementById('shareLink');
    shareLink.select();
    shareLink.setSelectionRange(0, 99999); // For mobile devices
    
    try {
        document.execCommand('copy');
        showToast('Link copied to clipboard!', 'success');
    } catch (err) {
        // Fallback for modern browsers
        navigator.clipboard.writeText(shareLink.value).then(() => {
            showToast('Link copied to clipboard!', 'success');
        }).catch(() => {
            showToast('Failed to copy link', 'error');
        });
    }
}

function copyFileId() {
    const fileIdInput = document.getElementById('fileIdDisplay');
    fileIdInput.select();
    fileIdInput.setSelectionRange(0, 99999); // For mobile devices
    
    try {
        document.execCommand('copy');
        showToast('File ID copied to clipboard!', 'success');
    } catch (err) {
        // Fallback for modern browsers
        navigator.clipboard.writeText(fileIdInput.value).then(() => {
            showToast('File ID copied to clipboard!', 'success');
        }).catch(() => {
            showToast('Failed to copy file ID', 'error');
        });
    }
}

function downloadQRCode() {
    const canvas = document.querySelector('#qrcode canvas');
    if (canvas) {
        const url = canvas.toDataURL('image/png');
        const a = document.createElement('a');
        a.href = url;
        a.download = `qr-code-${currentFileId}.png`;
        a.click();
        showToast('QR Code downloaded!', 'success');
    }
}

function testDownload() {
    if (currentFileId) {
        // Pre-fill the download form and switch to download section
        document.getElementById('fileId').value = currentFileId;
        if (currentPassword) {
            document.getElementById('downloadPassword').value = currentPassword;
        }
        showDownloadSection();
        showToast('Download form pre-filled with your file details', 'success');
    }
}

// Toast notifications
function showToast(message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    const icon = type === 'success' ? 'fa-check-circle' : 
                 type === 'error' ? 'fa-exclamation-circle' : 
                 'fa-exclamation-triangle';
    
    toast.innerHTML = `
        <i class="fas ${icon}"></i>
        <span>${message}</span>
    `;
    
    toastContainer.appendChild(toast);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    }, 5000);
    
    // Click to dismiss
    toast.addEventListener('click', () => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    });
}

// Handle URL parameters for direct download links
function handleUrlParameters() {
    const urlParams = new URLSearchParams(window.location.search);
    const fileId = urlParams.get('id');
    const password = urlParams.get('password');
    
    if (fileId) {
        // Pre-fill download form
        document.getElementById('fileId').value = fileId;
        if (password) {
            document.getElementById('downloadPassword').value = password;
        }
        showDownloadSection();
        showToast('Download link detected. Form pre-filled.', 'success');
    }
}

// Initialize URL parameter handling
document.addEventListener('DOMContentLoaded', handleUrlParameters);

// Add some visual feedback for better UX
function addVisualFeedback() {
    // Add loading states to buttons
    const buttons = document.querySelectorAll('button');
    buttons.forEach(button => {
        button.addEventListener('click', function() {
            if (!this.disabled) {
                this.style.transform = 'scale(0.98)';
                setTimeout(() => {
                    this.style.transform = '';
                }, 150);
            }
        });
    });
    
    // Add focus styles for better accessibility
    const inputs = document.querySelectorAll('input, textarea');
    inputs.forEach(input => {
        input.addEventListener('focus', function() {
            this.parentElement.style.transform = 'translateY(-2px)';
        });
        
        input.addEventListener('blur', function() {
            this.parentElement.style.transform = '';
        });
    });
}

function showProgress() {
    const progressContainer = document.getElementById('progressContainer');
    const progressFill = document.getElementById('progressFill');
    const progressText = document.getElementById('progressText');

    progressContainer.style.display = 'block';
    let progress = 0;
    const interval = setInterval(() => {
        progress += Math.random() * 10;
        if (progress > 100) progress = 100;
        
        progressFill.style.width = `${progress}%`;
        progressText.textContent = `${Math.floor(progress)}%`;

        if (progress >= 100) {
            clearInterval(interval);
            setTimeout(() => {
                progressContainer.style.display = 'none';
            }, 1000);
        }
    }, 200);
}

function validateFile(file) {
    const maxSize = 50 * 1024 * 1024; // Increased to 50MB

    if (!file) {
        throw new Error('No file selected');
    }

    if (file.size > maxSize) {
        throw new Error('File size must be less than 50MB');
    }
    
    // Allow all file types for flexibility
    return true;
}

function previewFile(file) {
    const previewContainer = document.getElementById('filePreviewContainer');
    const previewContent = document.getElementById('previewContent');

    previewContainer.style.display = 'block';
    
    
    previewContent.innerHTML = '';
    

    const fileType = file.type;
    const fileName = file.name;
    const fileSize = formatFileSize(file.size);
    
    if (fileType.startsWith('image/')) {
        // Image preview
        const reader = new FileReader();
        reader.onload = (e) => {
            previewContent.innerHTML = `
                <div class="preview-image">
                    <img src="${e.target.result}" alt="${fileName}" style="max-width: 100%; max-height: 300px; border-radius: 8px;">
                    <div class="preview-info">
                        <p><strong>${fileName}</strong></p>
                        <p>Size: ${fileSize}</p>
                        <p>Type: ${fileType}</p>
                    </div>
                </div>
            `;
        };
        reader.readAsDataURL(file);
    } else if (fileType.startsWith('text/') || fileName.endsWith('.txt') || fileName.endsWith('.md') || fileName.endsWith('.json')) {
        // Text file preview
        const reader = new FileReader();
        reader.onload = (e) => {
            const text = e.target.result;
            const previewText = text.length > 500 ? text.substring(0, 500) + '...' : text;
            previewContent.innerHTML = `
                <div class="preview-text">
                    <div class="preview-info">
                        <p><strong>${fileName}</strong></p>
                        <p>Size: ${fileSize}</p>
                        <p>Type: ${fileType}</p>
                    </div>
                    <div class="text-content">
                        <pre>${previewText}</pre>
                    </div>
                </div>
            `;
        };
        reader.readAsText(file);
    } else {
        // Generic file preview with icon
        const fileIcon = getFileIcon(fileType, fileName);
        previewContent.innerHTML = `
            <div class="preview-generic">
                <div class="file-icon-large">
                    <i class="${fileIcon}"></i>
                </div>
                <div class="preview-info">
                    <p><strong>${fileName}</strong></p>
                    <p>Size: ${fileSize}</p>
                    <p>Type: ${fileType || 'Unknown'}</p>
                </div>
            </div>
        `;
    }
}


function getFileIcon(fileType, fileName) {
    if (fileType.startsWith('image/')) {
        return 'fas fa-image';
    } 
    // Videos
    else if (fileType.startsWith('video/')) {
        return 'fas fa-video';
    } 
    // Audio
    else if (fileType.startsWith('audio/')) {
        return 'fas fa-music';
    } 
    // Text files
    else if (fileType.startsWith('text/') || fileName.endsWith('.txt') || fileName.endsWith('.md')) {
        return 'fas fa-file-alt';
    } 
    // Documents
    else if (fileName.endsWith('.pdf')) {
        return 'fas fa-file-pdf';
    } 
    else if (fileName.endsWith('.doc') || fileName.endsWith('.docx')) {
        return 'fas fa-file-word';
    } 
    else if (fileName.endsWith('.xls') || fileName.endsWith('.xlsx')) {
        return 'fas fa-file-excel';
    } 
    else if (fileName.endsWith('.ppt') || fileName.endsWith('.pptx')) {
        return 'fas fa-file-powerpoint';
    } 
    // Archives
    else if (fileName.endsWith('.zip') || fileName.endsWith('.rar') || fileName.endsWith('.7z') || fileName.endsWith('.tar') || fileName.endsWith('.gz')) {
        return 'fas fa-file-archive';
    } 
    // Code files
    else if (fileName.endsWith('.json') || fileName.endsWith('.xml') || fileName.endsWith('.html') || fileName.endsWith('.css') || fileName.endsWith('.js') || fileName.endsWith('.py') || fileName.endsWith('.go') || fileName.endsWith('.java') || fileName.endsWith('.cpp')) {
        return 'fas fa-file-code';
    }
    // Executables
    else if (fileName.endsWith('.exe') || fileName.endsWith('.msi') || fileName.endsWith('.dmg') || fileName.endsWith('.apk')) {
        return 'fas fa-cog';
    }
    // Default
    else {
        return 'fas fa-file';
    }
}


function uploadWithProgress(formData) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();

        xhr.upload.onprogress = (e) => {
            if (e.lengthComputable) {
                const percent = (e.loaded / e.total) * 100;
                console.log(`Upload progress: ${percent}%`);
            }
        };
        xhr.onload = () => {
            if (xhr.status >= 200 && xhr.status < 300) {
                resolve(xhr.responseText);
            } else {
                reject(new Error(`Server error: ${xhr.status} - ${xhr.responseText}`));
            }
        };
        xhr.onerror = () => reject(new Error('Network error - check if server is running'));

        xhr.open('POST', `${API_BASE_URL}/upload`);
        xhr.send(formData);
    });
}

document.addEventListener('DOMContentLoaded', addVisualFeedback);


function toggleDarkMode() {
    const darkModeToggle = document.getElementById('darkModeToggle');
    const body = document.body;
    
    if (darkModeToggle.checked) {
        body.classList.add('dark-mode');
        localStorage.setItem('darkMode', 'enabled');
        showToast('Dark mode enabled', 'success');
    } else {
        body.classList.remove('dark-mode');
        localStorage.setItem('darkMode', 'disabled');
        showToast('Light mode enabled', 'success');
    }
}

function initializeDarkMode() {
    const darkModeToggle = document.getElementById('darkModeToggle');
    const body = document.body;
    

    const darkModeEnabled = localStorage.getItem('darkMode') === 'enabled';
    
    if (darkModeEnabled) {
        body.classList.add('dark-mode');
        darkModeToggle.checked = true;
    } else {
        body.classList.remove('dark-mode');
        darkModeToggle.checked = false;
    }
}

// Add keyboard shortcuts
document.addEventListener('keydown', function(event) {
    // Ctrl/Cmd + U for upload mode
    if ((event.ctrlKey || event.metaKey) && event.key === 'u') {
        event.preventDefault();
        showUploadSection();
    }
    
    // Ctrl/Cmd + D for download mode
    if ((event.ctrlKey || event.metaKey) && event.key === 'd') {
        event.preventDefault();
        showDownloadSection();
    }
    
    // Escape to reset forms
    if (event.key === 'Escape') {
        if (successSection.style.display === 'block') {
            showUploadSection();
        } else if (downloadSection.style.display === 'block') {
            showUploadSection();
        }
    }
});

// Service worker disabled - uncomment if you add sw.js
// if ('serviceWorker' in navigator) {
//     window.addEventListener('load', function() {
//         navigator.serviceWorker.register('/sw.js')
//             .then(function(registration) {
//                 console.log('ServiceWorker registration successful');
//             })
//             .catch(function(err) {
//                 console.log('ServiceWorker registration failed');
//             });
//     });
// }
