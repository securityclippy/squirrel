import { env } from '$env/dynamic/public';

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

export interface Reminder {
	id: number;
	user_id: string;
	title: string;
	description?: string;
	scheduled_at: string;
	reminder_type: 'one-time' | 'persistent';
	notification_channels: string[];
	scheduled_time: string;
	scheduled_days_of_week?: number[];
	delivery_window_minutes: number;
	delivery_method: string;
	delivery_address: string;
	status: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateReminderRequest {
	title: string;
	description?: string;
	scheduled_at: string;
	reminder_type: 'one-time' | 'persistent';
	notification_channels: string[];
	scheduled_time: string;
	scheduled_days_of_week?: number[];
	delivery_window_minutes: number;
	delivery_method: string;
	delivery_address: string;
}

export interface UpdateReminderRequest {
	title?: string;
	description?: string;
	scheduled_at?: string;
	reminder_type?: 'one-time' | 'persistent';
	notification_channels?: string[];
	scheduled_time?: string;
	scheduled_days_of_week?: number[];
	delivery_window_minutes?: number;
	delivery_method?: string;
	delivery_address?: string;
	is_active?: boolean;
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
		const response = await fetch(`${API_BASE}/reminders`);
		return handleResponse(response);
	},

	async getReminder(id: number): Promise<Reminder> {
		const response = await fetch(`${API_BASE}/reminders/${id}`);
		return handleResponse(response);
	},

	async createReminder(data: CreateReminderRequest): Promise<Reminder> {
		console.log('Creating reminder with data:', data);
		console.log('Sending POST to:', `${API_BASE}/reminders`);
		try {
			const response = await fetch(`${API_BASE}/reminders`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
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
		const response = await fetch(`${API_BASE}/reminders/${id}`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		});
		return handleResponse(response);
	},

	async deleteReminder(id: number): Promise<void> {
		const response = await fetch(`${API_BASE}/reminders/${id}`, {
			method: 'DELETE',
		});
		return handleResponse(response);
	},
};