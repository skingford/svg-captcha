package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"svg-math-captcha/captcha"
)

// Server represents the HTTP server with captcha functionality
type Server struct {
	generator *captcha.CaptchaGenerator
	sessions  map[string]*Session
}

// Session stores captcha session data
type Session struct {
	Answer    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewServer creates a new server instance
func NewServer() *Server {
	config := &captcha.Config{
		MathMin:      1,
		MathMax:      10,
		MathOperator: "+",
		Width:        200,
		Height:       60,
		FontSize:     24,
		Noise:        2,
		Color:        true,
		Background:   "#f8f9fa",
	}

	return &Server{
		generator: captcha.NewCaptchaGenerator(config),
		sessions:  make(map[string]*Session),
	}
}

// generateCaptcha handles captcha generation requests
func (s *Server) generateCaptcha(w http.ResponseWriter, r *http.Request) {
	// Generate captcha
	result, err := s.generator.CreateMathExpr()
	if err != nil {
		log.Printf("Error generating captcha: %v", err)
		http.Error(w, "Failed to generate captcha", http.StatusInternalServerError)
		return
	}

	// Create session ID (in production, use a proper session management library)
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())

	// Store session
	s.sessions[sessionID] = &Session{
		Answer:    result.Text,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// Clean up expired sessions
	s.cleanupSessions()

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "captcha_session",
		Value:    sessionID,
		HttpOnly: true,
		MaxAge:   300, // 5 minutes
		Path:     "/",
	})

	// Return SVG
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result.Data))
}

