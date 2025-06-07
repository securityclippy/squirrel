import { env } from '$env/dynamic/private';

// Server-side API client for internal Docker communication
const getServerApiBase = () => {
	// Use environment variable or default to Docker service communication
	return env.BACKEND_API_URL || 'http://backend:8080/api';
};

const SERVER_API_BASE = getServerApiBase();

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

class ServerApiError extends Error {
	constructor(message: string, public status: number, public body?: any) {
		super(message);
		this.name = 'ServerApiError';
	}
}

async function handleServerResponse<T>(response: Response): Promise<T> {
	if (!response.ok) {
		let errorMessage = `HTTP ${response.status}`;
		let errorBody;
		
		try {
			errorBody = await response.text();
			errorMessage = errorBody || errorMessage;
		} catch {
			// If we can't read the response body, use the status
		}
		
		throw new ServerApiError(errorMessage, response.status, errorBody);
	}
	
	if (response.status === 204) {
		return undefined as T;
	}
	
	return response.json();
}

export const serverApi = {
	async getReminders(): Promise<{ reminders: Reminder[]; count: number }> {
		console.log('Server API: Fetching reminders from', `${SERVER_API_BASE}/reminders`);
		const response = await fetch(`${SERVER_API_BASE}/reminders`);
		return handleServerResponse(response);
	},

	async getReminder(id: number): Promise<Reminder> {
		console.log('Server API: Fetching reminder', id, 'from', `${SERVER_API_BASE}/reminders/${id}`);
		const response = await fetch(`${SERVER_API_BASE}/reminders/${id}`);
		return handleServerResponse(response);
	},

	async createReminder(data: CreateReminderRequest): Promise<Reminder> {
		console.log('Server API: Creating reminder at', `${SERVER_API_BASE}/reminders`);
		console.log('Server API: Request data:', JSON.stringify(data, null, 2));
		
		const response = await fetch(`${SERVER_API_BASE}/reminders`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		});
		
		console.log('Server API: Response status:', response.status);
		console.log('Server API: Response headers:', Object.fromEntries(response.headers.entries()));
		
		// Log the raw response text before parsing
		const responseText = await response.text();
		console.log('Server API: Raw response text:', responseText);
		
		// Create a new Response object since we've consumed the body
		const newResponse = new Response(responseText, {
			status: response.status,
			statusText: response.statusText,
			headers: response.headers,
		});
		
		return handleServerResponse(newResponse);
	},

	async updateReminder(id: number, data: UpdateReminderRequest): Promise<Reminder> {
		console.log('Server API: Updating reminder', id, data);
		const response = await fetch(`${SERVER_API_BASE}/reminders/${id}`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		});
		return handleServerResponse(response);
	},

	async deleteReminder(id: number): Promise<void> {
		console.log('Server API: Deleting reminder', id);
		const response = await fetch(`${SERVER_API_BASE}/reminders/${id}`, {
			method: 'DELETE',
		});
		return handleServerResponse(response);
	},
};