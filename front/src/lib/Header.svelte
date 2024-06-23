<script>
	export let heads;
</script>

<nav class="ml-20 mt-10 mb-5">
	<ul class="flex flex-wrap">
		{#each heads as { href, title, current }}
			{#if current === true}
				<li>
					<a
						class="underline decoration-2 decoration-pink-500 duration-300 text-gray-950 me-4 hover:underline md:me-6 dark:text-gray-50"
						{href}>{title}</a
					>
				</li>
			{:else}
				<li>
					<a
						class="underline decoration-2 decoration-pink-500 duration-300 text-gray-400 me-4 hover:underline md:me-6 dark:text-gray-650"
						{href}>{title}</a
					>
				</li>
			{/if}
		{/each}
	</ul>
</nav>
