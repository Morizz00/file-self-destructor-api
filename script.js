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
    
    // Navigation buttons
    uploadAnotherBtn.addEventListener('click', showUploadSection);
    testDownloadBtn.addEventListener('click', testDownload);
    uploadModeBtn.addEventListener('click', showUploadSection);
    downloadModeBtn.addEventListener('click', showDownloadSection);
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
    
    uploadBtn.disabled = true;
    uploadBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Uploading...';
    
    try {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('downloads', downloads);
        if (password) {
            formData.append('password', password);
        }
        
        const responseText = await uploadWithProgress(formData);
        const fileId = extractFileId(responseText);
        
        if (fileId) {
            currentFileId = fileId;
            currentPassword = password;
            showSuccessSection(file, downloads, password);
            showToast('File uploaded successfully!', 'success');
        } else {
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

function extractFileId(responseText) {
    const match = responseText.match(/\/file\/([a-f0-9]+)/);
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

function showSuccessSection(file, downloads, password) {
    uploadSection.style.display = 'none';
    successSection.style.display = 'block';
    downloadSection.style.display = 'none';
    
    // Update success section content
    document.getElementById('fileName').textContent = file.name;
    document.getElementById('fileSize').textContent = formatFileSize(file.size);
    document.getElementById('downloadsLeft').textContent = downloads;
    
    const shareLink = `${window.location.origin}/download.html?id=${currentFileId}`;
    document.getElementById('shareLink').value = shareLink;
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

function validateFile(file){
    const maxSize=10*1024*1024;
    const allowedTypes=['image/','text/','application/pdf'];

    if (!file){
        throw new Error('No file selected');
    }

    if(file.size>maxSize){
        throw new Error('File size must be less than 10MB');
    }
    
    if (!allowedTypes.some(type=>file.type.startsWith(type))){
        throw new Error('File type not allowed');
    }
    
    return true;
}

function previewFile(file){
    const reader=new FileReader();
    reader.onload=(e)=>{
        console.log('file preview ready');
    };
    reader.readAsDataURL(file);
}
function uploadWithProgress(formData){
    return new Promise((resolve,reject)=>{
        const xhr=new XMLHttpRequest();

        xhr.upload.onprogress=(e)=>{
            if (e.lengthComputable){
                const percent=(e.loaded/e.total)*100;
                console.log(`Upload progress: ${percent}%`);
            }
        };
        xhr.onload=()=>resolve(xhr.responseText);
        xhr.onerror=()=>reject(new Error('Upload failed'));

        xhr.open('POST', `${API_BASE_URL}/upload`);
        xhr.send(formData);
    });
}
// Initialize visual feedback
document.addEventListener('DOMContentLoaded', addVisualFeedback);

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

// Add service worker for offline functionality (optional)
if ('serviceWorker' in navigator) {
    window.addEventListener('load', function() {
        navigator.serviceWorker.register('/sw.js')
            .then(function(registration) {
                console.log('ServiceWorker registration successful');
            })
            .catch(function(err) {
                console.log('ServiceWorker registration failed');
            });
    });
}
