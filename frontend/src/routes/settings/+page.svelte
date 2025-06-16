<script>
	import { onMount } from 'svelte';
	import { authStore } from '$lib/auth0';
	import AuthGuard from '$lib/components/AuthGuard.svelte';

	let loading = false;
	let error = '';
	let success = '';
	
	// Notification Settings
	let emailAddresses = [
		{ id: 1, address: 'user@example.com', verified: true, primary: true, added: '2024-01-15' },
		{ id: 2, address: 'backup@example.com', verified: false, primary: false, added: '2024-02-01' }
	];
	let phoneNumbers = [
		{ id: 1, number: '+1234567890', verified: true, primary: true, added: '2024-01-15' },
		{ id: 2, number: '+1987654321', verified: false, primary: false, added: '2024-02-10' }
	];
	let newEmail = '';
	let newPhone = '';
	
	// Table sorting
	let emailSortBy = 'added';
	let emailSortDirection = 'desc';
	let phoneSortBy = 'added';
	let phoneSortDirection = 'desc';
	
	// Notification Preferences
	let emailNotifications = true;
	let smsNotifications = true;
	let pushNotifications = false;
	let reminderTime = '09:00';
	let timezone = 'America/New_York';
	
	// General Settings
	let defaultReminderWindow = 15;
	let language = 'en';
	let theme = 'light';

	const timezones = [
		'America/New_York',
		'America/Chicago',
		'America/Denver',
		'America/Los_Angeles',
		'Europe/London',
		'Europe/Paris',
		'Asia/Tokyo',
		'Australia/Sydney'
	];

	function addEmail() {
		if (!newEmail.trim()) {
			error = 'Please enter a valid email address';
			return;
		}
		
		if (emailAddresses.some(e => e.address === newEmail.trim())) {
			error = 'Email address already exists';
			return;
		}

		const newId = Math.max(...emailAddresses.map(e => e.id), 0) + 1;
		emailAddresses = [...emailAddresses, {
			id: newId,
			address: newEmail.trim(),
			verified: false,
			primary: emailAddresses.length === 0,
			added: new Date().toISOString().split('T')[0]
		}];
		newEmail = '';
		success = 'Email address added successfully';
		setTimeout(() => { success = ''; }, 3000);
	}

	function removeEmail(emailId) {
		const primaryCount = emailAddresses.filter(e => e.primary).length;
		const emailToRemove = emailAddresses.find(e => e.id === emailId);
		
		if (primaryCount <= 1 && emailToRemove?.primary) {
			error = 'You must have at least one primary email address';
			return;
		}
		
		emailAddresses = emailAddresses.filter(e => e.id !== emailId);
		success = 'Email address removed';
		setTimeout(() => { success = ''; }, 3000);
	}

	function addPhone() {
		if (!newPhone.trim()) {
			error = 'Please enter a valid phone number';
			return;
		}
		
		if (phoneNumbers.some(p => p.number === newPhone.trim())) {
			error = 'Phone number already exists';
			return;
		}

		const newId = Math.max(...phoneNumbers.map(p => p.id), 0) + 1;
		phoneNumbers = [...phoneNumbers, {
			id: newId,
			number: newPhone.trim(),
			verified: false,
			primary: phoneNumbers.length === 0,
			added: new Date().toISOString().split('T')[0]
		}];
		newPhone = '';
		success = 'Phone number added successfully';
		setTimeout(() => { success = ''; }, 3000);
	}

	function removePhone(phoneId) {
		phoneNumbers = phoneNumbers.filter(p => p.id !== phoneId);
		success = 'Phone number removed';
		setTimeout(() => { success = ''; }, 3000);
	}

	function setPrimaryEmail(emailId) {
		emailAddresses = emailAddresses.map(e => ({
			...e,
			primary: e.id === emailId
		}));
		success = 'Primary email updated';
		setTimeout(() => { success = ''; }, 3000);
	}

	function setPrimaryPhone(phoneId) {
		phoneNumbers = phoneNumbers.map(p => ({
			...p,
			primary: p.id === phoneId
		}));
		success = 'Primary phone updated';
		setTimeout(() => { success = ''; }, 3000);
	}

	function sortEmails(column) {
		if (emailSortBy === column) {
			emailSortDirection = emailSortDirection === 'asc' ? 'desc' : 'asc';
		} else {
			emailSortBy = column;
			emailSortDirection = 'asc';
		}
		
		emailAddresses = emailAddresses.sort((a, b) => {
			let aVal = a[column];
			let bVal = b[column];
			
			if (column === 'verified' || column === 'primary') {
				aVal = aVal ? 1 : 0;
				bVal = bVal ? 1 : 0;
			}
			
			if (emailSortDirection === 'asc') {
				return aVal > bVal ? 1 : -1;
			} else {
				return aVal < bVal ? 1 : -1;
			}
		});
	}

	function sortPhones(column) {
		if (phoneSortBy === column) {
			phoneSortDirection = phoneSortDirection === 'asc' ? 'desc' : 'asc';
		} else {
			phoneSortBy = column;
			phoneSortDirection = 'asc';
		}
		
		phoneNumbers = phoneNumbers.sort((a, b) => {
			let aVal = a[column];
			let bVal = b[column];
			
			if (column === 'verified' || column === 'primary') {
				aVal = aVal ? 1 : 0;
				bVal = bVal ? 1 : 0;
			}
			
			if (phoneSortDirection === 'asc') {
				return aVal > bVal ? 1 : -1;
			} else {
				return aVal < bVal ? 1 : -1;
			}
		});
	}

	function getSortIcon(currentSort, column, direction) {
		if (currentSort !== column) return '↕️';
		return direction === 'asc' ? '↑' : '↓';
	}

	async function saveSettings() {
		try {
			loading = true;
			error = '';
			
			// TODO: Implement API call to save settings
			// await api.updateUserSettings({
			//   emailAddresses,
			//   phoneNumbers,
			//   emailNotifications,
			//   smsNotifications,
			//   pushNotifications,
			//   reminderTime,
			//   timezone,
			//   defaultReminderWindow,
			//   language,
			//   theme
			// });
			
			// Simulate API call
			await new Promise(resolve => setTimeout(resolve, 1000));
			
			success = 'Settings saved successfully!';
			setTimeout(() => { success = ''; }, 5000);
		} catch (err) {
			error = 'Failed to save settings: ' + (err.message || 'Unknown error');
		} finally {
			loading = false;
		}
	}

	function clearError() {
		error = '';
	}
