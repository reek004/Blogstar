// API URL configuration
const API_URL = 'http://localhost:8080';

// DOM elements
const elements = {
    form: {
        contentType: document.getElementById('contentType'),
        topic: document.getElementById('topic'),
        topicError: document.getElementById('topicError'),
        tone: document.getElementById('tone'),
        length: document.getElementById('length'),
        additionalContext: document.getElementById('additionalContext'),
        generateBtn: document.getElementById('generateBtn'),
        loading: document.getElementById('loading')
    },
    status: document.getElementById('status'),
    tabs: document.querySelectorAll('.tab'),
    tabContents: document.querySelectorAll('.tab-content'),
    result: document.getElementById('result'),
    jsonResponse: document.getElementById('jsonResponse')
};

// Show a status message to the user
function showStatus(message, isError = false) {
    elements.status.textContent = message;
    elements.status.style.display = 'block';
    
    if (isError) {
        elements.status.className = 'error-msg';
    } else {
        elements.status.className = 'success';
    }
    
    // Auto-hide after 5 seconds
    setTimeout(() => {
        elements.status.style.display = 'none';
    }, 5000);
}

// Switch between tabs
function switchTab(tabName) {
    // Hide all tabs and remove active class
    elements.tabContents.forEach(tab => tab.classList.remove('active'));
    elements.tabs.forEach(tab => tab.classList.remove('active'));
    
    // Show the selected tab and add active class
    document.getElementById(tabName + 'Tab').classList.add('active');
    elements.tabs.forEach(tab => {
        if (tab.dataset.tab === tabName) {
            tab.classList.add('active');
        }
    });
}

// Validate form fields
function validateForm() {
    let isValid = true;
    const topic = elements.form.topic.value.trim();
    
    if (!topic) {
        elements.form.topicError.textContent = 'Topic is required';
        isValid = false;
    } else {
        elements.form.topicError.textContent = '';
    }
    
    return isValid;
}

// Event handler for generate button
async function generateContent(event) {
    console.log("Generate content function called");
    
    // Prevent any default browser behavior and propagation
    if (event) {
        event.preventDefault();
        event.stopPropagation();
    }
    
    // Validate form
    if (!validateForm()) {
        return false;
    }
    
    // Show loading and disable button
    elements.form.loading.style.display = 'inline';
    elements.form.generateBtn.disabled = true;
    
    // Clear previous results to prevent confusion
    elements.result.textContent = 'Generating content...';
    elements.jsonResponse.textContent = 'Waiting for response...';
    
    const requestData = {
        content_type: elements.form.contentType.value,
        topic: elements.form.topic.value.trim(),
        tone: elements.form.tone.value.trim(),
        length: parseInt(elements.form.length.value || '500'),
        additional_context: elements.form.additionalContext.value.trim()
    };
    
    try {
        console.log('Making request to API:', requestData);
        
        const response = await fetch(`${API_URL}/api/generate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        });
        
        console.log('Response status:', response.status);
        
        // Handle rate limiting (429 Too Many Requests)
        if (response.status === 429) {
            const retryAfter = response.headers.get('Retry-After') || 60;
            const errorMessage = `Rate limit exceeded. Please try again after ${retryAfter} seconds.`;
            elements.result.textContent = `Error: ${errorMessage}`;
            elements.jsonResponse.textContent = `Error: ${errorMessage}`;
            showStatus(errorMessage, true);
            return false;
        }
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        // Parse the response
        const data = await response.json();
        console.log('Response data:', data);
        
        // Display the full JSON in the response tab
        elements.jsonResponse.textContent = JSON.stringify(data, null, 2);
        
        // Handle the actual content display
        if (data.error) {
            elements.result.textContent = `Error: ${data.error}`;
            showStatus(`Error: ${data.error}`, true);
        } else {
            // Format and display the content
            elements.result.textContent = data.content;
            showStatus(`Content generated successfully and saved to: ${data.filename}`);
            
            // Switch to preview tab
            switchTab('preview');
        }
    } catch (error) {
        console.error('Error:', error);
        elements.result.textContent = `Error: ${error.message}`;
        elements.jsonResponse.textContent = `Error: ${error.message}`;
        showStatus(`Error: ${error.message}`, true);
    } finally {
        // Hide loading and enable button
        elements.form.loading.style.display = 'none';
        elements.form.generateBtn.disabled = false;
    }
    
    return false;
}

// Health check function to test API connection
async function checkApiHealth() {
    try {
        const response = await fetch(`${API_URL}/health`);
        if (response.ok) {
            console.log('API connection successful');
        } else {
            console.warn('API health check failed:', response.status);
            showStatus('Warning: API server appears to be offline. Please ensure the server is running.', true);
        }
    } catch (error) {
        console.error('API connection error:', error);
        showStatus('Error: Cannot connect to API server. Please ensure the server is running.', true);
    }
}

// Set up event listeners
function initializeApp() {
    // Button click handler
    elements.form.generateBtn.addEventListener('click', generateContent);
    
    // Tab switching
    elements.tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            switchTab(tab.dataset.tab);
        });
    });
    
    // Prevent form submission on enter key in inputs
    const inputs = document.querySelectorAll('input, textarea');
    inputs.forEach(input => {
        input.addEventListener('keydown', function(e) {
            if (e.key === 'Enter') {
                e.preventDefault();
                return false;
            }
        });
    });
    
    // Check API health when page loads
    checkApiHealth();
    
    console.log('MarloweQuill client initialized');
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', initializeApp);

// Prevent accidental page unload if content has been generated
window.addEventListener('beforeunload', function(e) {
    const resultText = elements.result.textContent;
    if (resultText && 
        resultText !== 'Results will appear here...' && 
        resultText !== 'Generating content...') {
        e.preventDefault();
        e.returnValue = 'You have generated content that will be lost if you leave the page.';
        return e.returnValue;
    }
});
