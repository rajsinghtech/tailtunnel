<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Chart, LineController, LineElement, PointElement, LinearScale, TimeScale, Title, Tooltip, Legend } from 'chart.js';
	import 'chartjs-adapter-date-fns';
	import type { PingResult } from '$lib/types/canary';

	Chart.register(LineController, LineElement, PointElement, LinearScale, TimeScale, Title, Tooltip, Legend);

	let { pingHistory }: { pingHistory: Map<string, Array<{ timestamp: Date; latency: number; connectionType: string }>> } = $props();

	let canvas: HTMLCanvasElement;
	let chart: Chart | null = null;

	function updateChart() {
		if (!chart || !canvas) return;

		const datasets = [];
		const colors = [
			'rgb(59, 130, 246)',   // blue
			'rgb(16, 185, 129)',   // green
			'rgb(245, 158, 11)',   // amber
			'rgb(239, 68, 68)',    // red
			'rgb(139, 92, 246)',   // purple
			'rgb(236, 72, 153)',   // pink
			'rgb(20, 184, 166)',   // teal
			'rgb(251, 146, 60)',   // orange
		];

		let colorIndex = 0;
		const sortedHosts = Array.from(pingHistory.keys()).sort();

		for (const [hostname, history] of pingHistory) {
			if (history.length === 0) continue;

			const color = colors[colorIndex % colors.length];
			colorIndex++;

			datasets.push({
				label: hostname,
				data: history.map(h => ({
					x: h.timestamp.getTime(),
					y: h.latency
				})),
				borderColor: color,
				backgroundColor: color.replace('rgb', 'rgba').replace(')', ', 0.1)'),
				borderWidth: 2,
				pointRadius: 2,
				pointHoverRadius: 4,
				tension: 0.3,
			});
		}

		chart.data.datasets = datasets;
		chart.update('none');
	}

	onMount(() => {
		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		chart = new Chart(ctx, {
			type: 'line',
			data: {
				datasets: []
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				interaction: {
					mode: 'nearest',
					axis: 'x',
					intersect: false
				},
				plugins: {
					title: {
						display: true,
						text: 'Latency Over Time',
						color: 'rgb(148, 163, 184)',
						font: {
							size: 16,
							weight: 'bold'
						}
					},
					legend: {
						display: true,
						position: 'bottom',
						labels: {
							color: 'rgb(148, 163, 184)',
							usePointStyle: true,
							padding: 10,
							font: {
								size: 11
							}
						},
						onClick: (e, legendItem, legend) => {
							const index = legendItem.datasetIndex;
							const ci = legend.chart;

							// Count how many are currently visible
							let visibleCount = 0;
							for (let i = 0; i < ci.data.datasets.length; i++) {
								if (ci.isDatasetVisible(i)) {
									visibleCount++;
								}
							}

							// If all are visible (initial state), hide all except the clicked one
							if (visibleCount === ci.data.datasets.length) {
								for (let i = 0; i < ci.data.datasets.length; i++) {
									ci.setDatasetVisibility(i, i === index);
								}
							} else {
								// Otherwise, toggle the clicked dataset
								const isCurrentlyVisible = ci.isDatasetVisible(index);
								ci.setDatasetVisibility(index, !isCurrentlyVisible);

								// Check if we just turned everything off
								let anyVisible = false;
								for (let i = 0; i < ci.data.datasets.length; i++) {
									if (ci.isDatasetVisible(i)) {
										anyVisible = true;
										break;
									}
								}

								// If nothing is visible, show all
								if (!anyVisible) {
									for (let i = 0; i < ci.data.datasets.length; i++) {
										ci.setDatasetVisibility(i, true);
									}
								}
							}

							ci.update();
						}
					},
					tooltip: {
						callbacks: {
							label: function(context) {
								return `${context.dataset.label}: ${context.parsed.y.toFixed(1)}ms`;
							}
						}
					}
				},
				scales: {
					x: {
						type: 'time',
						time: {
							unit: 'minute',
							displayFormats: {
								minute: 'HH:mm'
							}
						},
						title: {
							display: true,
							text: 'Time',
							color: 'rgb(148, 163, 184)'
						},
						ticks: {
							color: 'rgb(148, 163, 184)'
						},
						grid: {
							color: 'rgba(148, 163, 184, 0.1)'
						}
					},
					y: {
						beginAtZero: true,
						title: {
							display: true,
							text: 'Latency (ms)',
							color: 'rgb(148, 163, 184)'
						},
						ticks: {
							color: 'rgb(148, 163, 184)'
						},
						grid: {
							color: 'rgba(148, 163, 184, 0.1)'
						}
					}
				}
			}
		});

		updateChart();
	});

	$effect(() => {
		// Track the map and its size to ensure reactivity
		const _ = pingHistory.size;
		const __ = Array.from(pingHistory.keys());
		updateChart();
	});

	onDestroy(() => {
		if (chart) {
			chart.destroy();
		}
	});
</script>

<div class="rounded-lg border bg-card p-4">
	<canvas bind:this={canvas} style="height: 300px;"></canvas>
</div>
