document.addEventListener('DOMContentLoaded', function() {
    // Tab switching functionality
    const tabs = document.querySelectorAll('.tab');
    const tabContents = document.querySelectorAll('.tab-content');
    
    tabs.forEach(tab => {
        tab.addEventListener('click', function() {
            const targetTab = this.getAttribute('data-tab');
            
            // Update active tab
            tabs.forEach(t => t.classList.remove('active'));
            this.classList.add('active');
            
            // Update active content
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(`${targetTab}-tab`).classList.add('active');

            // Start polling for QR login status if QR tab is active
            if (targetTab === 'qr') {
                startQrPolling();
            } else {
                stopQrPolling();
            }
        });
    });
    
    // Secret key visibility toggle
    const toggleKey = document.getElementById('toggle-key');
    const secretKeyInput = document.getElementById('secretkey');
    
    if (toggleKey && secretKeyInput) {
        toggleKey.addEventListener('click', function() {
            const type = secretKeyInput.getAttribute('type') === 'password' ? 'text' : 'password';
            secretKeyInput.setAttribute('type', type);
            
            // Change the eye icon
            this.innerHTML = type === 'password' ? 
                '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>' : 
                '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>';
        });
    }
    
    // Form validation
    const loginForm = document.getElementById('login-form');
    const usernameInput = document.getElementById('username');
    const usernameError = document.getElementById('username-error');
    const secretKeyError = document.getElementById('secretkey-error');
    
    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            let hasError = false;
            
            // Reset error messages
            if (usernameError) {
                usernameError.textContent = '';
                usernameError.classList.remove('active');
            }
            
            if (secretKeyError) {
                secretKeyError.textContent = '';
                secretKeyError.classList.remove('active');
            }
            
            // Validate username
            if (usernameInput && !usernameInput.value.trim()) {
                if (usernameError) {
                    usernameError.textContent = 'Please enter your username';
                    usernameError.classList.add('active');
                }
                hasError = true;
            }
            
            // Validate secret key
            if (secretKeyInput && !secretKeyInput.value.trim()) {
                if (secretKeyError) {
                    secretKeyError.textContent = 'Please enter your secret key';
                    secretKeyError.classList.add('active');
                }
                hasError = true;
            }
            
            if (hasError) {
                e.preventDefault();
            }
        });
    }
    
    // QR code timer functionality
    const qrTimer = document.getElementById('qr-timer');
    let timeLeft = 300; // 5 minutes in seconds
    let timerInterval;
    
    function updateTimer() {
        if (qrTimer) {
            const minutes = Math.floor(timeLeft / 60);
            const seconds = timeLeft % 60;
            qrTimer.textContent = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
            
            if (timeLeft <= 0) {
                clearInterval(timerInterval);
                qrTimer.textContent = "Expired";
                qrTimer.style.color = "var(--error-color)";
                
                // Update status
                const qrStatus = document.getElementById('qr-status');
                if (qrStatus) {
                    qrStatus.innerHTML = '<p class="expired">QR code expired. Please generate a new one.</p>';
                }
                
                stopQrPolling();
            } else {
                timeLeft--;
            }
        }
    }
    
    if (qrTimer) {
        updateTimer(); // Initial call
        timerInterval = setInterval(updateTimer, 1000);
    }
    
    // QR code refresh button
    const refreshQrBtn = document.getElementById('refresh-qr');
    
    if (refreshQrBtn) {
        refreshQrBtn.addEventListener('click', function() {
            // Request a new QR code from the server
            fetch('/api/generate-login-qr?refresh=true')
                .then(response => {
                    if (response.ok) {
                        return response.json();
                    }
                    throw new Error('Failed to refresh QR code');
                })
                .then(data => {
                    // Update the QR code image
                    const qrImage = document.querySelector('.qr-code img');
                    if (qrImage) {
                        qrImage.src = data.qrUrl + '?t=' + new Date().getTime(); // Add timestamp to prevent caching
                        
                        // Reset timer
                        clearInterval(timerInterval);
                        timeLeft = 300;
                        qrTimer.style.color = "var(--primary-color)";
                        timerInterval = setInterval(updateTimer, 1000);
                        
                        // Reset status
                        const qrStatus = document.getElementById('qr-status');
                        if (qrStatus) {
                            qrStatus.innerHTML = '<p>Waiting for scan...</p>';
                        }
                        
                        // Start polling
                        qrToken = data.token;
                        startQrPolling();
                    }
                })
                .catch(error => {
                    console.error('Error refreshing QR code:', error);
                });
        });
    }
    
    // QR login polling
    let qrToken = null;
    let pollingInterval = null;
    
    // Extract token from QR code URL if present
    const qrImage = document.querySelector('.qr-code img');
    if (qrImage && qrImage.src) {
        const urlParts = qrImage.src.split('/');
        const filenameParts = urlParts[urlParts.length - 1].split('.');
        if (filenameParts.length > 0) {
            qrToken = filenameParts[0];
        }
    }
    
    function startQrPolling() {
        if (!qrToken) return;
        
        // Stop any existing polling
        stopQrPolling();
        
        // Start new polling
        pollingInterval = setInterval(() => {
            fetch(`/api/check-qr-login?token=${qrToken}`)
                .then(response => response.json())
                .then(data => {
                    if (data.status === 'scanned') {
                        // Update UI to show scanned status
                        const qrStatus = document.getElementById('qr-status');
                        if (qrStatus) {
                            qrStatus.innerHTML = '<p class="scanned">QR code scanned! Waiting for confirmation...</p>';
                        }
                    } else if (data.status === 'authenticated') {
                        // QR code was used for successful authentication
                        stopQrPolling();
                        
                        // Update UI
                        const qrStatus = document.getElementById('qr-status');
                        if (qrStatus) {
                            qrStatus.innerHTML = '<p class="success">Authentication successful! Redirecting...</p>';
                        }
                        
                        // Redirect to the dashboard or wherever the server indicates
                        setTimeout(() => {
                            window.location.href = data.redirectUrl || '/dashboard';
                        }, 1000);
                    }
                })
                .catch(error => {
                    console.error('Error checking QR login status:', error);
                });
        }, 2000); // Check every 2 seconds
    }
    
    function stopQrPolling() {
        if (pollingInterval) {
            clearInterval(pollingInterval);
            pollingInterval = null;
        }
    }
    
    // Start polling if we're on the QR tab
    if (document.querySelector('.tab[data-tab="qr"]').classList.contains('active')) {
        startQrPolling();
    }
});