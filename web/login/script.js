// DOM Elements
const signInForm = document.getElementById('signInForm');
const signUpForm = document.getElementById('signUpForm');
const toSignUp = document.getElementById('toSignUp');
const toSignIn = document.getElementById('toSignIn');

// Error message elements for Sign In
const signInEmailError = document.getElementById('signInEmailError');
const signInPasswordError = document.getElementById('signInPasswordError');

// Error message elements for Sign Up
const signUpNameError = document.getElementById('signUpNameError');
const signUpEmailError = document.getElementById('signUpEmailError');
const signUpPasswordError = document.getElementById('signUpPasswordError');

// Toggle between forms
toSignUp.addEventListener('click', () => {
    document.querySelector('.form-container').style.height = '600px';
    signInForm.classList.remove('active');
    setTimeout(() => {
        signUpForm.classList.add('active');
    }, 300);
});

toSignIn.addEventListener('click', () => {
    document.querySelector('.form-container').style.height = '500px';
    signUpForm.classList.remove('active');
    setTimeout(() => {
        signInForm.classList.add('active');
    }, 300);
});

// Sign In Validation
signInForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    // Reset message to empty
    setMessage('formMessage', '');
    
    // Reset error messages
    signInEmailError.textContent = '';
    signInPasswordError.textContent = '';
    
    const email = document.getElementById('signInEmail').value.trim();
    const password = document.getElementById('signInPassword').value;
    
    let isValid = true;
    
    // Email validation
    if (!email) {
        signInEmailError.textContent = 'Email is required';
        isValid = false;
    } else if (!isValidEmail(email)) {
        signInEmailError.textContent = 'Please enter a valid email address';
        isValid = false;
    }
    
    // Password validation
    if (!password) {
        signInPasswordError.textContent = 'Password is required';
        isValid = false;
    } else if (password.length < 6) {
        signInPasswordError.textContent = 'Password must be at least 6 characters';
        isValid = false;
    }
    
    if (!isValid) {
        // Set general error message
        setMessage('formMessage', 'Please correct the form errors.');
    } else {
        try {
            const response = await fetch('/signin', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('authToken', data.token);
                // Set success message
                setMessage('formMessage', 'Login successful! Redirecting...');
                window.location.href = '/drawboard'; // Redirect to the main app
            } else {
                // Set error message
                const errorText = await response.text();
                setMessage('formMessage', 'Login failed: ' + errorText);
            }
        } catch (error) {
            console.error('Login error:', error);
            setMessage('formMessage', 'An error occurred during login.');
        }
        signInForm.reset();
    }
});

// Sign Up Validation
signUpForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    // Reset message to empty
    setMessage('formMessage', '');
    
    // Reset error messages
    signUpNameError.textContent = '';
    signUpEmailError.textContent = '';
    signUpPasswordError.textContent = '';
    
    const name = document.getElementById('signUpName').value.trim();
    const email = document.getElementById('signUpEmail').value.trim();
    const password = document.getElementById('signUpPassword').value;
    
    let isValid = true;
    
    // Name validation
    if (!name) {
        signUpNameError.textContent = 'Name is required';
        isValid = false;
    } else if (name.length < 2) {
        signUpNameError.textContent = 'Name must be at least 2 characters';
        isValid = false;
    } else if (!/^[a-zA-Z\s]+$/.test(name)) {
        signUpNameError.textContent = 'Name can only contain letters and spaces';
        isValid = false;
    }
    
    // Email validation
    if (!email) {
        signUpEmailError.textContent = 'Email is required';
        isValid = false;
    } else if (!isValidEmail(email)) {
        signUpEmailError.textContent = 'Please enter a valid email address';
        isValid = false;
    }
    
    // Password validation
    if (!password) {
        signUpPasswordError.textContent = 'Password is required';
        isValid = false;
    } else if (password.length < 6) {
        signUpPasswordError.textContent = 'Password must be at least 6 characters';
        isValid = false;
    } else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(password)) {
        signUpPasswordError.textContent = 'Password must contain at least one uppercase letter, one lowercase letter, and one number';
        isValid = false;
    }
    
    if (!isValid) {
        // Set general error message
        setMessage('formMessage', 'Please correct the form errors.');
    } else {
        try {
            const response = await fetch('/signup', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, email, password }),
            });

            if (response.ok) {
                window.location.href = "/login/"
                event.target.reset();
            } else {
                const errorText = await response.text();
                setMessage('formMessage', 'Signup failed: ' + errorText);
            }
        } catch (error) {
            console.error('Signup error:', error);
            setMessage('formMessage', 'An error occurred during signup.');
        }
        signUpForm.reset();
        
        // Switch to sign in form after successful sign up
        setTimeout(() => {
            signUpForm.classList.remove('active');
            setTimeout(() => {
                signInForm.classList.add('active');
            }, 300);
        }, 500);
    }
});

// Helper function to validate email
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Real-time validation for better UX
document.getElementById('signInEmail').addEventListener('blur', function() {
    const email = this.value.trim();
    if (email && !isValidEmail(email)) {
        signInEmailError.textContent = 'Please enter a valid email address';
    } else {
        signInEmailError.textContent = '';
    }
});

document.getElementById('signUpName').addEventListener('blur', function() {
    const name = this.value.trim();
    if (name && name.length < 2) {
        signUpNameError.textContent = 'Name must be at least 2 characters';
    } else if (name && !/^[a-zA-Z\s]+$/.test(name)) {
        signUpNameError.textContent = 'Name can only contain letters and spaces';
    } else {
        signUpNameError.textContent = '';
    }
});

document.getElementById('signUpEmail').addEventListener('blur', function() {
    const email = this.value.trim();
    if (email && !isValidEmail(email)) {
        signUpEmailError.textContent = 'Please enter a valid email address';
    } else {
        signUpEmailError.textContent = '';
    }
});

document.getElementById('signUpPassword').addEventListener('blur', function() {
    const password = this.value;
    if (password && password.length < 6) {
        signUpPasswordError.textContent = 'Password must be at least 6 characters';
    } else if (password && !/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(password)) {
        signUpPasswordError.textContent = 'Password must contain at least one uppercase letter, one lowercase letter, and one number';
    } else {
        signUpPasswordError.textContent = '';
    }
});

// Add setMessage function
function setMessage(elementId, message) {
    const element = document.getElementById(elementId);
    if (element) {
        element.textContent = message;
    }
}
