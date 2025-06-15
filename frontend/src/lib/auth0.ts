import { createAuth0Client, type Auth0Client, type User } from '@auth0/auth0-spa-js';
import { writable, type Readable } from 'svelte/store';
import { browser } from '$app/environment';
import { 
	PUBLIC_AUTH0_DOMAIN, 
	PUBLIC_AUTH0_CLIENT_ID, 
	PUBLIC_AUTH0_AUDIENCE 
} from '$env/static/public';

interface AuthStore {
	isLoading: boolean;
	isAuthenticated: boolean;
	user: User | undefined;
	error: Error | undefined;
}

const createAuthStore = () => {
	const { subscribe, set, update } = writable<AuthStore>({
		isLoading: true,
		isAuthenticated: false,
		user: undefined,
		error: undefined
	});

	return {
		subscribe,
		set,
		update
	};
};

export const authStore = createAuthStore();

let auth0Client: Auth0Client | undefined;

export const initializeAuth0 = async (): Promise<Auth0Client> => {
	if (auth0Client) {
		return auth0Client;
	}

	try {
		auth0Client = await createAuth0Client({
			domain: PUBLIC_AUTH0_DOMAIN,
			clientId: PUBLIC_AUTH0_CLIENT_ID,
			authorizationParams: {
				audience: PUBLIC_AUTH0_AUDIENCE,
				redirect_uri: browser ? window.location.origin : undefined
			},
			cacheLocation: 'localstorage',
			useRefreshTokens: true
		});

		// Check if we're returning from a redirect
		if (browser && window.location.search.includes('code=')) {
			try {
				await auth0Client.handleRedirectCallback();
				// Remove the query parameters from the URL
				window.history.replaceState({}, document.title, window.location.pathname);
			} catch (error) {
				console.error('Error handling redirect callback:', error);
				authStore.update(state => ({ ...state, error: error as Error, isLoading: false }));
				return auth0Client;
			}
		}

		// Check if user is authenticated
		const isAuthenticated = await auth0Client.isAuthenticated();
		
		let user: User | undefined;
		if (isAuthenticated) {
			user = await auth0Client.getUser();
		}

		authStore.set({
			isLoading: false,
			isAuthenticated,
			user,
			error: undefined
		});

		return auth0Client;
	} catch (error) {
		console.error('Error initializing Auth0:', error);
		authStore.set({
			isLoading: false,
			isAuthenticated: false,
			user: undefined,
			error: error as Error
		});
		throw error;
	}
};

export const login = async (redirectPath?: string) => {
	if (!auth0Client) {
		throw new Error('Auth0 client not initialized');
	}

	try {
		await auth0Client.loginWithRedirect({
			authorizationParams: {
				redirect_uri: browser ? `${window.location.origin}${redirectPath || ''}` : undefined
			}
		});
	} catch (error) {
		console.error('Error during login:', error);
		authStore.update(state => ({ ...state, error: error as Error }));
		throw error;
	}
};

export const logout = async () => {
	if (!auth0Client) {
		throw new Error('Auth0 client not initialized');
	}

	try {
		await auth0Client.logout({
			logoutParams: {
				returnTo: browser ? window.location.origin : undefined
			}
		});
	} catch (error) {
		console.error('Error during logout:', error);
		authStore.update(state => ({ ...state, error: error as Error }));
		throw error;
	}
};

export const getAccessToken = async (): Promise<string | undefined> => {
	if (!auth0Client) {
		throw new Error('Auth0 client not initialized');
	}

	try {
		const token = await auth0Client.getTokenSilently({
			authorizationParams: {
				audience: PUBLIC_AUTH0_AUDIENCE
			}
		});
		return token;
	} catch (error) {
		console.error('Error getting access token:', error);
		// If we can't get a token silently, the user might need to re-authenticate
		authStore.update(state => ({ 
			...state, 
			isAuthenticated: false,
			user: undefined,
			error: error as Error 
		}));
		return undefined;
	}
};

export const refreshAuth = async () => {
	if (!auth0Client) {
		return;
	}

	try {
		authStore.update(state => ({ ...state, isLoading: true }));
		
		const isAuthenticated = await auth0Client.isAuthenticated();
		let user: User | undefined;
		
		if (isAuthenticated) {
			user = await auth0Client.getUser();
		}

		authStore.update(state => ({
			...state,
			isLoading: false,
			isAuthenticated,
			user,
			error: undefined
		}));
	} catch (error) {
		console.error('Error refreshing auth state:', error);
		authStore.update(state => ({
			...state,
			isLoading: false,
			error: error as Error
		}));
	}
};