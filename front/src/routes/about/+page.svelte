<script>
	let description = 'About status';
	import '../app.css';
	import Header from '$lib/Header.svelte';
</script>

<svelte:head>
	<title>Status · About</title>
	<meta name="description" content={description} />
</svelte:head>

<Header
	heads={[
		{ href: '/', title: 'Home', current: false },
		{ href: '/about', title: 'About', current: true }
	]}
/>

<div class="m-20 text-center text-gray-950 dark:text-gray-50">
	<div class="pb-10 text-3xl font-semibold">
		<h1>{description}</h1>
	</div>
	<div>
		<p>
			This is a minimal status page services. You can find source on <a
				href="https://github.com/r4iponce/status"
				class="underline decoration-2 decoration-pink-500 duration-300 hover:bg-pink-500">gitlab</a
			>. This is build with
			<a
				href="https://go.dev"
				class="underline decoration-2 decoration-pink-500 duration-300 hover:bg-pink-500">go</a
			>
			and
			<a
				href="https://svelte.dev"
				class="underline decoration-2 decoration-pink-500 duration-300 hover:bg-pink-500">svelte</a
			>.
		</p>
	</div>
</div>
