import { env } from '$env/dynamic/public';
import { getAccessToken } from './auth0';

// Get API base URL from environment variable or use dynamic detection
const getApiBase = () => {
	// First check for explicit environment variable
	if (env.PUBLIC_API_BASE_URL) {
		return env.PUBLIC_API_BASE_URL;
	}
	
	// Get configurable values from environment or use defaults
	const apiPort = env.PUBLIC_API_PORT || '8080';
	const apiPath = env.PUBLIC_API_PATH || '/api';
	
	if (typeof window !== 'undefined') {
		// Client-side: use the current host or localhost
		const host = window.location.hostname;
		if (host === 'localhost' || host === '127.0.0.1') {
			return `http://localhost:${apiPort}${apiPath}`;
		} else {
			return `http://${host}:${apiPort}${apiPath}`;
		}
	}
	// Server-side: use backend service name for Docker
	const backendHost = env.PUBLIC_BACKEND_HOST || 'backend';
	return `http://${backendHost}:${apiPort}${apiPath}`;
};

const API_BASE = getApiBase();

// Helper function to get authentication headers
async function getAuthHeaders(): Promise<HeadersInit> {
	const headers: HeadersInit = {
		'Content-Type': 'application/json',
	};

	try {
		const token = await getAccessToken();
		if (token) {
			headers['Authorization'] = `Bearer ${token}`;
		}
	} catch (error) {
		console.warn('Failed to get access token:', error);
		// Continue without auth header - backend will return 401
	}

	return headers;
}

export interface Reminder {
	id: number;
	user_id: number;
	title: string;
	description?: string;
	scheduled_at: string;
	reminder_type: 'one-time' | 'persistent' | 'recurring';
	notification_channels: string[];
	scheduled_time: string;
	scheduled_days_of_week?: number[];
	delivery_window_minutes: number;
	delivery_method: string;
	delivery_address: string;
	status: string;
	is_active: boolean;
	is_persistent: boolean;
	acknowledged_at?: string;
	reminder_interval_minutes?: number;
	last_reminded_at?: string;
	created_at: string;
	updated_at: string;
}

export interface CreateReminderRequest {
	title: string;
	description?: string;
	scheduled_at: string;
	reminder_type: 'one-time' | 'persistent' | 'recurring';
	notification_channels: string[];
	scheduled_time: string;
	scheduled_days_of_week?: number[];
	delivery_window_minutes: number;
	delivery_method: string;
	delivery_address: string;
	is_persistent: boolean;
	reminder_interval_minutes?: number;
}

export interface UpdateReminderRequest {
	title?: string;
	description?: string;
	scheduled_at?: string;
	reminder_type?: 'one-time' | 'persistent' | 'recurring';
	notification_channels?: string[];
	scheduled_time?: string;
	scheduled_days_of_week?: number[];
	delivery_window_minutes?: number;
	delivery_method?: string;
	delivery_address?: string;
	is_active?: boolean;
	is_persistent?: boolean;
	reminder_interval_minutes?: number;
}

export interface User {
	id: number;
	auth0_id: string;
	email: string;
	name: string;
	picture?: string;
	email_verified: boolean;
	created_at: string;
	updated_at: string;
	last_login_at?: string;
}

export interface UpdateUserRequest {
	name?: string;
	picture?: string;
	email_verified?: boolean;
	last_login_at?: string;
}

export interface APIKey {
	id: number;
	user_id: number;
	name: string;
	key_prefix: string;
	permissions: string[];
	expires_at?: string;
	last_used_at?: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateAPIKeyRequest {
	name: string;
	permissions: string[];
	expires_at?: string;
}

export interface CreateAPIKeyResponse {
	api_key: APIKey;
	key: string;
}

class ApiError extends Error {
	constructor(message: string, public status: number) {
		super(message);
		this.name = 'ApiError';
	}
}

async function handleResponse<T>(response: Response): Promise<T> {
	if (!response.ok) {
		const errorText = await response.text();
		throw new ApiError(errorText || `HTTP ${response.status}`, response.status);
	}
	
	if (response.status === 204) {
		return undefined as T;
	}
	
	return response.json();
}

export const api = {
	async getReminders(): Promise<{ reminders: Reminder[]; count: number }> {
		console.log('API_BASE:', API_BASE);
		console.log('Making request to:', `${API_BASE}/reminders`);
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/reminders`, { headers });
		return handleResponse(response);
	},

	async getReminder(id: number): Promise<Reminder> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/reminders/${id}`, { headers });
		return handleResponse(response);
	},

	async createReminder(data: CreateReminderRequest): Promise<Reminder> {
		console.log('Creating reminder with data:', data);
		console.log('Sending POST to:', `${API_BASE}/reminders`);
		try {
			const headers = await getAuthHeaders();
			const response = await fetch(`${API_BASE}/reminders`, {
				method: 'POST',
				headers,
				body: JSON.stringify(data),
			});
			console.log('Response status:', response.status);
			console.log('Response headers:', response.headers);
			return handleResponse(response);
		} catch (error) {
			console.error('Fetch error:', error);
			throw error;
		}
	},

	async updateReminder(id: number, data: UpdateReminderRequest): Promise<Reminder> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/reminders/${id}`, {
			method: 'PUT',
			headers,
			body: JSON.stringify(data),
		});
		return handleResponse(response);
	},

	async deleteReminder(id: number): Promise<void> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/reminders/${id}`, {
			method: 'DELETE',
			headers
		});
		return handleResponse(response);
	},

	async acknowledgeReminder(id: number): Promise<void> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/reminders/${id}/acknowledge`, {
			method: 'POST',
			headers
		});
		return handleResponse(response);
	},

	// User management
	async getUserProfile(): Promise<User> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/users/profile`, { headers });
		return handleResponse(response);
	},

	async updateUserProfile(data: UpdateUserRequest): Promise<User> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/users/profile`, {
			method: 'PUT',
			headers,
			body: JSON.stringify(data),
		});
		return handleResponse(response);
	},

	// API Key management
	async getAPIKeys(): Promise<{ api_keys: APIKey[]; count: number }> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/users/api-keys`, { headers });
		return handleResponse(response);
	},

	async createAPIKey(data: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/users/api-keys`, {
			method: 'POST',
			headers,
			body: JSON.stringify(data),
		});
		return handleResponse(response);
	},

	async revokeAPIKey(id: number): Promise<void> {
		const headers = await getAuthHeaders();
		const response = await fetch(`${API_BASE}/users/api-keys/${id}`, {
			method: 'DELETE',
			headers
		});
		return handleResponse(response);
	},
};