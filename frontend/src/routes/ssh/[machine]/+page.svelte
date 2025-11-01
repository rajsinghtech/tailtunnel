<script lang="ts">
	import Terminal from '$lib/components/Terminal.svelte';
	import { page } from '$app/stores';

	const machine = $derived($page.params.machine);
	const userParam = $derived($page.url.searchParams.get('user'));
	const user = (userParam || 'root') as string;
	const pageTitle = $derived(`SSH: ${user}@${machine} - TailTunnel`);
</script>

<svelte:head>
	<title>{pageTitle}</title>
</svelte:head>

<div class="flex h-screen flex-col">
	<div class="border-b bg-card p-4">
		<div class="container mx-auto flex items-center justify-between">
			<div>
				<a href="/" class="text-sm text-muted-foreground hover:text-foreground">
					← Back to machines
				</a>
				<h2 class="text-xl font-semibold mt-1">SSH: {user}@{machine}</h2>
			</div>
		</div>
	</div>
	<div class="flex-1 overflow-hidden bg-black p-4">
		<Terminal {machine} {user} />
	</div>
</div>
