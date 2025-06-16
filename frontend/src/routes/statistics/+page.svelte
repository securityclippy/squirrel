<script lang="ts">
	import { authStore } from '$lib/auth0';
	import AuthGuard from '$lib/components/AuthGuard.svelte';
	import type { PageData } from './$types';

	export let data: PageData;

	// Extract statistics from server data
	$: statistics = data.statistics;
	$: loadError = data.error;
	
	// Derived stats for easier access
	$: accountStats = {
		totalReminders: statistics.overview.total_reminders,
		activeReminders: statistics.overview.active_reminders,
		completedReminders: statistics.overview.completed_reminders,
		memberSince: statistics.overview.member_since
	};

	$: weeklyStats = {
		thisWeek: statistics.weekly.this_week,
		lastWeek: statistics.weekly.last_week,
		changePercent: statistics.weekly.change_percent
	};

	$: monthlyStats = {
		thisMonth: statistics.monthly.this_month,
		lastMonth: statistics.monthly.last_month,
		changePercent: statistics.monthly.change_percent
	};

	$: categoryStats = statistics.categories;

	function formatMemberSince(dateString) {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long'
		});
	}

	function getChangeColor(changePercent) {
		if (changePercent > 0) return '#059669'; // green
		if (changePercent < 0) return '#dc2626'; // red
		return '#6b7280'; // gray
	}

	function getChangeIcon(changePercent) {
		if (changePercent > 0) return '📈';
		if (changePercent < 0) return '📉';
		return '➡️';
	}
</script>

<svelte:head>
	<title>Statistics - Squirrel Reminders</title>
</svelte:head>

