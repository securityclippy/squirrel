# Frontend API Configuration

The application now uses **Server-Side Rendering (SSR)** with SvelteKit, which means:

- **Browser**: Never makes direct API calls to backend
- **Frontend Server**: Handles all API communication via internal Docker network
- **Backend**: Communicates directly with frontend server in same subnet

## Configuration Options

### Server-Side API (Primary)
For SSR communication between frontend and backend containers:
```bash
BACKEND_API_URL=http://backend:8080/api
```

### Client-Side API (Legacy/Fallback)
Only used for any remaining client-side API calls:
```bash
PUBLIC_API_BASE_URL=http://localhost:8080/api
```

### Individual Components (Optional)
Configure individual parts (used only if complete URLs are not set):
```bash
PUBLIC_API_PORT=8080
PUBLIC_API_PATH=/api
PUBLIC_BACKEND_HOST=backend
```

## Environment Examples

### Local Development
```bash
PUBLIC_API_BASE_URL=http://localhost:8080/api
```

### Production
```bash
PUBLIC_API_BASE_URL=https://api.yourdomain.com/api
```

### External Server
```bash
PUBLIC_API_BASE_URL=http://146.190.41.200:8080/api
```

### Docker Internal Communication
```bash
PUBLIC_BACKEND_HOST=backend
PUBLIC_API_PORT=8080
PUBLIC_API_PATH=/api
```

## How to Set Environment Variables

### Docker Compose
Update the `docker-compose.yml` file:
```yaml
frontend:
  environment:
    - PUBLIC_API_BASE_URL=http://your-server:8080/api
```

### .env File
Create a `.env` file in the project root:
```bash
PUBLIC_API_BASE_URL=http://localhost:8080/api
```

### Runtime
Set environment variables when starting:
```bash
PUBLIC_API_BASE_URL=http://localhost:8080/api npm run dev
```

## Automatic Detection

If no environment variables are set, the frontend automatically detects:
- **Client-side**: Uses the current hostname with port 8080
- **Server-side**: Uses Docker service name `backend:8080`