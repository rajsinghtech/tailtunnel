<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { Terminal } from '@xterm/xterm';
	import { FitAddon } from '@xterm/addon-fit';
	import '@xterm/xterm/css/xterm.css';

	let { machine, user }: { machine: string; user?: string } = $props();
	const sshUser = user || 'root';

	let terminalElement: HTMLDivElement;
	let terminal: Terminal;
	let ws: WebSocket;
	let fitAddon: FitAddon;

	onMount(() => {
		terminal = new Terminal({
			cursorBlink: true,
			fontSize: 14,
			fontFamily: 'Menlo, Monaco, "Courier New", monospace',
			theme: {
				background: '#000000',
				foreground: '#ffffff'
			}
		});

		fitAddon = new FitAddon();
		terminal.loadAddon(fitAddon);
		terminal.open(terminalElement);
		fitAddon.fit();

		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		ws = new WebSocket(`${protocol}//${window.location.host}/api/ws/ssh/${machine}?user=${encodeURIComponent(sshUser)}`);

		ws.onopen = () => {
			terminal.writeln(`Connected to ${sshUser}@${machine}`);
		};

		ws.onmessage = (event) => {
			if (typeof event.data === 'string') {
				terminal.write(event.data);
			} else {
				event.data.arrayBuffer().then((buffer: ArrayBuffer) => {
					terminal.write(new Uint8Array(buffer));
				});
			}
		};

		ws.onerror = (error) => {
			terminal.writeln('\r\nWebSocket error occurred');
		};

		ws.onclose = () => {
			terminal.writeln('\r\n\r\nConnection closed. Redirecting...');
			setTimeout(() => {
				goto('/');
			}, 500);
		};

		terminal.onData((data) => {
			if (ws.readyState === WebSocket.OPEN) {
				ws.send(data);
			}
		});

		const handleResize = () => {
			fitAddon.fit();
		};
		window.addEventListener('resize', handleResize);

		return () => {
			window.removeEventListener('resize', handleResize);
		};
	});

	onDestroy(() => {
		ws?.close();
		terminal?.dispose();
	});
</script>

<div bind:this={terminalElement} class="h-full w-full"></div>
