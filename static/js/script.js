document.addEventListener('DOMContentLoaded', function() {
    // DOM Elements
    const showLoginBtn = document.getElementById('show-login');
    const showRegisterBtn = document.getElementById('show-register');
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const loggedOutView = document.getElementById('logged-out-view');
    const loggedInView = document.getElementById('logged-in-view');
    const userEmail = document.getElementById('user-email');
    const logoutBtn = document.getElementById('logout-btn');
    const fileManagement = document.getElementById('file-management');
    const uploadForm = document.getElementById('upload-form');
    const filesList = document.getElementById('files-list');
    const messageDiv = document.getElementById('message');

    // API Endpoints
    const API_URL = 'http://localhost:8080';
    const ENDPOINTS = {
        REGISTER: `${API_URL}/register`,
        LOGIN: `${API_URL}/login`,
        UPLOAD: `${API_URL}/upload`,
        FILES: `${API_URL}/files`,
        FILE: (id) => `${API_URL}/files/${id}`,
        SHARE: (id) => `${API_URL}/share/${id}`
    };

    // Check if user is logged in
    function checkAuth() {
        const token = localStorage.getItem('token');
        if (token) {
            loggedOutView.classList.add('hidden');
            loggedInView.classList.remove('hidden');
            fileManagement.classList.remove('hidden');
            fetchFiles();
        } else {
            loggedOutView.classList.remove('hidden');
            loggedInView.classList.add('hidden');
            fileManagement.classList.add('hidden');
        }
    }

    // Show message
    function showMessage(text, isError = false) {
        messageDiv.textContent = text;
        messageDiv.className = isError ? 'error' : 'success';
        messageDiv.classList.remove('hidden');
        setTimeout(() => {
            messageDiv.classList.add('hidden');
        }, 3000);
    }

    // Register user
    async function registerUser(email, password) {
        try {
            const response = await fetch(ENDPOINTS.REGISTER, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Registration failed');
            }

            showMessage('Registration successful! Please login.');
            loginForm.classList.remove('hidden');
            registerForm.classList.add('hidden');
        } catch (error) {
            showMessage(error.message, true);
        }
    }

    // Login user
    async function loginUser(email, password) {
        try {
            const response = await fetch(ENDPOINTS.LOGIN, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Login failed');
            }

            localStorage.setItem('token', data.token);
            localStorage.setItem('userEmail', email);
            userEmail.textContent = email;
            showMessage('Login successful!');
            checkAuth();
        } catch (error) {
            showMessage(error.message, true);
        }
    }

    // Fetch user files
    async function fetchFiles() {
        try {
            const token = localStorage.getItem('token');
            if (!token) return;

            const response = await fetch(ENDPOINTS.FILES, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Failed to fetch files');
            }

            renderFiles(data.files || []);
        } catch (error) {
            showMessage(error.message, true);
        }
    }

    // Upload file
    async function uploadFile(file) {
        try {
            const token = localStorage.getItem('token');
            if (!token) return;

            const formData = new FormData();
            formData.append('file', file);

            const response = await fetch(ENDPOINTS.UPLOAD, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Upload failed');
            }

            showMessage('File uploaded successfully!');
            fetchFiles();
        } catch (error) {
            showMessage(error.message, true);
        }
    }

    // Delete file
    async function deleteFile(fileId) {
        try {
            const token = localStorage.getItem('token');
            if (!token) return;

            const response = await fetch(ENDPOINTS.FILE(fileId), {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.error || 'Delete failed');
            }

            showMessage('File deleted successfully!');
            fetchFiles();
        } catch (error) {
            showMessage(error.message, true);
        }
    }

    // Share file
async function shareFile(fileId) {
    try {
        const token = localStorage.getItem('token');
        if (!token) return;

        const response = await fetch(ENDPOINTS.SHARE(fileId), {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Share failed');
        }

        const shareUrl = `${API_URL}${data.url}`;
        navigator.clipboard.writeText(shareUrl);
        showMessage('Share link copied to clipboard!');
    } catch (error) {
        showMessage(error.message, true);
    }
}

function getFile(fileId) {
    const token = localStorage.getItem('token');
    if (!token) {
        showMessage('Please log in to download files.', true);
        return;
    }

    window.open(`${API_URL}/files/${fileId}`, '_blank');
}



    // Render files list
    function renderFiles(files) {
        filesList.innerHTML = '';
        
        if (files.length === 0) {
            filesList.innerHTML = '<p>No files uploaded yet.</p>';
            return;
        }

        files.forEach(file => {
            const fileItem = document.createElement('div');
            fileItem.className = 'file-item';
            
            const fileInfo = document.createElement('div');
            fileInfo.className = 'file-info';
            fileInfo.innerHTML = `
                <p><strong>${file.original_filename}</strong></p>
                <p>Size: ${formatFileSize(file.file_size)}</p>
            `;
            
            const fileActions = document.createElement('div');
            fileActions.className = 'file-actions';
            
            const downloadBtn = document.createElement('button');
            downloadBtn.textContent = 'Download';
            downloadBtn.addEventListener('click', () => {
                window.open(ENDPOINTS.FILE(file.id), '_blank');
            });
            
            const shareBtn = document.createElement('button');
            shareBtn.textContent = 'Share';
            shareBtn.className = 'share-btn';
            shareBtn.addEventListener('click', () => shareFile(file.id));
            
            const deleteBtn = document.createElement('button');
            deleteBtn.textContent = 'Delete';
            deleteBtn.className = 'delete-btn';
            deleteBtn.addEventListener('click', () => deleteFile(file.id));
            
            fileActions.appendChild(downloadBtn);
            fileActions.appendChild(shareBtn);
            fileActions.appendChild(deleteBtn);
            
            fileItem.appendChild(fileInfo);
            fileItem.appendChild(fileActions);
            filesList.appendChild(fileItem);
        });
    }

    // Format file size
    function formatFileSize(bytes) {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1048576) return (bytes / 1024).toFixed(2) + ' KB';
        if (bytes < 1073741824) return (bytes / 1048576).toFixed(2) + ' MB';
        return (bytes / 1073741824).toFixed(2) + ' GB';
    }

    // Event Listeners
    showLoginBtn.addEventListener('click', () => {
        loginForm.classList.remove('hidden');
        registerForm.classList.add('hidden');
    });

    showRegisterBtn.addEventListener('click', () => {
        registerForm.classList.remove('hidden');
        loginForm.classList.add('hidden');
    });

    document.getElementById('login').addEventListener('submit', (e) => {
        e.preventDefault();
        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;
        loginUser(email, password);
    });

    document.getElementById('register').addEventListener('submit', (e) => {
        e.preventDefault();
        const email = document.getElementById('register-email').value;
        const password = document.getElementById('register-password').value;
        registerUser(email, password);
    });

    logoutBtn.addEventListener('click', () => {
        localStorage.removeItem('token');
        localStorage.removeItem('userEmail');
        checkAuth();
        showMessage('Logged out successfully!');
    });

    uploadForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const fileInput = document.getElementById('file-input');
        if (fileInput.files.length > 0) {
            uploadFile(fileInput.files[0]);
            fileInput.value = '';
        }
    });

    // Initialize
    checkAuth();
    userEmail.textContent = localStorage.getItem('userEmail') || '';
});
