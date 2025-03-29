# MarloweQuill: AI-Powered Content Generation Platform

MarloweQuill is a powerful content generation platform that leverages Google's Gemini AI to create various types of content including blog posts, articles, social media posts, and scripts. The application provides both a web interface and API endpoints for seamless content generation.

##  Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Folder Structure](#folder-structure)
- [API Endpoints](#api-endpoints)
- [Setup & Installation](#setup--installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Rate Limiting](#rate-limiting)

##  Overview

MarloweQuill provides an intuitive interface for generating high-quality content using Google's Gemini AI models. Users can specify content type, topic, tone, length, and additional context to get customized content. The application handles API interactions, content generation, and file storage, making it easy to create and manage various content pieces.

## Features

- **Multiple Content Types**: Generate blog posts, articles, social media posts, and scripts
- **Customization Options**: Specify topic, tone, length, and additional context
- **Web Interface**: User-friendly client application for content creation
- **API Access**: REST API endpoints for programmatic content generation
- **Rate Limiting**: Built-in protection against API abuse
- **Content Storage**: Automatic saving of generated content to files
- **Error Handling**: Robust error handling and user feedback

##  Folder Structure

```
/home/reek/Projects/MarloweQuill/
├── cmd/
│   ├── marlowequill/     # Main application entry point
│   └── test/             # Test utilities
├── internal/
│   ├── api/              # Gemini API integration
│   ├── config/           # Configuration management
│   ├── content/          # Content generation logic
│   ├── middleware/       # HTTP middleware (rate limiting)
│   └── server/           # HTTP server implementation
├── client/               # Web client application
│   ├── css/              # Stylesheets
│   ├── js/               # JavaScript code
│   └── index.html        # Main HTML page
├── configs/              # Configuration files
├── generated_content/    # Output directory for generated content
└── README.md             # Project documentation
```

##  API Endpoints

### `POST /api/generate`

Generates content based on the provided parameters.

**Request Body:**
```json
{
  "content_type": "blog_post",
  "topic": "Artificial Intelligence in 2024",
  "tone": "professional",
  "length": 500,
  "additional_context": "Focus on recent developments"
}
```

**Parameters:**
- `content_type` (required): Type of content to generate (blog_post, article, social_media, script)
- `topic` (required): The main topic or subject of the content
- `tone` (optional): The tone of the content (professional, casual, formal, etc.)
- `length` (optional): Target word count for the generated content
- `additional_context` (optional): Any additional context or requirements

**Response:**
```json
{
  "content": "Generated content text goes here...",
  "filename": "generated_content/blog_post_20250329_104515.txt"
}
```

### `GET /health`

Checks if the API server is running.

**Response:**
```
Server is running
```

##  Setup & Installation

### Prerequisites

- Go 1.16 or higher
- Google Gemini API key

### Installation Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/reek004/MarloweQuill.git
   cd MarloweQuill
   ```

2. Create a `.env` file in the root directory:
   ```
   GEMINI_API_KEY=your_gemini_api_key_here
   ```

3. Build the application:
   ```bash
   go build -o marlowequill ./cmd/marlowequill
   ```

4. Run the server:
   ```bash
   ./marlowequill
   ```

5. Access the web client:
   - Open `client/index.html` in a web browser
   - Alternatively, serve the client directory using a web server

## ⚙️ Configuration

Configuration is managed through the `configs/config.yaml` file:

```yaml
gemini_api_key: ${GEMINI_API_KEY}
default_models:
  - gemini-1.5-pro
max_tokens: 1024
rate_limit:
  requests_per_minute: 5
```

- `gemini_api_key`: API key for Gemini (can be overridden via environment variable)
- `default_models`: List of AI models to try in order of preference
- `max_tokens`: Maximum tokens for content generation
- `rate_limit.requests_per_minute`: Maximum API requests allowed per minute per IP

## 🚀 Usage

### Using the Web Interface

1. Open the MarloweQuill client in a web browser
2. Select the content type (blog post, article, social media, script)
3. Enter the required topic and optional parameters
4. Click "Generate Content"
5. View the generated content in the preview tab
6. The content is automatically saved to the `generated_content` directory

### Using the Client Application

#### Setting Up the Client

1. **Serve the client files**: You can use any web server to serve the files in the `client` directory. For example:
   ```bash
   # Using Python's built-in HTTP server
   cd /path/to/MarloweQuill/client
   python3 -m http.server 8000
   
   # Or using Node.js with http-server
   npm install -g http-server
   cd /path/to/MarloweQuill/client
   http-server -p 8000
   ```

2. **Access the client**: Open your web browser and navigate to:
   ```
   http://localhost:8000
   ```

3. **Ensure the API is running**: The client needs to connect to the MarloweQuill API server. Make sure the server is running on port 8080 (default) or update the `API_URL` in `client/js/app.js` if using a different port.

#### Using the Interface

1. **Select content type**: Choose the type of content you want to generate from the dropdown menu:
   - Blog Post: Longer-form content suitable for blogs (500-1000 words)
   - Article: Structured content with sections and details
   - Social Media: Short, engaging posts for platforms like Twitter or Instagram
   - Script: Dialogue or presentation scripts

2. **Enter a topic**: Type in the main subject or topic for your content.

3. **Specify tone (optional)**: Set the tone of the content, such as:
   - Professional
   - Casual
   - Formal
   - Humorous
   - Educational
   - Inspirational

4. **Adjust length (optional)**: Set the target word count for your content.

5. **Add context (optional)**: Provide additional details, requirements, or specific points you want included in the generated content.

6. **Generate content**: Click the "Generate Content" button and wait for the response. The generation process typically takes 5-15 seconds depending on the content length and server load.

7. **View results**: Once generation is complete, you'll see:
   - A success message indicating where the content was saved
   - The generated content in the "Content Preview" tab
   - The full API response in the "Full Response" tab

8. **Copy or use the content**: You can copy the generated text directly from the preview window. The content is also automatically saved to a file in the `generated_content` directory on the server.

#### Tips for Better Results

- **Be specific with your topic**: Instead of "Technology," try "The Impact of Quantum Computing on Cybersecurity in 2024"
- **Use the tone parameter**: Setting an appropriate tone significantly improves the context-appropriateness of the content
- **Provide context**: Use the additional context field to specify angles, points to cover, or information to include
- **Adjust length based on content type**: Blog posts work well at 600-1000 words, while social media posts should be much shorter (50-100 words)

#### Troubleshooting

- **API Connection Error**: If you see "Cannot connect to API server," make sure the MarloweQuill server is running
- **Rate Limit Exceeded**: The default setting limits requests to 5 per minute per IP; wait and try again
- **Content Not Generated**: Check the API response tab for detailed error messages
- **CORS Issues**: If you're accessing the client from a different domain, you may need to adjust the CORS settings in the server

### Using the API Directly

```bash
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "content_type": "blog_post",
    "topic": "Artificial Intelligence in 2024",
    "tone": "professional",
    "length": 500,
    "additional_context": "Focus on recent developments"
  }'
```

##  Rate Limiting

MarloweQuill implements rate limiting to prevent API abuse:

- Default: 5 requests per minute per IP address
- Rate limit can be configured in `config.yaml`
- When exceeded, the server returns HTTP 429 with Retry-After header
- Client displays user-friendly messages for rate limiting errors

##  Security Considerations

- The API server is configured with CORS to allow cross-origin requests
- Rate limiting helps prevent denial-of-service attacks
- API keys are stored in environment variables for security

## Contribution & Development

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