</script>

<svelte:head>
	<title>Settings - Squirrel Reminders</title>
</svelte:head>

<AuthGuard>
	<div class="settings-page">
		<div class="header">
			<h1>Settings</h1>
			<p class="subtitle">Manage your notification preferences and account settings</p>
		</div>

		{#if error}
			<div class="alert alert-error">
				{error}
				<button class="close-btn" on:click={clearError}>×</button>
			</div>
		{/if}

		{#if success}
			<div class="alert alert-success">
				{success}
				<button class="close-btn" on:click={() => success = ''}>×</button>
			</div>
		{/if}

		<!-- Notification Settings Section -->
		<div class="section">
			<h2>📧 Notification Methods</h2>
			<p class="section-description">Manage how you receive reminder notifications</p>

			<!-- Email Addresses -->
			<div class="subsection">
				<h3>Email Addresses</h3>
				<div class="data-table-container">
					<table class="data-table">
						<thead>
							<tr>
								<th class="sortable" on:click={() => sortEmails('address')}>
									Email Address {getSortIcon(emailSortBy, 'address', emailSortDirection)}
								</th>
								<th class="sortable" on:click={() => sortEmails('verified')}>
									Status {getSortIcon(emailSortBy, 'verified', emailSortDirection)}
								</th>
								<th class="sortable" on:click={() => sortEmails('primary')}>
									Primary {getSortIcon(emailSortBy, 'primary', emailSortDirection)}
								</th>
								<th class="sortable" on:click={() => sortEmails('added')}>
									Added {getSortIcon(emailSortBy, 'added', emailSortDirection)}
								</th>
								<th>Actions</th>
							</tr>
						</thead>
						<tbody>
							{#each emailAddresses as email}
								<tr>
									<td class="email-cell">{email.address}</td>
									<td>
										<span class="status-badge {email.verified ? 'verified' : 'unverified'}">
											{email.verified ? '✓ Verified' : '⚠ Unverified'}
										</span>
									</td>
									<td>
										{#if email.primary}
											<span class="primary-badge">★ Primary</span>
										{:else}
											<button 
												class="btn btn-sm btn-secondary"
												on:click={() => setPrimaryEmail(email.id)}
											>
												Set Primary
											</button>
										{/if}
									</td>
									<td>{email.added}</td>
									<td>
										<button 
											class="btn btn-sm btn-danger" 
											on:click={() => removeEmail(email.id)}
											disabled={email.primary && emailAddresses.filter(e => e.primary).length <= 1}
										>
											Remove
										</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
				
				<div class="add-contact">
					<input 
						type="email" 
						placeholder="Add new email address"
						bind:value={newEmail}
						on:keydown={(e) => e.key === 'Enter' && addEmail()}
					/>
					<button class="btn btn-primary" on:click={addEmail}>Add Email</button>
				</div>
			</div>

			<!-- Phone Numbers -->
			<div class="subsection">
				<h3>Phone Numbers</h3>
				<div class="data-table-container">
					<table class="data-table">
						<thead>
							<tr>
								<th class="sortable" on:click={() => sortPhones('number')}>
									Phone Number {getSortIcon(phoneSortBy, 'number', phoneSortDirection)}
								</th>
								<th class="sortable" on:click={() => sortPhones('verified')}>
									Status {getSortIcon(phoneSortBy, 'verified', phoneSortDirection)}
								</th>
								<th class="sortable" on:click={() => sortPhones('primary')}>
									Primary {getSortIcon(phoneSortBy, 'primary', phoneSortDirection)}
								</th>
								<th class="sortable" on:click={() => sortPhones('added')}>
									Added {getSortIcon(phoneSortBy, 'added', phoneSortDirection)}
								</th>
								<th>Actions</th>
							</tr>
						</thead>
						<tbody>
							{#each phoneNumbers as phone}
								<tr>
									<td class="phone-cell">{phone.number}</td>
									<td>
										<span class="status-badge {phone.verified ? 'verified' : 'unverified'}">
											{phone.verified ? '✓ Verified' : '⚠ Unverified'}
										</span>
									</td>
									<td>
										{#if phone.primary}
											<span class="primary-badge">★ Primary</span>
										{:else}
											<button 
												class="btn btn-sm btn-secondary"
												on:click={() => setPrimaryPhone(phone.id)}
											>
												Set Primary
											</button>
										{/if}
									</td>
									<td>{phone.added}</td>
									<td>
										<button 
											class="btn btn-sm btn-danger" 
											on:click={() => removePhone(phone.id)}
										>
											Remove
										</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
				
				<div class="add-contact">
					<input 
						type="tel" 
						placeholder="Add new phone number (+1234567890)"
						bind:value={newPhone}
						on:keydown={(e) => e.key === 'Enter' && addPhone()}
					/>
					<button class="btn btn-primary" on:click={addPhone}>Add Phone</button>
				</div>
			</div>
		</div>

		<!-- Notification Preferences -->
		<div class="section">
			<h2>🔔 Notification Preferences</h2>
			<p class="section-description">Choose which types of notifications you want to receive</p>

			<div class="setting-group">
				<label class="setting-item">
					<input type="checkbox" bind:checked={emailNotifications} />
					<span class="checkmark"></span>
					<div class="setting-content">
						<strong>Email Notifications</strong>
						<p>Receive reminders via email</p>
					</div>
				</label>

				<label class="setting-item">
					<input type="checkbox" bind:checked={smsNotifications} />
					<span class="checkmark"></span>
					<div class="setting-content">
						<strong>SMS Notifications</strong>
						<p>Receive reminders via text message</p>
					</div>
				</label>

				<label class="setting-item">
					<input type="checkbox" bind:checked={pushNotifications} />
					<span class="checkmark"></span>
					<div class="setting-content">
						<strong>Push Notifications</strong>
						<p>Receive browser push notifications</p>
					</div>
				</label>
			</div>
		</div>

		<!-- Time & Location Settings -->
		<div class="section">
			<h2>🕐 Time & Location</h2>
			<p class="section-description">Set your default time preferences</p>

			<div class="form-grid">
				<div class="form-group">
					<label for="reminderTime">Default Reminder Time</label>
					<input 
						type="time" 
						id="reminderTime"
						bind:value={reminderTime}
					/>
				</div>

				<div class="form-group">
					<label for="timezone">Timezone</label>
					<select id="timezone" bind:value={timezone}>
						{#each timezones as tz}
							<option value={tz}>{tz}</option>
						{/each}
					</select>
				</div>

				<div class="form-group">
					<label for="reminderWindow">Default Reminder Window (minutes)</label>
					<input 
						type="number" 
						id="reminderWindow"
						bind:value={defaultReminderWindow}
						min="1"
						max="60"
					/>
				</div>
			</div>
		</div>

		<!-- General Settings -->
		<div class="section">
			<h2>⚙️ General</h2>
			<p class="section-description">General application preferences</p>

			<div class="form-grid">
				<div class="form-group">
					<label for="language">Language</label>
					<select id="language" bind:value={language}>
						<option value="en">English</option>
						<option value="es">Español</option>
						<option value="fr">Français</option>
						<option value="de">Deutsch</option>
					</select>
				</div>

				<div class="form-group">
					<label for="theme">Theme</label>
					<select id="theme" bind:value={theme}>
						<option value="light">Light</option>
						<option value="dark">Dark</option>
						<option value="system">System</option>
					</select>
				</div>
			</div>
		</div>

		<!-- Save Button -->
		<div class="save-section">
			<button 
				class="btn btn-primary btn-large" 
				on:click={saveSettings}
				disabled={loading}
			>
				{loading ? 'Saving...' : 'Save Settings'}
			</button>
		</div>
	</div>
</AuthGuard>

<style>
	.settings-page {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
	}

	.header {
		margin-bottom: 2rem;
	}

	.header h1 {
		margin: 0 0 0.5rem 0;
		color: #111827;
		font-size: 2rem;
		font-weight: 600;
	}

	.subtitle {
		margin: 0;
		color: #6b7280;
		font-size: 1rem;
	}

	.section {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		margin-bottom: 2rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		border: 1px solid #e5e7eb;
	}

	.section h2 {
		margin: 0 0 0.5rem 0;
		color: #374151;
		font-size: 1.5rem;
		font-weight: 600;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.section-description {
		margin: 0 0 2rem 0;
		color: #6b7280;
		font-size: 0.875rem;
	}

	.subsection {
		margin-bottom: 2rem;
	}

	.subsection:last-child {
		margin-bottom: 0;
	}

	.subsection h3 {
		margin: 0 0 1rem 0;
		color: #374151;
		font-size: 1.125rem;
		font-weight: 600;
	}

	.data-table-container {
		overflow-x: auto;
		margin-bottom: 1.5rem;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		background: white;
	}

	.data-table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.875rem;
	}

	.data-table thead {
		background: #f9fafb;
	}

	.data-table th {
		padding: 0.75rem 1rem;
		text-align: left;
		font-weight: 600;
		color: #374151;
		border-bottom: 1px solid #e5e7eb;
		white-space: nowrap;
	}

	.data-table th.sortable {
		cursor: pointer;
		user-select: none;
		transition: background-color 0.2s;
	}

	.data-table th.sortable:hover {
		background: #f3f4f6;
	}

	.data-table td {
		padding: 0.75rem 1rem;
		border-bottom: 1px solid #f3f4f6;
		vertical-align: middle;
	}

	.data-table tbody tr:hover {
		background: #f9fafb;
	}

	.data-table tbody tr:last-child td {
		border-bottom: none;
	}

	.email-cell, .phone-cell {
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
		font-weight: 500;
		color: #1f2937;
	}

	.status-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-badge.verified {
		background: #d1fae5;
		color: #065f46;
	}

	.status-badge.unverified {
		background: #fef3c7;
		color: #92400e;
	}

	.primary-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.25rem 0.5rem;
		background: #dbeafe;
		color: #1e40af;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.add-contact {
		display: flex;
		gap: 0.75rem;
	}

	.add-contact input {
		flex: 1;
		padding: 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 8px;
		font-size: 0.875rem;
	}

	.setting-group {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.setting-item {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		cursor: pointer;
		padding: 1rem;
		border-radius: 8px;
		transition: background-color 0.2s;
	}

	.setting-item:hover {
		background: #f9fafb;
	}

	.setting-item input[type="checkbox"] {
		position: absolute;
		opacity: 0;
		cursor: pointer;
	}

	.checkmark {
		width: 20px;
		height: 20px;
		background: white;
		border: 2px solid #d1d5db;
		border-radius: 4px;
		position: relative;
		transition: all 0.2s;
		flex-shrink: 0;
		margin-top: 2px;
	}

	.setting-item input[type="checkbox"]:checked + .checkmark {
		background: #3b82f6;
		border-color: #3b82f6;
	}

	.setting-item input[type="checkbox"]:checked + .checkmark:after {
		content: '';
		position: absolute;
		left: 6px;
		top: 2px;
		width: 6px;
		height: 10px;
		border: solid white;
		border-width: 0 2px 2px 0;
		transform: rotate(45deg);
	}

	.setting-content {
		flex: 1;
	}

	.setting-content strong {
		display: block;
		color: #111827;
		margin-bottom: 0.25rem;
	}

	.setting-content p {
		margin: 0;
		color: #6b7280;
		font-size: 0.875rem;
	}

	.form-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1.5rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-group label {
		font-weight: 600;
		color: #374151;
		font-size: 0.875rem;
	}

	.form-group input,
	.form-group select {
		padding: 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 8px;
		font-size: 0.875rem;
		background: white;
	}

	.form-group input:focus,
	.form-group select:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.save-section {
		display: flex;
		justify-content: center;
		padding: 2rem 0;
	}

	.alert {
		padding: 1rem;
		border-radius: 8px;
		margin-bottom: 1.5rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.alert-error {
		background: #fef2f2;
		color: #991b1b;
		border: 1px solid #fecaca;
	}

	.alert-success {
		background: #f0fdf4;
		color: #166534;
		border: 1px solid #bbf7d0;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 1.5rem;
		cursor: pointer;
		opacity: 0.7;
		margin-left: 1rem;
	}

	.close-btn:hover {
		opacity: 1;
	}

	.btn {
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 8px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		text-decoration: none;
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
	}

	.btn-primary {
		background: #3b82f6;
		color: white;
	}

	.btn-primary:hover {
		background: #2563eb;
	}

	.btn-danger {
		background: #dc2626;
		color: white;
	}

	.btn-danger:hover {
		background: #b91c1c;
	}

	.btn-sm {
		padding: 0.5rem 1rem;
		font-size: 0.75rem;
	}

	.btn-large {
		padding: 1rem 2rem;
		font-size: 1rem;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	@media (max-width: 768px) {
		.settings-page {
			padding: 1rem;
		}

		.section {
			padding: 1.5rem;
		}

		.form-grid {
			grid-template-columns: 1fr;
		}

		.add-contact {
			flex-direction: column;
		}

		.data-table {
			font-size: 0.75rem;
		}

		.data-table th,
		.data-table td {
			padding: 0.5rem 0.75rem;
		}

		.email-cell, .phone-cell {
			font-size: 0.75rem;
		}
	}
</style>