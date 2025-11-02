<script lang="ts">
	import MachineCard from '$lib/components/MachineCard.svelte';
	import MachineService from '$lib/services/machine-service';
	import type { Machine, MachineListResponse } from '$lib/types/machine';
	import { onMount } from 'svelte';

	const pageTitle = 'TailTunnel - SSH Machines';

	let machinesData = $state<MachineListResponse | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let searchQuery = $state('');

	const sortedAndFilteredMachines = $derived(() => {
		if (!machinesData?.machines) return [];

		const query = searchQuery.toLowerCase().trim();
		let filtered = machinesData.machines;

		if (query) {
			filtered = filtered.filter(machine =>
				machine.hostName.toLowerCase().includes(query) ||
				machine.dnsName.toLowerCase().includes(query) ||
				machine.userLogin.toLowerCase().includes(query) ||
				machine.userDisplay.toLowerCase().includes(query) ||
				machine.tags.some(tag => tag.toLowerCase().includes(query))
			);
		}

		return [...filtered].sort((a, b) =>
			a.hostName.localeCompare(b.hostName, undefined, { sensitivity: 'base' })
		);
	});

	async function loadMachines() {
		try {
			loading = true;
			error = null;
			const service = new MachineService();
			machinesData = await service.list();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load machines';
			console.error('Failed to load machines:', e);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadMachines();
	});
</script>

<svelte:head>
	<title>{pageTitle}</title>
</svelte:head>

<div class="container mx-auto p-4 md:p-6">
	<div class="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
		<div>
			<h1 class="text-2xl md:text-3xl font-bold tracking-tight">SSH Machines</h1>
			<p class="text-muted-foreground text-sm md:text-base mt-1">Connect to your Tailscale machines</p>
		</div>
		<button
			onclick={loadMachines}
			disabled={loading}
			class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50 self-start md:self-auto"
		>
			{loading ? 'Loading...' : 'Refresh'}
		</button>
	</div>

	{#if !error && machinesData && machinesData.machines.length > 0}
		<div class="mb-6">
			<div class="relative">
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search by name, user, or tag..."
					class="w-full rounded-md border border-input bg-background px-4 py-2.5 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
				/>
				{#if searchQuery}
					<button
						onclick={() => searchQuery = ''}
						aria-label="Clear search"
						class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
					>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<line x1="18" y1="6" x2="6" y2="18"></line>
							<line x1="6" y1="6" x2="18" y2="18"></line>
						</svg>
					</button>
				{/if}
			</div>
			{#if searchQuery && sortedAndFilteredMachines().length > 0}
				<p class="mt-2 text-sm text-muted-foreground">
					Found {sortedAndFilteredMachines().length} machine{sortedAndFilteredMachines().length === 1 ? '' : 's'}
				</p>
			{/if}
		</div>
	{/if}

	{#if error}
		<div class="rounded-lg border border-destructive bg-destructive/10 p-4 text-destructive">
			<p class="font-semibold">Error</p>
			<p class="text-sm">{error}</p>
		</div>
	{:else if loading && !machinesData}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Loading machines...</p>
		</div>
	{:else if machinesData}
		{#if machinesData.machines.length === 0}
			<div class="rounded-lg border bg-card p-8 text-center">
				<p class="text-muted-foreground">No SSH-enabled machines found on your tailnet</p>
			</div>
		{:else if sortedAndFilteredMachines().length === 0}
			<div class="rounded-lg border bg-card p-8 text-center">
				<p class="text-muted-foreground">No machines match your search</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
				{#each sortedAndFilteredMachines() as machine (machine.nodeKey)}
					<MachineCard {machine} />
				{/each}
			</div>
		{/if}
	{/if}
</div>