<AuthGuard>
	<div class="statistics-page">
		<div class="header">
			<h1>📊 Statistics</h1>
			<p class="subtitle">Track your reminder activity and progress</p>
		</div>

		{#if loadError}
			<div class="alert alert-error">
				{loadError}
			</div>
		{/if}

		<!-- Overview Statistics -->
		<div class="section">
			<h2>📈 Overview</h2>
			<div class="stats-grid">
				<div class="stat-card primary">
					<div class="stat-number">{accountStats.totalReminders}</div>
					<div class="stat-label">Total Reminders</div>
				</div>
				
				<div class="stat-card success">
					<div class="stat-number">{accountStats.activeReminders}</div>
					<div class="stat-label">Active Reminders</div>
				</div>
				
				<div class="stat-card info">
					<div class="stat-number">{accountStats.completedReminders}</div>
					<div class="stat-label">Completed</div>
				</div>
				
				<div class="stat-card neutral">
					<div class="stat-number">{formatMemberSince(accountStats.memberSince)}</div>
					<div class="stat-label">Member Since</div>
				</div>
			</div>
		</div>

		<!-- Time-based Statistics -->
		<div class="section">
			<h2>📅 Activity Trends</h2>
			<div class="trends-grid">
				<div class="trend-card">
					<div class="trend-header">
						<h3>This Week</h3>
						<span class="trend-icon">{getChangeIcon(weeklyStats.changePercent)}</span>
					</div>
					<div class="trend-number">{weeklyStats.thisWeek}</div>
					<div class="trend-change" style="color: {getChangeColor(weeklyStats.changePercent)}">
						{weeklyStats.changePercent > 0 ? '+' : ''}{weeklyStats.changePercent}% from last week
					</div>
				</div>

				<div class="trend-card">
					<div class="trend-header">
						<h3>This Month</h3>
						<span class="trend-icon">{getChangeIcon(monthlyStats.changePercent)}</span>
					</div>
					<div class="trend-number">{monthlyStats.thisMonth}</div>
					<div class="trend-change" style="color: {getChangeColor(monthlyStats.changePercent)}">
						{monthlyStats.changePercent > 0 ? '+' : ''}{monthlyStats.changePercent}% from last month
					</div>
				</div>
			</div>
		</div>

		<!-- Category Breakdown -->
		<div class="section">
			<h2>🏷️ Categories</h2>
			<div class="category-stats">
				{#each categoryStats as category}
					<div class="category-item">
						<div class="category-info">
							<span class="category-name">{category.name}</span>
							<span class="category-count">{category.count} reminders</span>
						</div>
						<div class="category-bar">
							<div class="category-fill" style="width: {category.percentage}%"></div>
						</div>
						<span class="category-percentage">{category.percentage}%</span>
					</div>
				{:else}
					<div class="empty-state">
						<p>No reminder categories found. Create some reminders to see category breakdown!</p>
					</div>
				{/each}
			</div>
		</div>

	</div>
</AuthGuard>

<style>
	.statistics-page {
		max-width: 1000px;
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
		margin: 0 0 1.5rem 0;
		color: #374151;
		font-size: 1.5rem;
		font-weight: 600;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1.5rem;
	}

	.stat-card {
		border-radius: 12px;
		padding: 1.5rem;
		text-align: center;
		border: 1px solid #e2e8f0;
	}

	.stat-card.primary {
		background: linear-gradient(135deg, #ddd6fe, #c7d2fe);
		border-color: #c7d2fe;
	}

	.stat-card.success {
		background: linear-gradient(135deg, #dcfce7, #bbf7d0);
		border-color: #bbf7d0;
	}

	.stat-card.info {
		background: linear-gradient(135deg, #dbeafe, #bfdbfe);
		border-color: #bfdbfe;
	}

	.stat-card.neutral {
		background: linear-gradient(135deg, #f8fafc, #e2e8f0);
		border-color: #e2e8f0;
	}

	.stat-number {
		font-size: 2rem;
		font-weight: 700;
		color: #1e293b;
		margin-bottom: 0.5rem;
	}

	.stat-label {
		color: #64748b;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.trends-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.trend-card {
		background: #f8fafc;
		border-radius: 12px;
		padding: 1.5rem;
		border: 1px solid #e2e8f0;
	}

	.trend-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}

	.trend-header h3 {
		margin: 0;
		color: #374151;
		font-size: 1.125rem;
		font-weight: 600;
	}

	.trend-icon {
		font-size: 1.5rem;
	}

	.trend-number {
		font-size: 2.5rem;
		font-weight: 700;
		color: #1e293b;
		margin-bottom: 0.5rem;
	}

	.trend-change {
		font-size: 0.875rem;
		font-weight: 500;
	}

	.category-stats {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.category-item {
		display: grid;
		grid-template-columns: 1fr 2fr auto;
		align-items: center;
		gap: 1rem;
		padding: 1rem;
		background: #f9fafb;
		border-radius: 8px;
		border: 1px solid #e5e7eb;
	}

	.category-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.category-name {
		font-weight: 600;
		color: #374151;
	}

	.category-count {
		font-size: 0.875rem;
		color: #6b7280;
	}

	.category-bar {
		background: #e5e7eb;
		height: 8px;
		border-radius: 4px;
		overflow: hidden;
	}

	.category-fill {
		height: 100%;
		background: linear-gradient(90deg, #3b82f6, #8b5cf6);
		transition: width 0.3s ease;
	}

	.category-percentage {
		font-weight: 600;
		color: #374151;
		font-size: 0.875rem;
		min-width: 40px;
		text-align: right;
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

	@media (max-width: 768px) {
		.statistics-page {
			padding: 1rem;
		}

		.section {
			padding: 1.5rem;
		}

		.stats-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.trends-grid {
			grid-template-columns: 1fr;
		}

		.category-item {
			grid-template-columns: 1fr;
			gap: 0.75rem;
		}

		.category-bar {
			order: 3;
		}

		.category-percentage {
			text-align: left;
		}

	}

	@media (max-width: 480px) {
		.stats-grid {
			grid-template-columns: 1fr;
		}
	}
</style>