// validateCaptcha handles captcha validation requests
func (s *Server) validateCaptcha(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Answer string `json:"answer"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get session from cookie
	cookie, err := r.Cookie("captcha_session")
	if err != nil {
		log.Printf("No session cookie found: %v", err)
		http.Error(w, "No captcha session found", http.StatusBadRequest)
		return
	}

	// Find session
	session, exists := s.sessions[cookie.Value]
	if !exists {
		http.Error(w, "Invalid or expired session", http.StatusBadRequest)
		return
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, cookie.Value)
		http.Error(w, "Captcha expired", http.StatusBadRequest)
		return
	}

	// Validate answer
	isValid := captcha.ValidateAnswer(session.Answer, request.Answer)

	response := struct {
		Valid   bool   `json:"valid"`
		Message string `json:"message"`
	}{
		Valid:   isValid,
		Message: "Validation result",
	}

	if isValid {
		response.Message = "Captcha validation successful"
		// Remove session after successful validation
		delete(s.sessions, cookie.Value)

		// Clear the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "captcha_session",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
			Path:     "/",
		})
	} else {
		response.Message = "Captcha validation failed"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// serveDemoPage serves the HTML demo page
func (s *Server) serveDemoPage(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SVG Math Captcha Demo</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .captcha-container {
            text-align: center;
            margin: 30px 0;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 8px;
            border: 1px solid #dee2e6;
        }
        .captcha-image {
            display: inline-block;
            border: 2px solid #ddd;
            border-radius: 4px;
            background: white;
            margin-bottom: 15px;
        }
        .form-group {
            margin: 20px 0;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }
        input[type="text"] {
            width: 200px;
            padding: 12px;
            font-size: 16px;
            border: 2px solid #ddd;
            border-radius: 4px;
            transition: border-color 0.3s;
        }
        input[type="text"]:focus {
            outline: none;
            border-color: #007bff;
        }
        button {
            padding: 12px 20px;
            font-size: 16px;
            margin: 5px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        .btn-primary {
            background-color: #007bff;
            color: white;
        }
        .btn-primary:hover {
            background-color: #0056b3;
        }
        .btn-secondary {
            background-color: #6c757d;
            color: white;
        }
        .btn-secondary:hover {
            background-color: #545b62;
        }
        .result {
            margin: 20px 0;
            padding: 15px;
            border-radius: 4px;
            font-weight: bold;
        }
        .success {
            background-color: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .error {
            background-color: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .info {
            background-color: #d1ecf1;
            color: #0c5460;
            border: 1px solid #bee5eb;
            margin-bottom: 20px;
        }
        .loading {
            opacity: 0.6;
            pointer-events: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔢 SVG Math Captcha Demo</h1>
        
        <div class="info">
            <strong>Info:</strong> This demo shows a Go implementation of math captcha similar to node svg-captcha library.
            Solve the math problem shown in the image below.
        </div>
        
        <div class="captcha-container">
            <div class="captcha-image">
                <img id="captcha" src="/captcha" alt="Math Captcha" onload="hideLoading()" onerror="showError()" />
            </div>
            <br>
            <button class="btn-secondary" onclick="refreshCaptcha()">🔄 Refresh Captcha</button>
        </div>
        
        <form id="captchaForm">
            <div class="form-group">
                <label for="answer">Enter your answer:</label>
                <input type="text" id="answer" name="answer" placeholder="Type the answer here" required>
            </div>
            <div class="form-group">
                <button type="submit" class="btn-primary">✓ Validate Answer</button>
            </div>
        </form>
        
        <div id="result"></div>
        
        <div id="loading" style="display: none;">
            <div class="info">Loading new captcha...</div>
        </div>
    </div>
    
    <script>
        function refreshCaptcha() {
            showLoading();
            document.getElementById('captcha').src = '/captcha?' + new Date().getTime();
            document.getElementById('answer').value = '';
            document.getElementById('result').innerHTML = '';
        }
        
        function showLoading() {
            document.getElementById('loading').style.display = 'block';
            document.querySelector('.captcha-container').classList.add('loading');
        }
        
        function hideLoading() {
            document.getElementById('loading').style.display = 'none';
            document.querySelector('.captcha-container').classList.remove('loading');
        }
        
        function showError() {
            hideLoading();
            document.getElementById('result').innerHTML = 
                '<div class="result error">❌ Failed to load captcha. Please try refreshing.</div>';
        }
        
        document.getElementById('captchaForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const answer = document.getElementById('answer').value.trim();
            const resultDiv = document.getElementById('result');
            
            if (!answer) {
                resultDiv.innerHTML = '<div class="result error">❌ Please enter an answer</div>';
                return;
            }
            
            try {
                const response = await fetch('/validate', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ answer: answer })
                });
                
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                
                const data = await response.json();
                
                if (data.valid) {
                    resultDiv.innerHTML = '<div class="result success">✅ ' + data.message + '</div>';
                    setTimeout(() => {
                        refreshCaptcha();
                    }, 2000);
                } else {
                    resultDiv.innerHTML = '<div class="result error">❌ ' + data.message + '</div>';
                }
            } catch (error) {
                console.error('Error:', error);
                resultDiv.innerHTML = '<div class="result error">❌ Network error. Please try again.</div>';
            }
        });
        
        // Focus on answer input when page loads
        document.addEventListener('DOMContentLoaded', function() {
            document.getElementById('answer').focus();
        });
        
        // Allow Enter key to submit form
        document.getElementById('answer').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                e.preventDefault();
                document.getElementById('captchaForm').dispatchEvent(new Event('submit'));
            }
        });
    </script>
</body>
</html>`

	t, err := template.New("demo").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, nil)
}

// apiStatus returns server status information
func (s *Server) apiStatus(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status         string          `json:"status"`
		Version        string          `json:"version"`
		ActiveSessions int             `json:"active_sessions"`
		Config         *captcha.Config `json:"config"`
	}{
		Status:         "ok",
		Version:        "1.0.0",
		ActiveSessions: len(s.sessions),
		Config:         s.generator.GetConfig(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// cleanupSessions removes expired sessions
func (s *Server) cleanupSessions() {
	now := time.Now()
	for sessionID, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, sessionID)
		}
	}
}

// CORS middleware for API endpoints
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	server := NewServer()

	// Routes
	http.HandleFunc("/", server.serveDemoPage)
	http.HandleFunc("/captcha", corsMiddleware(server.generateCaptcha))
	http.HandleFunc("/validate", corsMiddleware(server.validateCaptcha))
	http.HandleFunc("/status", corsMiddleware(server.apiStatus))

	// Start cleanup routine
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			server.cleanupSessions()
		}
	}()

	port := ":8080"
	fmt.Printf("🚀 SVG Math Captcha Server starting on http://localhost%s\n", port)
	fmt.Printf("📱 Visit http://localhost%s for the demo\n", port)
	fmt.Printf("🔍 API Status: http://localhost%s/status\n", port)
	fmt.Printf("📊 Captcha API: http://localhost%s/captcha\n", port)
	fmt.Printf("✅ Validate API: http://localhost%s/validate\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
