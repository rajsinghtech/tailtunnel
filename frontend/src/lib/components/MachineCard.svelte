<script lang="ts">
	import type { Machine } from '$lib/types/machine';
	import { cn } from '$lib/utils/style';

	let { machine }: { machine: Machine } = $props();
	let selectedUser = $state('root');
	let showCustomInput = $state(false);
	let customUser = $state('');

	const commonUsers = ['root', 'ubuntu', 'admin', 'ec2-user', 'custom'];

	function handleUserChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		selectedUser = target.value;
		showCustomInput = selectedUser === 'custom';
	}

	function handleConnect() {
		const user = selectedUser === 'custom' ? customUser : selectedUser;
		if (!user) return;
		window.location.href = `/ssh/${encodeURIComponent(machine.dnsName)}?user=${encodeURIComponent(user)}`;
	}
</script>

<div class={cn(
	'rounded-lg border bg-card shadow-sm transition-all hover:shadow-md hover:border-primary/50 flex flex-col h-full',
	!machine.online && 'opacity-60'
)}>
	<div class="p-4 flex-1 flex flex-col">
		<div class="flex items-start justify-between mb-3">
			<div class="flex-1 min-w-0">
				<h3 class="text-base font-semibold truncate mb-1">{machine.hostName}</h3>
				<p class="text-xs text-muted-foreground truncate">{machine.dnsName}</p>
			</div>
			<span class={cn(
				'ml-2 inline-flex items-center rounded-full px-2 py-0.5 text-xs font-semibold whitespace-nowrap',
				machine.online
					? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100'
					: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-100'
			)}>
				{machine.online ? 'Online' : 'Offline'}
			</span>
		</div>

		<div class="space-y-2 mb-3 flex-1">
			<div class="flex items-center gap-2">
				<span class="inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium">
					{machine.os}
				</span>
			</div>

			{#if machine.userDisplay || machine.userLogin}
				<div class="flex items-center gap-1 text-xs text-muted-foreground">
					<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="flex-shrink-0">
						<path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
						<circle cx="12" cy="7" r="4"></circle>
					</svg>
					<span class="truncate">{machine.userDisplay || machine.userLogin}</span>
				</div>
			{/if}

			{#if machine.tags && machine.tags.length > 0}
				<div class="flex flex-wrap gap-1">
					{#each machine.tags.slice(0, 3) as tag}
						<span class="inline-flex items-center rounded-md bg-blue-50 dark:bg-blue-950 px-1.5 py-0.5 text-xs font-medium text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800">
							{tag.replace('tag:', '')}
						</span>
					{/each}
					{#if machine.tags.length > 3}
						<span class="inline-flex items-center rounded-md bg-muted px-1.5 py-0.5 text-xs font-medium text-muted-foreground">
							+{machine.tags.length - 3}
						</span>
					{/if}
				</div>
			{/if}
		</div>

		<div class="space-y-2">
			<div class="flex items-center gap-2">
				<label for="user-{machine.nodeKey}" class="text-xs text-muted-foreground whitespace-nowrap">
					User:
				</label>
				<select
					id="user-{machine.nodeKey}"
					onchange={handleUserChange}
					disabled={!machine.online}
					class="flex-1 rounded-md border border-input bg-background px-2 py-1 text-xs focus:outline-none focus:ring-2 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
				>
					{#each commonUsers as user}
						<option value={user} selected={user === selectedUser}>
							{user}
						</option>
					{/each}
				</select>
			</div>

			{#if showCustomInput}
				<input
					type="text"
					bind:value={customUser}
					placeholder="Enter username"
					disabled={!machine.online}
					class="w-full rounded-md border border-input bg-background px-2 py-1 text-xs focus:outline-none focus:ring-2 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
				/>
			{/if}

			<button
				onclick={handleConnect}
				disabled={!machine.online}
				class={cn(
					'w-full rounded-md px-4 py-2 text-sm font-medium transition-colors',
					machine.online
						? 'bg-primary text-primary-foreground hover:bg-primary/90'
						: 'cursor-not-allowed bg-muted text-muted-foreground'
				)}
			>
				Connect SSH
			</button>
		</div>
	</div>
</div>
